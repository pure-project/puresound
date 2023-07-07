package main

import (
	"errors"
	"log"
	"os"
	"reflect"
	"sync"
	"syscall"
	"time"
	"unsafe"

	. "github.com/pure-project/puresound/private/system/winmm"
)

type player struct {
	handle  WaveOutHandle
	buffers [][]byte
	headers []WaveHeader
	mtx     sync.Mutex
	file    *os.File
	reset   bool
}

func (p *player) callback(hdl WaveOutHandle, msg uint32, inst, p1, p2 unsafe.Pointer) uintptr {
	switch msg {
	case WOM_OPEN:
		//handle = hdl   //TODO: this can solve
		log.Println("wave out open.")

	case WOM_CLOSE:
		log.Println("wave out close.")

	case WOM_DONE:
		hdr := (*WaveHeader)(p1)
		//log.Printf("wave out done: user=%d flag=%x reseting=%t", hdr.User, hdr.Flags, reseting)

		p.mtx.Lock()
		defer p.mtx.Unlock()

		if !p.reset {
			n, _ := p.file.Read(*(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
				Data: uintptr(hdr.Data),
				Len:  1600,
				Cap:  1600,
			})))
			hdr.BufferLength = uint32(n)
			//TODO 注释下面这句就行，不注释只有debug能跑，WHY?
			log.Printf("wave out done: user=%d flag=%x next=%d", hdr.User, hdr.Flags, hdr.BufferLength)
			err := WaveOutWrite(hdl, hdr)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return 0
}

func (p *player) Open(file *os.File) error {
	p.file = file

	format := WaveFormatEx{
		FormatTag:     WAVE_FORMAT_PCM,
		BitsPerSample: 16,
		SamplesPerSec: 16000,
		Channels:      1,

		BytesPerSec: 32000,
		BlockAlign:  2,
	}

	err := WaveOutOpen(&p.handle, WAVE_MAPPER, &format, syscall.NewCallback(p.callback), 0, CALLBACK_FUNCTION)
	if err != nil {
		return err
	}

	if p.handle == 0 {
		return errors.New("opened handle is NULL")
	}

	err = WaveOutPause(p.handle)
	if err != nil {
		return err
	}

	p.buffers = make([][]byte, 4)
	p.headers = make([]WaveHeader, 4)

	for i := range p.headers {
		p.buffers[i] = make([]byte, 1600)
		buf := (*reflect.SliceHeader)(unsafe.Pointer(&p.buffers [i]))

		hdr := &p.headers[i]
		hdr.Data = unsafe.Pointer(buf.Data)
		n, _ := file.Read(p.buffers[i])
		hdr.BufferLength = uint32(n)
		hdr.User = unsafe.Pointer(uintptr(i))

		err = WaveOutPrepareHeader(p.handle, hdr)
		if err != nil {
			log.Println("WaveInPrepareHeader err:", err)
			return err
		}

		err = WaveOutWrite(p.handle, hdr)
		if err != nil {
			log.Println("WaveInAddBuffer err:", err)
			return err
		}
	}

	return nil
}

func (p *player) Close() error {
	p.mtx.Lock()
	p.reset = true
	p.mtx.Unlock()

	err := WaveOutReset(p.handle)
	if err != nil {
		return err
	}

	for i := range p.headers {
		hdr := &p.headers[i]
		err = WaveOutUnprepareHeader(p.handle, hdr)
		if err != nil {
			return err
		}
	}

	return WaveOutClose(p.handle)
}

func (p *player) Play() error {
	return WaveOutRestart(p.handle)
}

func play(filename string) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		log.Println("open file err:", err)
		return
	}
	defer file.Close()

	p := new(player)
	err = p.Open(file)
	if err != nil {
		log.Println(err)
		return
	}

	err = p.Play()
	if err != nil {
		log.Println(err)
		return
	}

	time.Sleep(5 * time.Second)

	err = p.Close()
	if err != nil {
		log.Println(err)
	}

	return

	var (
		handle WaveOutHandle
		mtx sync.Mutex
		reseting bool
	)

	format := WaveFormatEx{
		FormatTag: WAVE_FORMAT_PCM,
		BitsPerSample: 16,
		SamplesPerSec: 8000,
		Channels:      1,

		BytesPerSec:   16000,
		BlockAlign:    2,
	}

	err = WaveOutOpen(new(WaveOutHandle), WAVE_MAPPER, &format, syscall.NewCallback(func(hdl WaveOutHandle, msg uint32, inst, p1, p2 unsafe.Pointer) uintptr {
		switch msg {
		case WOM_OPEN:
			handle = hdl   //TODO: this can solve
			log.Println("wave out open.")

		case WOM_CLOSE:
			log.Println("wave out close.")

		case WOM_DONE:
			hdr := (*WaveHeader)(p1)
			//log.Printf("wave out done: user=%d flag=%x reseting=%t", hdr.User, hdr.Flags, reseting)

			mtx.Lock()
			defer mtx.Unlock()

			if !reseting {
				n, _ := file.Read(*(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
					Data: uintptr(hdr.Data),
					Len:  1600,
					Cap:  1600,
				})))
				hdr.BufferLength = uint32(n)
				//TODO 注释下面这句就行，不注释只有debug能跑，WHY?
				log.Printf("wave out done: user=%d flag=%x next=%d", hdr.User, hdr.Flags, hdr.BufferLength)
				err := WaveOutWrite(hdl, hdr)
				if err != nil {
					log.Println(err)
				}
			}
		}

		return 0
	}), 0, CALLBACK_FUNCTION)

	if err != nil {
		log.Println(err)
		return
	}

	if handle == 0 {
		log.Println("opened handle is NULL!")
		return
	}

	defer func() {
		err := WaveOutClose(handle)
		if err != nil {
			log.Println(err)
		}
	}()

	err = WaveOutPause(handle)
	if err != nil {
		log.Printf("waveOutPause %d err: %v", uintptr(handle), err)
		return
	}

	bufs := make([][]byte, 4)
	hdrs := make([]WaveHeader, 4)

	for i := range hdrs {
		bufs[i] = make([]byte, 1600)
		buf := (*reflect.SliceHeader)(unsafe.Pointer(&bufs[i]))

		hdr := &hdrs[i]
		hdr.Data = unsafe.Pointer(buf.Data)
		n, _ := file.Read(bufs[i])
		hdr.BufferLength = uint32(n)
		hdr.User = unsafe.Pointer(uintptr(i))

		err := WaveOutPrepareHeader(handle, hdr)
		if err != nil {
			log.Println("WaveInPrepareHeader err:", err)
			return
		}
		defer WaveOutUnprepareHeader(handle, hdr)

		err = WaveOutWrite(handle, hdr)
		if err != nil {
			log.Println("WaveInAddBuffer err:", err)
			return
		}
	}

	err = WaveOutRestart(handle)
	if err != nil {
		log.Println(err)
		return
	}

	time.Sleep(3 * time.Second)

	mtx.Lock()
	reseting = true
	mtx.Unlock()

	err = WaveOutReset(handle)
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	go play("D:/output.pcm")
	time.Sleep(500 * time.Millisecond)
	play("D:/output.pcm")
	time.Sleep(5 * time.Second)
}