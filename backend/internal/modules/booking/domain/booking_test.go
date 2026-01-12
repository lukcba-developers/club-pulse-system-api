package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuestDetails_ScanValue(t *testing.T) {
	gd := GuestDetails{
		{Name: "Guest 1", DNI: "12345678A", FeeAmount: 10.0},
	}

	// Test Value
	val, err := gd.Value()
	assert.NoError(t, err)
	assert.NotNil(t, val)

	// Test Scan
	var gd2 GuestDetails
	err = gd2.Scan(val)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(gd2))
	assert.Equal(t, gd[0].Name, gd2[0].Name)

	// Test Nil Value
	nilGd := GuestDetails(nil)
	val, err = nilGd.Value()
	assert.NoError(t, err)
	assert.Equal(t, "[]", val)

	// Test Nil Scan
	var gd3 GuestDetails
	err = gd3.Scan(nil)
	assert.NoError(t, err)
	assert.Nil(t, gd3)

	// Test Invalid Scan Type
	err = gd3.Scan(123)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed")
}

func TestBookingConstants(t *testing.T) {
	assert.Equal(t, BookingStatus("PENDING_PAYMENT"), BookingStatusPendingPayment)
	assert.Equal(t, BookingStatus("CONFIRMED"), BookingStatusConfirmed)
	assert.Equal(t, BookingStatus("CANCELLED"), BookingStatusCancelled)
}

func TestRecurrenceConstants(t *testing.T) {
	assert.Equal(t, RecurrenceType("CLASS"), RecurrenceTypeClass)
	assert.Equal(t, RecurrenceType("MAINTENANCE"), RecurrenceTypeMaintenance)
	assert.Equal(t, RecurrenceType("FIXED"), RecurrenceTypeFixed)
}
