package impl

import (
	"syscall"

	"puresound/private/system/winmm"
)

type playDevice struct {
	name string
	dev  uint32
}

func (d *playDevice) IsPlayDevice() bool {
	return true
}

func (d *playDevice) Name() string {
	return d.name
}

func (d *playDevice) Handle() interface{} {
	return d.dev
}

func DefaultPlay() (*playDevice, error) {
	caps := &winmm.WaveOutCaps{}
	err := winmm.WaveOutGetDevCaps(winmm.WAVE_MAPPER, caps)
	if err != nil {
		return nil, err
	}

	return &playDevice{
		name: syscall.UTF16ToString(caps.Pname[:]),
		dev:  winmm.WAVE_MAPPER,
	}, nil
}

func ListPlay(found func(interface{})) (err error) {
	caps := &winmm.WaveOutCaps{}
	num := winmm.WaveOutGetNumDevs()
	for dev := uint32(0); dev < num; dev++ {
		err = winmm.WaveOutGetDevCaps(dev, caps)
		if err != nil {
			return
		}

		found(&playDevice{
			name: syscall.UTF16ToString(caps.Pname[:]),
			dev: dev,
		})
	}
	return
}
