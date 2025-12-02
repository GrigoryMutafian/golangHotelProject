package usecase

import (
	"context"
	"errors"
	"golangHotelProject/internal/delivery/handlers/dto"
	md "golangHotelProject/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRoomRepository struct {
	mock.Mock
}

func (m *MockRoomRepository) CreateRoom(ctx context.Context, room md.Room) error {
	args := m.Called(ctx, room)
	return args.Error(0)
}

func (m *MockRoomRepository) IsNumberExists(ctx context.Context, number int) (bool, error) {
	args := m.Called(ctx, number)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoomRepository) ListRoom(ctx context.Context) ([]md.Room, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return []md.Room{}, args.Error(1)
	}
	return args.Get(0).([]md.Room), args.Error(1)
}

func (m *MockRoomRepository) FilterRoom(ctx context.Context, filter map[string]interface{}) (map[string][]int, error) {
	args := m.Called(ctx, filter)
	var responses map[string][]int
	if args.Get(0) != nil {
		responses = args.Get(0).(map[string][]int)
	}
	return responses, args.Error(1)
}

func (m *MockRoomRepository) PatchRoom(ctx context.Context, id int, p dto.RoomPatch) error {
	args := m.Called(ctx, id, p)
	return args.Error(0)
}

func (m *MockRoomRepository) DeleteRoom(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoomRepository) IsOccupied(ctx context.Context, roomID int) (bool, error) {
	args := m.Called(ctx, roomID)
	return args.Bool(0), args.Error(1)
}

func TestCreateRoom_Success(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	mockRepo.On("IsNumberExists", mock.Anything, 1).Return(false, nil)
	mockRepo.On("CreateRoom", mock.Anything, mock.Anything).Return(nil)

	uc := NewRoomUsecase(mockRepo)

	room := md.Room{
		Number:         1,
		RoomCount:      1,
		IsOccupied:     false,
		Floor:          1,
		SleepingPlaces: 1,
		RoomType:       "Standard",
		NeedCleaning:   false,
	}

	err := uc.AddRoom(context.Background(), room)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateRoom_NumberAlreadyExists(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	mockRepo.On("IsNumberExists", mock.Anything, 1).Return(true, nil)

	uc := NewRoomUsecase(mockRepo)

	room := md.Room{
		Number:         1,
		RoomCount:      1,
		IsOccupied:     false,
		Floor:          1,
		SleepingPlaces: 1,
		RoomType:       "Standard",
		NeedCleaning:   false,
	}

	err := uc.AddRoom(context.Background(), room)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "CreateRoom")
}

func TestCreateRoom_InvalidNumber(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	uc := NewRoomUsecase(mockRepo)

	room := md.Room{
		Number:         -1,
		RoomCount:      1,
		IsOccupied:     false,
		Floor:          1,
		SleepingPlaces: 1,
		RoomType:       "Standard",
		NeedCleaning:   false,
	}

	err := uc.AddRoom(context.Background(), room)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "IsNumberExists")
	mockRepo.AssertNotCalled(t, "CreateRoom")
}

func TestCreateRoom_InvalidRoomCountNumber(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	mockRepo.On("CreateRoom", mock.Anything, mock.Anything).Return(nil)

	uc := NewRoomUsecase(mockRepo)

	room := md.Room{
		Number:         1,
		RoomCount:      0,
		IsOccupied:     false,
		Floor:          1,
		SleepingPlaces: 1,
		RoomType:       "Standard",
		NeedCleaning:   false,
	}

	err := uc.AddRoom(context.Background(), room)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "CreateRoom")
}

func TestCreateRoom_InvalidFloorNumber(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	uc := NewRoomUsecase(mockRepo)

	room := md.Room{
		Number:         1,
		RoomCount:      1,
		IsOccupied:     false,
		Floor:          0,
		SleepingPlaces: 1,
		RoomType:       "Standard",
		NeedCleaning:   false,
	}

	err := uc.AddRoom(context.Background(), room)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "IsNumberExists")
	mockRepo.AssertNotCalled(t, "CreateRoom")
}

func TestCreateRoom_InvalidSleepingPlacesNumber(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	uc := NewRoomUsecase(mockRepo)

	room := md.Room{
		Number:         -1,
		RoomCount:      1,
		IsOccupied:     false,
		Floor:          1,
		SleepingPlaces: 0,
		RoomType:       "Standard",
		NeedCleaning:   false,
	}

	err := uc.AddRoom(context.Background(), room)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "IsNumberExists")
	mockRepo.AssertNotCalled(t, "CreateRoom")
}

func TestCreateRoom_RoomTypeIsEmpty(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	uc := NewRoomUsecase(mockRepo)

	room := md.Room{
		Number:         -1,
		RoomCount:      1,
		IsOccupied:     false,
		Floor:          1,
		SleepingPlaces: 0,
		RoomType:       "",
		NeedCleaning:   false,
	}

	err := uc.AddRoom(context.Background(), room)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "IsNumberExists")
	mockRepo.AssertNotCalled(t, "CreateRoom")

}

func TestCreateRoom_DatabaseErrorWhenCheckingNumberExisting(t *testing.T) {
	mockRepo := new(MockRoomRepository)
	mockRepo.On("IsNumberExists", mock.Anything, 1).Return(false, errors.New("database connection failed"))
	mockRepo.On("CreateRoom", mock.Anything, mock.Anything).Return(nil)

	uc := NewRoomUsecase(mockRepo)

	room := md.Room{
		Number:         1,
		RoomCount:      1,
		IsOccupied:     false,
		Floor:          1,
		SleepingPlaces: 1,
		RoomType:       "Standard",
		NeedCleaning:   false,
	}

	err := uc.AddRoom(context.Background(), room)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "CreateRoom")
}

func TestCreateRoom_DatabaseErrorWhenCreatingRoom(t *testing.T) {
	mockRepo := new(MockRoomRepository)
	mockRepo.On("IsNumberExists", mock.Anything, 1).Return(false, nil)
	mockRepo.On("CreateRoom", mock.Anything, mock.Anything).Return(errors.New("database connection failed"))

	uc := NewRoomUsecase(mockRepo)

	room := md.Room{
		Number:         1,
		RoomCount:      1,
		IsOccupied:     false,
		Floor:          1,
		SleepingPlaces: 1,
		RoomType:       "Standard",
		NeedCleaning:   false,
	}

	err := uc.AddRoom(context.Background(), room)

	assert.Error(t, err)
}

func TestGetList_Success(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	rooms := []md.Room{
		{
			Number:         1,
			RoomCount:      1,
			IsOccupied:     false,
			Floor:          1,
			SleepingPlaces: 2,
			RoomType:       "Standard",
			NeedCleaning:   false,
		},
		{
			Number:         2,
			RoomCount:      1,
			IsOccupied:     true,
			Floor:          2,
			SleepingPlaces: 3,
			RoomType:       "Deluxe",
			NeedCleaning:   false,
		},
	}

	mockRepo.On("ListRoom", mock.Anything).Return(rooms, nil)

	uc := NewRoomUsecase(mockRepo)

	result, err := uc.GetList(context.Background())

	assert.NoError(t, err)

	assert.Len(t, result, 2)

	assert.Equal(t, 1, result[0].Number)
	assert.Equal(t, "Standard", result[0].RoomType)

	assert.Equal(t, 2, result[1].Number)
	assert.Equal(t, "Deluxe", result[1].RoomType)

	mockRepo.AssertExpectations(t)
}

func TestGetList_ListIsClear(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	rooms := []md.Room{}

	mockRepo.On("ListRoom", mock.Anything).Return(rooms, nil)

	uc := NewRoomUsecase(mockRepo)

	result, err := uc.GetList(context.Background())

	assert.Error(t, err)

	assert.Len(t, result, 0)
	mockRepo.AssertExpectations(t)
}

func TestGetList_DatabaseConnectionError(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	rooms := []md.Room{}
	mockRepo.On("ListRoom", mock.Anything).Return(rooms, errors.New("DB connection failed"))

	uc := NewRoomUsecase(mockRepo)

	_, err := uc.GetList(context.Background())

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestGetFilteredRooms_Success(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	filter := map[string]interface{}{
		"floor": 2,
	}

	expectedResponse := map[string][]int{
		"floor": {1, 2, 3},
	}

	mockRepo.On("FilterRoom", mock.Anything, filter).Return(expectedResponse, nil)

	uc := NewRoomUsecase(mockRepo)

	result, err := uc.GetFilteredRooms(context.Background(), filter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, len(result["floor"]))
	assert.Contains(t, result["floor"], 1)
	mockRepo.AssertExpectations(t)
}

func TestGetFilteredRooms_FilterIsEmpty(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	filter := map[string]interface{}{
		"floor": 999,
	}

	expectedResponse := map[string][]int{}

	mockRepo.On("FilterRoom", mock.Anything, filter).Return(expectedResponse, nil)

	uc := NewRoomUsecase(mockRepo)

	result, err := uc.GetFilteredRooms(context.Background(), filter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
	mockRepo.AssertExpectations(t)
}

func TestGetFilteredRooms_DatabaseError(t *testing.T) {
	mockRepo := new(MockRoomRepository)
	filter := map[string]interface{}{
		"floorr": "999",
	}

	mockRepo.On("FilterRoom", mock.Anything, filter).Return(map[string][]int{}, errors.New("column `floorr` doesnt exist"))
	uc := NewRoomUsecase(mockRepo)

	result, err := uc.GetFilteredRooms(context.Background(), filter)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "doesnt exist")
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetFilteredRooms_Multiplyfilters(t *testing.T) {
	mockRepo := new(MockRoomRepository)
	filter := map[string]interface{}{
		"floor":           1,
		"sleeping_places": 1,
	}

	expectedResponse := map[string][]int{
		"floor":           {1, 2, 3},
		"sleeping_places": {1, 2},
	}

	mockRepo.On("FilterRoom", mock.Anything, filter).Return(expectedResponse, nil)
	uc := NewRoomUsecase(mockRepo)
	result, err := uc.GetFilteredRooms(context.Background(), filter)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, 3, len(result["floor"]))
	assert.Equal(t, 2, len(result["sleeping_places"]))
	assert.Contains(t, result["floor"], 1)
	assert.Contains(t, result["sleeping_places"], 1)
	mockRepo.AssertExpectations(t)
}

func TestRemoveRoom_SuccessRemove(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	id := 1
	mockRepo.On("IsOccupied", mock.Anything, id).Return(false, nil)
	mockRepo.On("DeleteRoom", mock.Anything, id).Return(nil)

	uc := NewRoomUsecase(mockRepo)
	err := uc.RemoveRoom(context.Background(), id)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRemoveRoom_InvalidID(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	id := 0

	uc := NewRoomUsecase(mockRepo)
	err := uc.RemoveRoom(context.Background(), id)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID must be more than 0")
	mockRepo.AssertNotCalled(t, "IsOccupied")
	mockRepo.AssertNotCalled(t, "DeleteRoom")
}

func TestRemoveRoom_RoomOccupied(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	id := 1
	mockRepo.On("IsOccupied", mock.Anything, id).Return(true, nil)

	uc := NewRoomUsecase(mockRepo)
	err := uc.RemoveRoom(context.Background(), id)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "room is occupied")
	mockRepo.AssertNotCalled(t, "DeleteRoom")
	mockRepo.AssertExpectations(t)
}

func TestRemoveRoom_RoomOccupiedDatabaseError(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	id := 1
	mockRepo.On("IsOccupied", mock.Anything, id).Return(false, errors.New("db error"))

	uc := NewRoomUsecase(mockRepo)
	err := uc.RemoveRoom(context.Background(), id)
	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "DeleteRoom")
	mockRepo.AssertExpectations(t)
}

func TestRemoveRoom_DeleteRoomDatabaseError(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	id := 1
	mockRepo.On("IsOccupied", mock.Anything, id).Return(false, nil)
	mockRepo.On("DeleteRoom", mock.Anything, id).Return(errors.New("db error"))

	uc := NewRoomUsecase(mockRepo)
	err := uc.RemoveRoom(context.Background(), id)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPatchRoom(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	id := 1
	floor := 2
	sleepingPlaces := 3
	roomType := "Deluxe"
	roomCount := 1
	isOccupied := false
	needCleaning := true

	patch := dto.RoomPatch{
		Floor:          &floor,
		SleepingPlaces: &sleepingPlaces,
		RoomType:       &roomType,
		RoomCount:      &roomCount,
		IsOccupied:     &isOccupied,
		NeedCleaning:   &needCleaning,
	}

	mockRepo.On("PatchRoom", mock.Anything, id, patch).Return(nil)

	uc := NewRoomUsecase(mockRepo)
	err := uc.PatchRoom(context.Background(), id, patch)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPatchRoom_OnlyFloor(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	id := 1
	floor := 2

	patch := dto.RoomPatch{
		Floor: &floor,
	}

	mockRepo.On("PatchRoom", mock.Anything, id, patch).Return(nil)

	uc := NewRoomUsecase(mockRepo)
	err := uc.PatchRoom(context.Background(), id, patch)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPatchRoom_EmptyPatch(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	id := 1
	patch := dto.RoomPatch{}

	uc := NewRoomUsecase(mockRepo)

	err := uc.PatchRoom(context.Background(), id, patch)

	assert.NoError(t, err)
	mockRepo.AssertNotCalled(t, "PatchRoom")
}

func TestPatchRoom_InvalidID(t *testing.T) {
	tests := []struct {
		name string
		id   int
	}{
		{"ID is zero", 0},
		{"ID is negative", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRoomRepository)
			uc := NewRoomUsecase(mockRepo)

			floor := 2
			patch := dto.RoomPatch{Floor: &floor}

			err := uc.PatchRoom(context.Background(), tt.id, patch)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid id")
			mockRepo.AssertNotCalled(t, "PatchRoom")
		})
	}
}

func TestPatchRoom_InvalidFloor(t *testing.T) {
	mockRepo := new(MockRoomRepository)
	floor := -1
	id := 1
	patch := dto.RoomPatch{
		Floor: &floor,
	}

	uc := NewRoomUsecase(mockRepo)
	err := uc.PatchRoom(context.Background(), id, patch)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "PatchRoom")

}

func TestPatchRoom_InvalidSleepingPlaces(t *testing.T) {
	mockRepo := new(MockRoomRepository)
	sleepingPlaces := -1
	id := 1
	patch := dto.RoomPatch{
		SleepingPlaces: &sleepingPlaces,
	}

	uc := NewRoomUsecase(mockRepo)
	err := uc.PatchRoom(context.Background(), id, patch)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "PatchRoom")

}

func TestPatchRoom_InvalidRoomType(t *testing.T) {
	mockRepo := new(MockRoomRepository)
	roomType := "1233"
	id := 1
	patch := dto.RoomPatch{
		RoomType: &roomType,
	}

	uc := NewRoomUsecase(mockRepo)
	err := uc.PatchRoom(context.Background(), id, patch)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "PatchRoom")

}

func TestPatchRoom_InvalidRoomCount(t *testing.T) {
	mockRepo := new(MockRoomRepository)
	roomCount := -1
	id := 1
	patch := dto.RoomPatch{
		RoomCount: &roomCount,
	}

	uc := NewRoomUsecase(mockRepo)
	err := uc.PatchRoom(context.Background(), id, patch)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "PatchRoom")

}

func TestPatchRoom_InvalidDatas(t *testing.T) {
	mockRepo := new(MockRoomRepository)
	needCleaning := true
	isOccupied := true
	id := 1
	patch := dto.RoomPatch{
		NeedCleaning: &needCleaning,
		IsOccupied:   &isOccupied,
	}

	uc := NewRoomUsecase(mockRepo)
	err := uc.PatchRoom(context.Background(), id, patch)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "PatchRoom")

}

func TestPatchRoom_DataBaseError(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	id := 1
	floor := 2

	patch := dto.RoomPatch{
		Floor: &floor,
	}

	mockRepo.On("PatchRoom", mock.Anything, id, patch).Return(errors.New("databese error"))

	uc := NewRoomUsecase(mockRepo)
	err := uc.PatchRoom(context.Background(), id, patch)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPatchRoom_RoomNotFound(t *testing.T) {
	mockRepo := new(MockRoomRepository)

	id := 99999
	floor := 2

	patch := dto.RoomPatch{
		Floor: &floor,
	}

	mockRepo.On("PatchRoom", mock.Anything, id, patch).Return(errors.New("room not found"))

	uc := NewRoomUsecase(mockRepo)
	err := uc.PatchRoom(context.Background(), id, patch)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
