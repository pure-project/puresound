package impl

import "github.com/jfreymuth/pulse"

type playDevice struct {
	sink *pulse.Sink
}

func (d *playDevice) IsPlayDevice() bool {
	return true
}

func (d *playDevice) Name() string {
	return d.sink.Name()
}

func (d *playDevice) Handle() interface{} {
	return d.sink
}

func DefaultPlay() (*playDevice, error) {
	client, err := pulse.NewClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	sink, err := client.DefaultSink()
	if err != nil {
		return nil, err
	}

	return &playDevice{ sink }, nil
}

func ListPlay(found func(interface{})) (err error) {
	client, err := pulse.NewClient()
	if err != nil {
		return
	}
	defer client.Close()

	sinks, err := client.ListSinks()
	if err != nil {
		return
	}

	for _, sink := range sinks {
		found(&playDevice{ sink })
	}

	return
}
