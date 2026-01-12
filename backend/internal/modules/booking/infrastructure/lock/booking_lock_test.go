package lock

import (
	"context"
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

type BookingLockTestSuite struct {
	suite.Suite
	mockRedis  *MockRedisClient
	bl         *BookingLock
	ctx        context.Context
	facilityID string
	userID     string
	start      time.Time
	end        time.Time
	ttl        time.Duration
}

func (s *BookingLockTestSuite) SetupTest() {
	s.mockRedis = new(MockRedisClient)
	s.bl = &BookingLock{
		redis: s.mockRedis,
	}
	s.start = time.Now()
	s.end = s.start.Add(1 * time.Hour)
	s.ttl = 5 * time.Minute
	s.facilityID = "fac-1"
	s.userID = "user-1"
	s.ctx = context.Background()
}

func (s *BookingLockTestSuite) TestAcquireLock_Success() {
	s.mockRedis.On("SetNX", s.ctx, mock.Anything, s.userID, s.ttl).Return(true, nil).Once()

	acquired, err := s.bl.AcquireLock(s.ctx, s.facilityID, s.start, s.end, s.userID, s.ttl)
	s.NoError(err)
	s.True(acquired)
}

func (s *BookingLockTestSuite) TestAcquireLock_Failed() {
	s.mockRedis.On("SetNX", s.ctx, mock.Anything, s.userID, s.ttl).Return(false, nil).Once()

	acquired, err := s.bl.AcquireLock(s.ctx, s.facilityID, s.start, s.end, s.userID, s.ttl)
	s.NoError(err)
	s.False(acquired)
}

func (s *BookingLockTestSuite) TestReleaseLock_Success() {
	s.mockRedis.On("Del", s.ctx, mock.Anything).Return(nil).Once()

	err := s.bl.ReleaseLock(s.ctx, s.facilityID, s.start, s.end)
	s.NoError(err)
}

func (s *BookingLockTestSuite) TestIsLocked_True() {
	s.mockRedis.On("Exists", s.ctx, mock.Anything).Return(true, nil).Once()

	locked, err := s.bl.IsLocked(s.ctx, s.facilityID, s.start, s.end)
	s.NoError(err)
	s.True(locked)
}

func (s *BookingLockTestSuite) TestGetLockHolder_Success() {
	s.mockRedis.On("Get", s.ctx, mock.Anything).Return(s.userID, nil).Once()

	holder, err := s.bl.GetLockHolder(s.ctx, s.facilityID, s.start, s.end)
	s.NoError(err)
	s.Equal(s.userID, holder)
}

func (s *BookingLockTestSuite) TestExtendLock_Success() {
	cmd := redis.NewCmd(s.ctx)
	cmd.SetVal(int64(1))
	s.mockRedis.On("Eval", s.ctx, mock.Anything, mock.Anything, mock.Anything).Return(cmd).Once()

	err := s.bl.ExtendLock(s.ctx, s.facilityID, s.start, s.end, s.userID, s.ttl)
	s.NoError(err)
}

func (s *BookingLockTestSuite) TestExtendLock_Failed() {
	cmd := redis.NewCmd(s.ctx)
	cmd.SetVal(int64(0))
	s.mockRedis.On("Eval", s.ctx, mock.Anything, mock.Anything, mock.Anything).Return(cmd).Once()

	err := s.bl.ExtendLock(s.ctx, s.facilityID, s.start, s.end, s.userID, s.ttl)
	s.Error(err)
	s.Contains(err.Error(), "failed")
}

func (s *BookingLockTestSuite) TestNewBookingLock() {
	lock := NewBookingLock()
	s.NotNil(lock)
}

func TestBookingLockTestSuite(t *testing.T) {
	suite.Run(t, new(BookingLockTestSuite))
}
