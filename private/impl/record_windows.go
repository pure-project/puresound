package impl

import (
	"syscall"

	"puresound/private/system/winmm"
)

type recordDevice struct {
	name string
	dev  uint32
}

func (d *recordDevice) IsRecordDevice() bool {
	return true
}

func (d *recordDevice) Name() string {
	return d.name
}

func (d *recordDevice) Handle() interface{} {
	return d.dev
}

func DefaultRecord() (*recordDevice, error) {
	caps := &winmm.WaveInCaps{}
	err := winmm.WaveInGetDevCaps(winmm.WAVE_MAPPER, caps)
	if err != nil {
		return nil, err
	}

	return &recordDevice{
		name: syscall.UTF16ToString(caps.Pname[:]),
		dev:  winmm.WAVE_MAPPER,
	}, nil
}

func ListRecord(found func(interface{})) (err error) {
	caps := &winmm.WaveInCaps{}
	num := winmm.WaveInGetNumDevs()
	for dev := uint32(0); dev < num; dev++ {
		err = winmm.WaveInGetDevCaps(dev, caps)
		if err != nil {
			return
		}

		found(&recordDevice{
			name: syscall.UTF16ToString(caps.Pname[:]),
			dev:  dev,
		})
	}
	return
}
