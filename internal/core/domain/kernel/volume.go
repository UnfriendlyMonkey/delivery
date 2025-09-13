package kernel

import "delivery/internal/pkg/errs"

const MinVolume = 1

type Volume int

func NewVolume(size int) (*Volume, error){
	if size < MinVolume {
		return nil, errs.NewValueIsInvalidError("volume")
	}
	v := Volume(size)
	return &v, nil
}

func (v *Volume) IsValid() bool {
	return int(*v) >= MinVolume
}

func (v *Volume) FitsTo(target *Volume) bool {
	if !target.IsValid() {
		return false
	}
	return *v <= *target
}
