// Package kernel for Value Objects common for all app
package kernel

import (
	"fmt"
	"math/rand"
	"time"
	"delivery/internal/pkg/errs"
)

const (
	MinX = 1
	MaxX = 10
	MinY = 1
	MaxY = 10
)

func NewLocation(x, y uint8) (Location, error) {
	if x < MinX || x > MaxX {
		return Location{}, errs.NewValueIsOutOfRangeError("x", x, MinX, MaxX)
	}
	if y < MinY || y > MaxY {
		return Location{}, errs.NewValueIsOutOfRangeError("y", y, MinY, MaxY)
	}
	loc := Location{
		x:     x,
		y:     y,
		valid: true,
	}
	return loc, nil
}

func MinLocation() Location {
	loc, _ := NewLocation(MinX, MinY)
	return loc
}

func MaxLocation() Location {
	loc, _ := NewLocation(MaxX, MaxY)
	return loc
}

func RandomLocation() (Location, error) {
	randomizer := rand.New(rand.NewSource(time.Now().Unix()))
	randx := randomizer.Intn(MaxX - MinX + 1) + MinX
	randy := randomizer.Intn(MaxY - MinY + 1) + MinY
	loc, err := NewLocation(uint8(randx), uint8(randy))
	if err != nil {
		return Location{}, err
	}
	return loc, nil
}

type Location struct {
	x     uint8
	y     uint8
	valid bool
}

func (l Location) X() uint8 {
	return l.x
}

func (l Location) Y() uint8 {
	return l.y
}

func (l Location) Equal(target Location) bool {
	return l == target
}

func (l Location) IsValid() bool {
	return l.valid
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (l Location) Distance(target Location) (uint8, error) {
	if !target.IsValid() {
		return 0, fmt.Errorf("target location: %v is not valid", target)
	}
	distX := abs(int(l.X()) - int(target.X()))
	distY := abs(int(l.Y()) - int(target.Y()))
	return uint8(distX) + uint8(distY), nil
}
