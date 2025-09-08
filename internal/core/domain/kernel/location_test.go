package kernel_test

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/pkg/errs"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewLocCorrectParams(t *testing.T) {
	// Arrange

	// Act
	loc, err := kernel.NewLocation(3, 4)

	// Assert
	assert.NotEmpty(t, loc, "new location should not be empty")
	assert.NoError(t, err, "should not be error creating new location")
	assert.True(t, loc.IsValid(), "new location should be valid")
	assert.Equal(t, uint8(3), loc.X(), "location X shall be 3")
	assert.Equal(t, uint8(4), loc.Y(), "location Y shall be 4")
}

func Test_NewLocWithOutOfRangeReturnError(t *testing.T) {
	// Arrange
	tests := map[string]struct {
		x        uint8
		y        uint8
		expected error
	}{
		"wrong_x": {
			x: uint8(15),
			y: uint8(4),
			expected: errs.NewValueIsOutOfRangeError("x", uint8(15), kernel.MinX, kernel.MaxX),
		},
		"wrong_y": {
			x: uint8(3),
			y: uint8(0),
			expected: errs.NewValueIsOutOfRangeError("y", uint8(0), kernel.MinY, kernel.MaxY),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			_, err := kernel.NewLocation(test.x, test.y)

			// Assert
			if err.Error() != test.expected.Error() {
				t.Errorf("expected %v, got %v", test.expected, err)
			}
		})
	}
}

func Test_RandomLocation(t *testing.T) {
	for i := range 11 {
		t.Run(fmt.Sprintf("test RandomLocation %d", i), func(t *testing.T) {
			got, err := kernel.RandomLocation()
			assert.NotEmpty(t, got, "random location shouldn't be empty")
			assert.NoError(t, err, "shouldn't be errors creating random location")
		})
	}
}

func Test_MinLocation(t *testing.T) {
	// Arrange
	// Act
	loc := kernel.MinLocation()

	// Assert
	assert.NotEmpty(t, loc, "new location should not be empty")
	assert.Equal(t, loc.X(), uint8(kernel.MinX))
	assert.Equal(t, loc.Y(), uint8(kernel.MinY))
}

func Test_MaxLocation(t *testing.T) {
	// Arrange
	// Act
	loc := kernel.MaxLocation()

	// Assert
	assert.NotEmpty(t, loc, "new location should not be empty")
	assert.Equal(t, loc.X(), uint8(kernel.MaxX))
	assert.Equal(t, loc.Y(), uint8(kernel.MaxY))
}

func Test_LocationsEqual(t *testing.T) {
	// Arrange
	// Act
	loc1 := kernel.MaxLocation()
	loc2 := kernel.MaxLocation()

	// Assert
	assert.True(t, loc1.Equal(loc2), "two max locations should be equal")
}

func Test_LocationsNotEqual(t *testing.T) {
	// Arrange
	// Act
	loc1 := kernel.MaxLocation()
	loc2 := kernel.MinLocation()

	// Assert
	assert.False(t, loc1.Equal(loc2), "min and max locations should not be equal")
}

func Test_MinMaxLocationsDistance(t *testing.T) {
	// Arrange
	// Act
	loc1 := kernel.MaxLocation()
	loc2 := kernel.MinLocation()
	d1, err1 := loc1.Distance(loc2)
	d2, err2 := loc2.Distance(loc1)
	
	// Assert
	assert.Equal(t, uint8(d1), uint8(18), "distance between min and max should be 18")
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, uint8(d1), uint8(d2), fmt.Sprintf("expected equal, got: %d != %d", d1, d2))
}
