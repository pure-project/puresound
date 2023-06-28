package impl

import (
	"io"
	"reflect"
	"sync"
	"syscall"
	"unsafe"

	"puresound/private/system/winmm"
)

type player struct {
	bufSize int
	reader  io.Reader

	hwo     winmm.WaveOutHandle
	headers []winmm.WaveHeader
	buffers [][]byte
	pool    chan *winmm.WaveHeader

	state   int

	stop    chan bool
	wg      sync.WaitGroup

	overCount  uint
	writeCount uint
}

var playerMap sync.Map  //map[uintptr]*player
const playerBufferCount = 10
const (
	playerIdle = iota
	playerPlaying
	playerPaused
	playerOver
	playerStopped
	playerClosed
)

func NewPlayer(sampleBits, sampleRate, channels, bufferBytes int, reader io.Reader, device ...DeviceHandler) (*player, error) {
	dev := uint32(winmm.WAVE_MAPPER)
	if len(device) != 0 {
		dev, _ = device[0].Handle().(uint32)
	}

	f := newWaveFormatEx(sampleBits, sampleRate, channels)
	p := new(player)
	p.init(bufferBytes, reader)

	err := winmm.WaveOutOpen(&p.hwo, dev, f, waveOutCallback, uintptr(unsafe.Pointer(p)), winmm.CALLBACK_FUNCTION)
	if err != nil {
		return nil, err
	}

	err = p.initBuffers()
	if err != nil {
		return nil, err
	}

	playerMap.Store(p.hwo, p)
	return p, nil
}

func (p *player) Close() (err error) {
	if p.state != playerClosed {
		p.Stop()
		p.state = playerClosed
		playerMap.Delete(p.hwo)
		p.closeBuffers()
		err = winmm.WaveOutClose(p.hwo)
	}
	return
}

func (p *player) Start() error {
	if p.state != playerIdle && p.state != playerStopped {
		return nil
	}

	p.resetPool()

	err := winmm.WaveOutRestart(p.hwo)
	if err != nil {
		return err
	}

	p.state = playerPlaying
	p.wg.Add(1)
	go p.doData()
	return nil
}

func (p *player) Stop() (err error) {
	if p.state == playerIdle || p.state == playerStopped || p.state == playerClosed {
		return nil
	}

	p.state = playerStopped

	p.stop <- true
	p.wg.Wait()

	err = winmm.WaveOutReset(p.hwo)
	if err != nil {
		return
	}

	return
}

func (p *player) Pause() error {
	if p.state != playerPlaying {
		return nil
	}

	err := winmm.WaveOutPause(p.hwo)
	if err != nil {
		return err
	}

	p.state = playerPaused

	return nil
}

func (p *player) Resume() error {
	if p.state != playerPaused {
		return nil
	}

	err := winmm.WaveOutRestart(p.hwo)
	if err != nil {
		return err
	}

	p.state = playerPlaying

	return nil
}

func (p *player) Playing() bool {
	return p.state == playerPlaying
}

func (p *player) init(bufSize int, reader io.Reader) {
	p.bufSize = bufSize
	p.reader = reader
	p.stop = make(chan bool, 1)
}

func (p *player) initBuffers() error {
	p.buffers = make([][]byte, playerBufferCount)
	p.headers = make([]winmm.WaveHeader, playerBufferCount)
	p.pool    = make(chan *winmm.WaveHeader, playerBufferCount)

	for i := 0; i < playerBufferCount; i++ {
		p.buffers[i] = make([]byte, p.bufSize)
		bufHdr := *(*reflect.SliceHeader)(unsafe.Pointer(&p.buffers[i]))
		header := &p.headers[i]
		header.Data = unsafe.Pointer(bufHdr.Data)
		header.BufferLength = uint32(bufHdr.Len)
		header.User = unsafe.Pointer(uintptr(i))

		err := winmm.WaveOutPrepareHeader(p.hwo, header)
		if err != nil {
			return err
		}

		p.pool <- header
	}

	return nil
}

func (p *player) closeBuffers() {
	for i := 0; i < playerBufferCount; i++ {
		header := &p.headers[i]
		_ = winmm.WaveOutUnprepareHeader(p.hwo, header)
	}
}

func (p *player) resetPool() {
	if p.state != playerStopped {
		return
	}

	for {
		select {
		case <- p.pool:
			continue
		default:
		}
		break
	}

	for i := 0; i < playerBufferCount; i++ {
		header := &p.headers[i]
		p.pool <- header
	}
}

func (p *player) doData() {
	defer func() {
		p.state = playerOver
		p.wg.Done()
	}()

	for {
		select {
		case <- p.stop:
			return

		case hdr := <- p.pool:
			buf := p.buffers[uintptr(hdr.User)]
			n, err := p.reader.Read(buf)
			if err != nil {
				break
			}

			if n != 0 {
				hdr.BufferLength = uint32(n)
				err = winmm.WaveOutWrite(p.hwo, hdr)
				if err != nil {
					return
				}
				p.writeCount++
			}

			continue
		}
		break
	}

	for {
		select {
		case <- p.stop:
			return

		case <- p.pool:
			if p.writeCount == p.overCount {
				return
			}
		}
	}
}

var waveOutCallback = syscall.NewCallback(func(hwo winmm.WaveOutHandle, msg uint32, inst, p1, p2 unsafe.Pointer) uintptr {
	switch msg {
	case winmm.WOM_OPEN:
		p := (*player)(inst)
		p.hwo = hwo

	case winmm.WOM_DONE:
		ip, ok := playerMap.Load(hwo)
		if !ok {
			break
		}

		p := ip.(*player)
		hdr := (*winmm.WaveHeader)(p1)
		p.overCount++
		p.pool <- hdr
	}
	return 0
})