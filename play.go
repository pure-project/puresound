package puresound

import "puresound/private/impl"

type PlayDevice interface {
	IsPlayDevice() bool
	Name() string
	Handle() interface{}
}

func DefaultPlay() (PlayDevice, error) {
	return impl.DefaultPlay()
}

func ListPlay() (dev []PlayDevice, err error) {
	err = impl.ListPlay(func(v interface{}) {
		dev = append(dev, v.(PlayDevice))
	})
	return
}