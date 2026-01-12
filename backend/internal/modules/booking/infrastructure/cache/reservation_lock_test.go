package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockRedisClient is a mock of platformRedis.Client
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return m.Called(ctx, key, value, ttl).Error(0)
}
func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}
func (m *MockRedisClient) Del(ctx context.Context, keys ...string) error {
	return m.Called(ctx, keys).Error(0)
}
func (m *MockRedisClient) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}
func (m *MockRedisClient) Incr(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return int64(args.Int(0)), args.Error(1)
}
func (m *MockRedisClient) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return m.Called(ctx, key, ttl).Error(0)
}
func (m *MockRedisClient) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	args := m.Called(ctx, key, value, ttl)
	return args.Bool(0), args.Error(1)
}
func (m *MockRedisClient) Scan(ctx context.Context, pattern string) ([]string, error) {
	args := m.Called(ctx, pattern)
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockRedisClient) Publish(ctx context.Context, channel string, message interface{}) error {
	return m.Called(ctx, channel, message).Error(0)
}
func (m *MockRedisClient) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return m.Called(ctx, channel).Get(0).(*redis.PubSub)
}
func (m *MockRedisClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	return m.Called(ctx, script, keys, args).Get(0).(*redis.Cmd)
}
func (m *MockRedisClient) LPush(ctx context.Context, key string, values ...interface{}) error {
	return m.Called(ctx, key, values).Error(0)
}
func (m *MockRedisClient) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	args := m.Called(ctx, key, start, stop)
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockRedisClient) LTrim(ctx context.Context, key string, start, stop int64) error {
	return m.Called(ctx, key, start, stop).Error(0)
}
func (m *MockRedisClient) Ping(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}
func (m *MockRedisClient) Close() error {
	return m.Called().Error(0)
}

type ReservationLockTestSuite struct {
	suite.Suite
	mockRedis  *MockRedisClient
	rl         *ReservationLock
	ctx        context.Context
	ttl        time.Duration
	facilityID string
	userID     string
	start      time.Time
	end        time.Time
}

func (s *ReservationLockTestSuite) SetupTest() {
	s.mockRedis = new(MockRedisClient)
	s.ttl = 5 * time.Minute
	s.rl = &ReservationLock{
		redis: s.mockRedis,
		ttl:   s.ttl,
	}
	s.facilityID = "fac-1"
	s.userID = "user-1"
	s.start = time.Now()
	s.end = s.start.Add(1 * time.Hour)
	s.ctx = context.Background()
}

func (s *ReservationLockTestSuite) TestAcquire_Success() {
	s.mockRedis.On("SetNX", s.ctx, mock.Anything, mock.Anything, s.ttl).Return(true, nil).Once()
	s.mockRedis.On("Set", s.ctx, mock.Anything, mock.Anything, s.ttl).Return(nil).Once()

	lockID, err := s.rl.Acquire(s.ctx, s.facilityID, s.userID, s.start, s.end)
	s.NoError(err)
	s.NotEmpty(lockID)
}

func (s *ReservationLockTestSuite) TestAcquire_AlreadyLocked() {
	s.mockRedis.On("SetNX", s.ctx, mock.Anything, mock.Anything, s.ttl).Return(false, nil).Once()

	lockID, err := s.rl.Acquire(s.ctx, s.facilityID, s.userID, s.start, s.end)
	s.Error(err)
	s.Empty(lockID)
	s.Contains(err.Error(), "already locked")
}

func (s *ReservationLockTestSuite) TestAcquire_Error() {
	s.mockRedis.On("SetNX", s.ctx, mock.Anything, mock.Anything, s.ttl).Return(false, errors.New("redis error")).Once()

	lockID, err := s.rl.Acquire(s.ctx, s.facilityID, s.userID, s.start, s.end)
	s.Error(err)
	s.Empty(lockID)
}

func (s *ReservationLockTestSuite) TestRelease_Success() {
	s.mockRedis.On("Del", s.ctx, mock.Anything).Return(nil).Once()

	err := s.rl.Release(s.ctx, s.facilityID, s.start, s.end)
	s.NoError(err)
}

func (s *ReservationLockTestSuite) TestIsLocked_Locked() {
	s.mockRedis.On("Exists", s.ctx, mock.Anything).Return(true, nil).Once()

	locked, err := s.rl.IsLocked(s.ctx, s.facilityID, s.start, s.end)
	s.NoError(err)
	s.True(locked)
}

func (s *ReservationLockTestSuite) TestIsLocked_Unlocked() {
	s.mockRedis.On("Exists", s.ctx, mock.Anything).Return(false, nil).Once()

	locked, err := s.rl.IsLocked(s.ctx, s.facilityID, s.start, s.end)
	s.NoError(err)
	s.False(locked)
}

func (s *ReservationLockTestSuite) TestIsLocked_RedisError() {
	s.mockRedis.On("Exists", s.ctx, mock.Anything).Return(false, errors.New("redis failure")).Once()

	locked, err := s.rl.IsLocked(s.ctx, s.facilityID, s.start, s.end)
	s.NoError(err) // Fails open
	s.False(locked)
}

func (s *ReservationLockTestSuite) TestReleaseByID_Success() {
	lockID := "lock-123"
	s.mockRedis.On("Get", s.ctx, mock.Anything).Return(lockID, nil).Once()
	s.mockRedis.On("Del", s.ctx, mock.Anything).Return(nil).Twice() // One for lock, one for info

	err := s.rl.ReleaseByID(s.ctx, lockID, s.facilityID, s.start, s.end)
	s.NoError(err)
}

func (s *ReservationLockTestSuite) TestReleaseByID_OwnershipMismatch() {
	lockID := "lock-123"
	s.mockRedis.On("Get", s.ctx, mock.Anything).Return("other-lock", nil).Once()

	err := s.rl.ReleaseByID(s.ctx, lockID, s.facilityID, s.start, s.end)
	s.Error(err)
	s.Contains(err.Error(), "owned by different user")
}

func (s *ReservationLockTestSuite) TestExtendLock_Success() {
	lockID := "lock-123"
	s.mockRedis.On("Get", s.ctx, mock.Anything).Return(lockID, nil).Once()
	s.mockRedis.On("Expire", s.ctx, mock.Anything, s.ttl).Return(nil).Once()

	err := s.rl.ExtendLock(s.ctx, lockID, s.facilityID, s.start, s.end)
	s.NoError(err)
}

func (s *ReservationLockTestSuite) TestNewReservationLock() {
	lock := NewReservationLock(5 * time.Minute)
	s.NotNil(lock)
}

func TestReservationLock(t *testing.T) {
	suite.Run(t, new(ReservationLockTestSuite))
}
