package courier_test

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/pkg/errs"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	NameBag = "bag"
	EmptyName = ""
	StandardBag = 10
	VolumeOK = 5
	VolumeBigger = 13
	VolumeWrong = -1
)

func Test_NewStoragePlaceCorrectParams(t *testing.T) {
	// Arrange

	// Act
	sp, err := courier.NewStoragePlace(NameBag, StandardBag)

	// Assert
	assert.NotEmpty(t, sp, "new StoragePlace should not be empty")
	assert.NoError(t, err, "should be no error creating new StoragePlace")
	assert.Equal(t, NameBag, sp.Name(), "name shell be 'bag'")
	assert.Equal(t, StandardBag, sp.TotalVolume(), "total volume shall be 10")
	assert.Nil(t, sp.OrderID(), "order id shall be empty")
}

func Test_NewStoragePlaceWithErrors(t *testing.T) {
	// Arrange
	tests := map[string]struct {
		name string
		volume int
		expected error
	}{
		"wrong_name": {
			name: EmptyName,
			volume: VolumeOK,
			expected: errs.NewValueIsRequiredError("name"),
		},
		"wrong_volume": {
			name: NameBag,
			volume: VolumeWrong,
			expected: errs.NewValueIsInvalidError("totalVolume"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			_, err := courier.NewStoragePlace(test.name, test.volume)

			// Assert
			if err.Error() != test.expected.Error() {
				t.Errorf("expected %v, got %v", test.expected, err)
			}
		})
	}
}

func Test_StoragePlaceCanStoreOK(t *testing.T) {
	// Arrange
	
	// Act
	sp, _ := courier.NewStoragePlace(NameBag, StandardBag)
	ok, err := sp.CanStore(VolumeOK)

	// Assert
	assert.True(t, ok, "Bag (10) must be able to store 5")
	assert.NoError(t, err, "no error expected here")
}

func Test_StoragePlaceCanStoreTooBig(t *testing.T) {
	// Arrange
	
	// Act
	sp, _ := courier.NewStoragePlace(NameBag, StandardBag)
	ok, err := sp.CanStore(VolumeBigger)

	// Assert
	assert.False(t, ok, "Bag (10) must not be able to store 15")
	assert.NoError(t, err, "no error expected here")
}

func Test_StoragePlaceCanStoreWrongVolume(t *testing.T) {
	// Arrange
	
	// Act
	sp, _ := courier.NewStoragePlace(NameBag, StandardBag)
	ok, err := sp.CanStore(VolumeWrong)

	// Assert
	assert.False(t, ok, "Cannot place wrong volume to storage_place")
	assert.Error(t, err, "Should return error here")
}

func Test_EqualNot(t *testing.T) {
	// Arrange

	// Act
	sp1, _ := courier.NewStoragePlace(NameBag, VolumeOK)
	sp2, _ := courier.NewStoragePlace(NameBag, VolumeOK)

	// Assert 
	assert.False(t, sp1.Equal(sp2), "different SPs must not be Equal")
}

func Test_StoreOK(t *testing.T) {
	sp, _ := courier.NewStoragePlace(NameBag, StandardBag)
	err := sp.Store(uuid.New(), VolumeOK)
	assert.NoError(t, err, "No error expected placing 5 to empty bag 10")
}


func Test_StoreOccupied(t *testing.T) {
	sp, _ := courier.NewStoragePlace(NameBag, StandardBag)
	err := sp.Store(uuid.New(), VolumeOK)
	assert.NoError(t, err, "No error expected placing 5 to empty bag 10")
	err = sp.Store(uuid.New(), VolumeOK)
	assert.Error(t, err, "Error expected placing new volume to occupied SP")
	assert.ErrorIs(t, err, courier.ErrStoragePlaceNotEmptyOrLargeEnough, fmt.Sprintf("expected %v, got %v", courier.ErrStoragePlaceNotEmptyOrLargeEnough, err))
}

func Test_ClearOK(t *testing.T) {
	sp, _ := courier.NewStoragePlace(NameBag, StandardBag)
	id := uuid.New()
	_ = sp.Store(id, VolumeOK)
	err := sp.Clear(id)
	assert.NoError(t, err, "No error expected clearing SP with correct ID")
}

func Test_ClearWrongID(t *testing.T) {
	sp, _ := courier.NewStoragePlace(NameBag, StandardBag)
	id := uuid.New()
	_ = sp.Store(id, VolumeOK)
	err := sp.Clear(uuid.New())
	assert.Error(t, err, "Should be error clearing SP with wrong ID")
}
