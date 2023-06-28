package puresound

import "puresound/private/impl"

type RecordDevice interface {
	IsRecordDevice() bool
	Name() string
	Handle() interface{}
}

func DefaultRecord() (RecordDevice, error) {
	return impl.DefaultRecord()
}

func ListRecord() (dev []RecordDevice, err error) {
	err = impl.ListRecord(func(v interface{}) {
		dev = append(dev, v.(RecordDevice))
	})
	return
}
