package impl

import (
	"github.com/jfreymuth/pulse"
	"github.com/jfreymuth/pulse/proto"
)

const defaultPlayRecordFormat = proto.FormatInt16LE
var (
	defaultRecordSampleRate = pulse.RecordSampleRate(16000)
	defaultRecordBufferSize = pulse.RecordBufferFragmentSize(1600)

	defaultPlaySampleRate = pulse.PlaybackSampleRate(16000)
	defaultPlayBufferSize = pulse.PlaybackRawOption(func(stream *proto.CreatePlaybackStream) {
		stream.BufferTargetLength = 1600
	})
)

func newPlaybackFormatOptions(sampleBits, sampleRate, channels, bufferBytes int, dev ...DeviceHandler) (f byte, opts []pulse.PlaybackOption) {
	switch sampleBits {
	case 8:
		f = proto.FormatUint8
	case 16:
		f = proto.FormatInt16LE
	default:
		f = defaultPlayRecordFormat
	}

	if sampleRate != 0 {
		opts = append(opts, pulse.PlaybackSampleRate(sampleRate))
	} else {
		opts = append(opts, defaultPlaySampleRate)
	}

	switch channels {
	case 1:
		opts = append(opts, pulse.PlaybackMono)
	case 2:
		opts = append(opts, pulse.PlaybackStereo)
	default:
		opts = append(opts, pulse.PlaybackMono)
	}

	if bufferBytes != 0 {
		opts = append(opts, pulse.PlaybackRawOption(func(stream *proto.CreatePlaybackStream) {
			stream.BufferTargetLength = uint32(bufferBytes)
		}))
	} else {
		opts = append(opts, defaultPlayBufferSize)
	}

	if len(dev) != 0 {
		if sink, ok := dev[0].Handle().(*pulse.Sink); ok {
			opts = append(opts, pulse.PlaybackSink(sink))
		}
	}

	return
}

func newRecordFormatOptions(sampleBits, sampleRate, channels, bufferBytes int, dev ...DeviceHandler) (f byte, opts []pulse.RecordOption) {
	switch sampleBits {
	case 8:
		f = proto.FormatUint8
	case 16:
		f = proto.FormatInt16LE
	default:
		f = byte(defaultPlayRecordFormat)
	}

	if sampleRate != 0 {
		opts = append(opts, pulse.RecordSampleRate(int(sampleRate)))
	} else {
		opts = append(opts, defaultRecordSampleRate)
	}

	switch channels {
	case 0, 1:
		opts = append(opts, pulse.RecordMono)
	default:
		opts = append(opts, pulse.RecordStereo)
	}

	if bufferBytes != 0 {
		opts = append(opts, pulse.RecordBufferFragmentSize(uint32(bufferBytes)))
	} else {
		opts = append(opts, defaultRecordBufferSize)
	}

	if len(dev) != 0 {
		if source, ok := dev[0].Handle().(*pulse.Source); ok {
			opts = append(opts, pulse.RecordSource(source))
		}
	}

	return
}
