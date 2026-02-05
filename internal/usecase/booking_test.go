package usecase

import (
	"context"
	"testing"
	"time"

	"golangHotelProject/internal/delivery/handlers/dto"
	"golangHotelProject/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBookingRepository struct {
	mock.Mock
}

func (m *MockBookingRepository) CreateBooking(ctx context.Context, b model.Booking) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockBookingRepository) GettingStatus(ctx context.Context, guestID int) (bool, error) {
	args := m.Called(ctx, guestID)
	return args.Bool(0), args.Error(1)
}

func (m *MockBookingRepository) ArrivalStatusOfRoom(ctx context.Context, roomID int) (bool, error) {
	args := m.Called(ctx, roomID)
	return args.Bool(0), args.Error(1)
}

func (m *MockBookingRepository) ReadBookingByID(ctx context.Context, id int) (model.Booking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return model.Booking{}, args.Error(1)
	}
	return args.Get(0).(model.Booking), args.Error(1)
}

func (m *MockBookingRepository) PatchBooking(ctx context.Context, patch dto.BookingPatch) error {
	args := m.Called(ctx, patch)
	return args.Error(0)
}

func (m *MockBookingRepository) ListColumn(ctx context.Context) ([]model.Booking, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return []model.Booking{}, args.Error(1)
	}
	return args.Get(0).([]model.Booking), args.Error(1)
}

func (m *MockBookingRepository) FilterBookings(ctx context.Context, filter map[string]interface{}) (map[string][]int, error) {
	args := m.Called(ctx, filter)
	var responses map[string][]int
	if args.Get(0) != nil {
		responses = args.Get(0).(map[string][]int)
	}
	return responses, args.Error(1)
}

func (m *MockBookingRepository) DeleteBooking(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestBookingCreate_Success(t *testing.T) {
	mockRepo := new(MockBookingRepository)

	start := time.Now()
	end := start.Add(48 * time.Hour)
	booking := model.Booking{
		RoomID:     1,
		GuestID:    10,
		Start_date: start,
		End_date:   end,
		Status:     "confirmed",
	}

	mockRepo.On("GettingStatus", mock.Anything, booking.GuestID).Return(false, nil)
	mockRepo.On("ArrivalStatusOfRoom", mock.Anything, booking.RoomID).Return(false, nil)
	mockRepo.On("CreateBooking", mock.Anything, booking).Return(nil)

	uc := NewBookingUsecase(mockRepo)

	err := uc.CreateBooking(context.Background(), booking)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBookingReadByID_Success(t *testing.T) {
	mockRepo := new(MockBookingRepository)

	start := time.Now()
	end := start.Add(24 * time.Hour)
	expected := model.Booking{
		ID:         1,
		RoomID:     2,
		GuestID:    3,
		Start_date: start,
		End_date:   end,
		Status:     "confirmed",
	}

	mockRepo.On("ReadBookingByID", mock.Anything, expected.ID).Return(expected, nil)

	uc := NewBookingUsecase(mockRepo)

	result, err := uc.ReadByIDUsecase(context.Background(), expected.ID)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestBookingPatchByID_Success(t *testing.T) {
	mockRepo := new(MockBookingRepository)

	now := time.Now()
	end := now.Add(72 * time.Hour)
	oldBooking := model.Booking{
		ID:         1,
		RoomID:     1,
		GuestID:    2,
		Start_date: now,
		End_date:   end,
		Status:     "confirmed",
	}

	newStatus := "cancelled"
	patch := dto.BookingPatch{
		ID:     &oldBooking.ID,
		Status: &newStatus,
	}

	mockRepo.On("ReadBookingByID", mock.Anything, oldBooking.ID).Return(oldBooking, nil)
	mockRepo.On("PatchBooking", mock.Anything, mock.Anything).Return(nil)

	uc := NewBookingUsecase(mockRepo)

	err := uc.PatchBookingByID(context.Background(), patch)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "PatchBooking", mock.Anything, mock.MatchedBy(func(p dto.BookingPatch) bool {
		return p.ID != nil && *p.ID == oldBooking.ID && p.Status != nil && *p.Status == newStatus
	}))
}

func TestBookingGetList_Success(t *testing.T) {
	mockRepo := new(MockBookingRepository)

	bookings := []model.Booking{
		{
			ID:         1,
			RoomID:     1,
			GuestID:    1,
			Status:     "confirmed",
			Start_date: time.Now(),
			End_date:   time.Now().Add(24 * time.Hour),
		},
	}

	mockRepo.On("ListColumn", mock.Anything).Return(bookings, nil)

	uc := NewBookingUsecase(mockRepo)

	result, err := uc.GetList(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, bookings[0].ID, result[0].ID)
	mockRepo.AssertExpectations(t)
}

func TestBookingGetFiltered_Success(t *testing.T) {
	mockRepo := new(MockBookingRepository)

	filter := map[string]interface{}{
		"room_id": 1,
	}
	expected := map[string][]int{
		"room_id": {1, 2},
	}

	mockRepo.On("FilterBookings", mock.Anything, filter).Return(expected, nil)

	uc := NewBookingUsecase(mockRepo)

	result, err := uc.GetFilteredBookings(context.Background(), filter)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestBookingRemove_Success(t *testing.T) {
	mockRepo := new(MockBookingRepository)

	mockRepo.On("DeleteBooking", mock.Anything, 1).Return(nil)

	uc := NewBookingUsecase(mockRepo)

	err := uc.RemoveBooking(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
