package impl

import "github.com/jfreymuth/pulse"

type recordDevice struct {
	source *pulse.Source
}

func (d *recordDevice) IsRecordDevice() bool {
	return true
}

func (d *recordDevice) Name() string {
	return d.source.Name()
}

func (d *recordDevice) Handle() interface{} {
	return d.source
}

func DefaultRecord() (dev *recordDevice, err error) {
	client, err := pulse.NewClient()
	if err != nil {
		return
	}
	defer client.Close()

	source, err := client.DefaultSource()
	if err != nil {
		return
	}

	return &recordDevice{ source }, nil
}

func ListRecord(found func(interface{})) (err error) {
	client, err := pulse.NewClient()
	if err != nil {
		return
	}
	defer client.Close()

	sources, err := client.ListSources()
	if err != nil {
		return
	}

	for _, source := range sources {
		found(&recordDevice{source})
	}

	return
}
