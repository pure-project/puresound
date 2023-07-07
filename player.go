package puresound

import (
	"io"

	"github.com/pure-project/puresound/private/impl"
)

type Player interface {
	Start() error
	Stop() error
	Pause() error
	Resume() error
	Playing() bool
	Close() error
}

func NewPlayer(sampleBits, sampleRate, channels, bufferBytes int, reader io.Reader, device ...PlayDevice) (Player, error) {
	if len(device) == 0 {
		return impl.NewPlayer(sampleBits, sampleRate, channels, bufferBytes, reader)
	}
	return impl.NewPlayer(sampleBits, sampleRate, channels, bufferBytes, reader, device[0])
}
