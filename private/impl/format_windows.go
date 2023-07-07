package impl

import (
	"github.com/pure-project/puresound/private/system/winmm"
)

func newWaveFormatEx(sampleBits, sampleRate, channels int) *winmm.WaveFormatEx {
	//default 16k16bit
	f := &winmm.WaveFormatEx{
		FormatTag: winmm.WAVE_FORMAT_PCM,
		SamplesPerSec: 16000,
		BitsPerSample: 16,
		Channels:      1,
		BlockAlign:    2,
		BytesPerSec:   32000,
	}

	if sampleBits != 0 {
		f.BitsPerSample = uint16(sampleBits)
	}

	if sampleRate != 0 {
		f.SamplesPerSec = uint32(sampleRate)
	}

	if channels != 0 {
		f.Channels = uint16(channels)
	}

	f.BlockAlign = f.BitsPerSample * f.Channels / 8
	f.BytesPerSec = uint32(f.BlockAlign) * f.SamplesPerSec

	return f
}
