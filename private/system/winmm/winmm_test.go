package winmm

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"sync"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

func TestRecordBasic(t *testing.T) {
	num := WaveInGetNumDevs()
	t.Logf("dev num: %d", num)

	var i uint32
	for ; i < num; i++ {
		caps := WaveInCaps{}
		err := WaveInGetDevCaps(i, &caps)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%d: %s", i, syscall.UTF16ToString(caps.Pname[:]))
	}

	var handle WaveInHandle
	format := WaveFormatEx{
		FormatTag: WAVE_FORMAT_PCM,
		BitsPerSample: 16,
		SamplesPerSec: 8000,
		BytesPerSec:   16000,
		Channels:      1,
		BlockAlign:    2,
	}

	err := WaveInOpen(&handle, 0, &format, syscall.NewCallback(func(hdl WaveInHandle, msg uint32, inst, p1, p2 uintptr) uintptr {
		switch msg {
		case WIM_OPEN:
			t.Log("wave in open")
		case WIM_CLOSE:
			t.Log("wave in close")
		}
		return 0
	}), 0, CALLBACK_FUNCTION)

	if err != nil {
		t.Fatal(err)
	}

	defer WaveInClose(handle)

	var c uint8
	fmt.Scan(&c)
}

func TestRecord(t *testing.T) {
	file, err := os.OpenFile("output.pcm", os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0666)
	if err != nil {
		t.Fatal("open file err:", err)
		return
	}
	defer file.Close()

	format := WaveFormatEx{
		FormatTag: WAVE_FORMAT_PCM,
		BitsPerSample: 16,
		SamplesPerSec: 8000,
		Channels:      1,

		BytesPerSec:   16000,
		BlockAlign:    2,
	}

	var handle WaveInHandle
	err = WaveInOpen(&handle, WAVE_MAPPER, &format, syscall.NewCallback(func(hdl WaveInHandle, msg uint32, inst, p1, p2 uintptr) uintptr {
		switch msg {
		case WIM_OPEN:
			t.Log("wave in open")
		case WIM_CLOSE:
			t.Log("wave in close")
		case WIM_DATA:
			hdr := (*WaveHeader)(unsafe.Pointer(p1))
			t.Logf("wave in data: user=%d bytes=%d flag=%x", hdr.User, hdr.BytesRecorded, hdr.Flags)
			if hdr.BytesRecorded != 0 {
				_, err := file.Write(*(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
					Data: uintptr(hdr.Data),
					Len:  int(hdr.BytesRecorded),
				})))
				if err != nil {
					t.Log("write file err:", err)
				}
				err = WaveInAddBuffer(hdl, hdr)
				if err != nil {
					t.Log("callback WaveInAddBuffer err:", err)
				}
			}
		}
		return 0
	}), 0, CALLBACK_FUNCTION)

	if err != nil {
		t.Fatal(err)
		return
	}

	defer WaveInClose(handle)

	bufs := make([][]byte, 4)
	hdrs := make([]WaveHeader, 4)
	for i := range hdrs {
		bufs[i] = make([]byte, 1600)
		buf := *(*reflect.SliceHeader)(unsafe.Pointer(&bufs[i]))

		hdr := &hdrs[i]
		hdr.Data = unsafe.Pointer(buf.Data)
		hdr.BufferLength = uint32(buf.Len)
		hdr.User = unsafe.Pointer(uintptr(i))

		err := WaveInPrepareHeader(handle, hdr)
		if err != nil {
			t.Fatal("WaveInPrepareHeader err:", err)
			return
		}
		defer WaveInUnprepareHeader(handle, hdr)

		err = WaveInAddBuffer(handle, hdr)
		if err != nil {
			t.Fatal("WaveInAddBuffer err:", err)
			return
		}
	}

	err = WaveInStart(handle)
	if err != nil {
		t.Fatal("WaveInStart err:", err)
		return
	}

	time.Sleep(10 * time.Second)

	err = WaveInStop(handle)
	if err != nil {
		t.Fatal("WaveInStop err:", err)
		return
	}

	err = WaveInReset(handle)
	if err != nil {
		t.Fatal("WaveInReset err:", err)
		return
	}
}

func TestPlayBasic(t *testing.T) {
	num := WaveOutGetNumDevs()
	t.Logf("dev num: %d", num)

	var i uint32
	for ; i < num; i++ {
		caps := WaveOutCaps{}
		err := WaveOutGetDevCaps(i, &caps)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%d: %s", i, syscall.UTF16ToString(caps.Pname[:]))
	}
}

func TestWaveOutOpen(t *testing.T) {
	file, err := os.OpenFile("D:/output.pcm", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal("open file err:", err)
		return
	}
	defer file.Close()

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

	err = WaveOutOpen(&handle, WAVE_MAPPER, &format, syscall.NewCallback(func(hdl WaveOutHandle, msg uint32, inst, p1, p2 unsafe.Pointer) uintptr {
		switch msg {
		case WOM_OPEN:
			t.Log("wave out open.")

		case WOM_CLOSE:
			t.Log("wave out close.")

		case WOM_DONE:
			hdr := (*WaveHeader)(p1)
			//t.Logf("wave out done: user=%d flag=%x reseting=%t", hdr.User, hdr.Flags, reseting)

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
				//t.Logf("wave out done: user=%d flag=%x next=%d", hdr.User, hdr.Flags, hdr.BufferLength)
				err := WaveOutWrite(hdl, hdr)
				if err != nil {
					t.Log(err)
				}
			}
		}

		return 0
	}), 0, CALLBACK_FUNCTION)

	if err != nil {
		t.Fatal(err)
		return
	}

	if handle == 0 {
		t.Errorf("opened handle is NULL!")
		return
	}

	defer func() {
		err := WaveOutClose(handle)
		if err != nil {
			t.Log(err)
		}
	}()

	err = WaveOutPause(handle)
	if err != nil {
		t.Fatalf("waveOutPause %d err: %v", uintptr(handle), err)
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
			t.Fatal("WaveInPrepareHeader err:", err)
			return
		}
		defer WaveOutUnprepareHeader(handle, hdr)

		err = WaveOutWrite(handle, hdr)
		if err != nil {
			t.Fatal("WaveInAddBuffer err:", err)
			return
		}
	}

	err = WaveOutRestart(handle)
	if err != nil {
		t.Fatal(err)
		return
	}

	time.Sleep(3 * time.Second)

	mtx.Lock()
	reseting = true
	mtx.Unlock()

	err = WaveOutReset(handle)
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestPlay(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	file, err := os.OpenFile("D:/output.pcm", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal("open file err:", err)
		return
	}
	defer file.Close()

	format := WaveFormatEx{
		FormatTag: WAVE_FORMAT_PCM,
		BitsPerSample: 16,
		SamplesPerSec: 8000,
		Channels:      1,

		BytesPerSec:   16000,
		BlockAlign:    2,
	}

	var (
		handle WaveOutHandle
		mtx sync.Mutex
		reseting bool
	)

	err = WaveOutOpen(&handle, WAVE_MAPPER, &format, syscall.NewCallback(func(hdl WaveOutHandle, msg uint32, inst, p1, p2 unsafe.Pointer) uintptr {
		switch msg {
		case WOM_OPEN:
			t.Log("wave out open.")

		case WOM_CLOSE:
			t.Log("wave out close.")

		case WOM_DONE:
			hdr := (*WaveHeader)(p1)
			mtx.Lock()
			defer mtx.Unlock()

			if !reseting {
				n, _ := file.Read(*(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
					Data: uintptr(hdr.Data),
					Len:  1600,
				})))
				hdr.BufferLength = uint32(n)
				t.Logf("wave out done: user=%d flag=%x next=%d", hdr.User, hdr.Flags, hdr.BufferLength)
				err := WaveOutWrite(hdl, hdr)
				if err != nil {
					t.Log(err)
				}

			} else {
				//t.Logf("wave out done: user=%d flag=%x reseting", hdr.User, hdr.Flags)
			}
		}

		return 0
	}), 0, CALLBACK_FUNCTION)

	if err != nil {
		t.Fatal(err)
		return
	}

	//time.Sleep(time.Second)

	if handle == 0 {
		t.Errorf("waveOutOpen return NULL handle.")
		return
	}

	defer WaveOutClose(handle)

	err = WaveOutPause(handle)
	if err != nil {
		t.Fatalf("waveOutPause %d err: %v", uintptr(handle), err)
		return
	}

	bufs := make([][]byte, 4)
	hdrs := make([]WaveHeader, 4)
	for i := range hdrs {
		bufs[i] = make([]byte, 1600)
		buf := *(*reflect.SliceHeader)(unsafe.Pointer(&bufs[i]))

		hdr := &hdrs[i]
		hdr.Data = unsafe.Pointer(buf.Data)
		n, _ := file.Read(bufs[i])
		hdr.BufferLength = uint32(n)
		hdr.User = unsafe.Pointer(uintptr(i))

		err := WaveOutPrepareHeader(handle, hdr)
		if err != nil {
			t.Fatal("WaveInPrepareHeader err:", err)
			return
		}
		defer WaveOutUnprepareHeader(handle, hdr)

		err = WaveOutWrite(handle, hdr)
		if err != nil {
			t.Fatal("WaveInAddBuffer err:", err)
			return
		}
	}

	err = WaveOutRestart(handle)
	if err != nil {
		t.Fatal(err)
		return
	}

	time.Sleep(10 * time.Second)

	mtx.Lock()
	reseting = true
	mtx.Unlock()

	err = WaveOutReset(handle)
	if err != nil {
		t.Fatal(err)
		return
	}

}