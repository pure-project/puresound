package puresound

import (
	"io"

	"github.com/pure-project/puresound/private/impl"
)

type Recorder interface {
	Start() error
	Stop() error
	Close() error
}

func NewRecorder(sampleBits, sampleRate, channels, bufferBytes int, writer io.Writer, device ...RecordDevice) (Recorder, error) {
	if len(device) == 0 {
		return impl.NewRecorder(sampleBits, sampleRate, channels, bufferBytes, writer)
	}
	return impl.NewRecorder(sampleBits, sampleRate, channels, bufferBytes, writer, device[0])
}