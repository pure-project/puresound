package impl

import (
	"io"
	"sync"

	"github.com/jfreymuth/pulse"
)

type recorder struct {
	once    sync.Once
	client  *pulse.Client
	stream  *pulse.RecordStream
}

func NewRecorder(sampleBits, sampleRate, channels, bufferBytes int, writer io.Writer, dev ...DeviceHandler) (*recorder, error) {
	f, opts := newRecordFormatOptions(sampleBits, sampleRate, channels, bufferBytes, dev...)

	client, err := pulse.NewClient()
	if err != nil {
		return nil, err
	}

	stream, err := client.NewRecord(pulse.NewWriter(writer, f), opts...)
	if err != nil {
		return nil, err
	}

	return &recorder{
		client: client,
		stream: stream,
	}, nil
}

func (r *recorder) Close() (err error) {
	r.once.Do(func() {
		r.stream.Close()
		r.client.Close()
	})
	return nil
}

func (r *recorder) Start() error {
	r.stream.Start()
	return nil
}

func (r *recorder) Stop() error {
	r.stream.Stop()
	return nil
}
