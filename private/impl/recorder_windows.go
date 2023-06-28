package impl

import (
	"io"
	"reflect"
	"sync"
	"syscall"
	"unsafe"

	"puresound/private/system/winmm"
)

type recorder struct {
	once    sync.Once
	writer  io.Writer
	closed  bool
	hwi     winmm.WaveInHandle
	headers []winmm.WaveHeader
	buffers [][]byte
}

var recorderMap sync.Map //map[uintptr]*recorder
const recorderBufferCount = 10

func NewRecorder(sampleBits, sampleRate, channels, bufferBytes int, writer io.Writer, device ...DeviceHandler) (*recorder, error) {
	dev := uint32(winmm.WAVE_MAPPER)
	if len(device) != 0 {
		dev, _ = device[0].Handle().(uint32)
	}

	f := newWaveFormatEx(sampleBits, sampleRate, channels)
	r := new(recorder)
	r.writer = writer

	err := winmm.WaveInOpen(&r.hwi, dev, f, waveInCallback, uintptr(unsafe.Pointer(r)), winmm.CALLBACK_FUNCTION)
	if err != nil {
		return nil, err
	}

	err = r.initBuffers(bufferBytes)
	if err != nil {
		return nil, err
	}

	recorderMap.Store(r.hwi, r)
	return r, nil
}

func (r *recorder) Close() (err error) {
	if !r.closed {
		r.closed = true
		recorderMap.Delete(r.hwi)
		_ = winmm.WaveInReset(r.hwi)
		r.closeBuffers()
		err = winmm.WaveInClose(r.hwi)
	}
	return
}

func (r *recorder) Start() (err error) {
	err = winmm.WaveInStart(r.hwi)
	return
}

func (r *recorder) Stop() error {
	return winmm.WaveInStop(r.hwi)
}

func (r *recorder) initBuffers(bufferBytes int) error {
	r.buffers = make([][]byte, recorderBufferCount)
	r.headers = make([]winmm.WaveHeader, recorderBufferCount)

	for i := 0; i < recorderBufferCount; i++ {
		r.buffers[i] = make([]byte, bufferBytes)
		bufHdr := *(*reflect.SliceHeader)(unsafe.Pointer(&r.buffers[i]))
		header := &r.headers[i]
		header.Data = unsafe.Pointer(bufHdr.Data)
		header.BufferLength = uint32(bufHdr.Len)
		header.User = unsafe.Pointer(uintptr(i))

		err := winmm.WaveInPrepareHeader(r.hwi, header)
		if err != nil {
			return err
		}

		err = winmm.WaveInAddBuffer(r.hwi, header)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *recorder) closeBuffers() {
	for i := 0; i < recorderBufferCount; i++ {
		header := &r.headers[i]
		_ = winmm.WaveInUnprepareHeader(r.hwi, header)
	}
}

var waveInCallback = syscall.NewCallback(func(hwi winmm.WaveInHandle, msg uint32, inst, p1, p2 unsafe.Pointer) uintptr {
	switch msg {
	case winmm.WIM_OPEN:
		r := (*recorder)(inst)
		r.hwi = hwi

	case winmm.WIM_DATA:
		ir, ok := recorderMap.Load(hwi)
		if !ok {
			//TODO: log it
			break
		}
		r := ir.(*recorder)

		hdr := (*winmm.WaveHeader)(p1)
		buf := r.buffers[int(uintptr(hdr.User))]

		if hdr.BytesRecorded != 0 {
			r.writer.Write(buf)
			err := winmm.WaveInAddBuffer(hwi, hdr)
			if err != nil {
				//TODO: log it
			}
		}
	}

	return 0
})
