package impl

import (
	"io"
	"sync"

	"github.com/jfreymuth/pulse"
)

type player struct {
	once   sync.Once
	client *pulse.Client
	stream *pulse.PlaybackStream
}

func NewPlayer(sampleBits, sampleRate, channels, bufferBytes int, reader io.Reader, dev ...DeviceHandler) (*player, error) {
	f, opts := newPlaybackFormatOptions(sampleBits, sampleRate, channels, bufferBytes, dev...)

	client, err := pulse.NewClient()
	if err != nil {
		return nil, err
	}

	stream, err := client.NewPlayback(pulse.NewReader(reader, f), opts...)
	if err != nil {
		return nil, err
	}

	return &player{
		client: client,
		stream: stream,
	}, nil
}

func (p *player) Close() error {
	p.once.Do(func() {
		p.stream.Close()
		p.client.Close()
	})
	return nil
}

func (p *player) Start() error {
	p.stream.Start()
	return nil
}

func (p *player) Stop() error {
	p.stream.Stop()
	return nil
}

func (p *player) Pause() error {
	p.stream.Pause()
	return nil
}

func (p *player) Resume() error {
	p.stream.Resume()
	return nil
}

func (p *player) Playing() bool {
	return p.stream.Running()
}
