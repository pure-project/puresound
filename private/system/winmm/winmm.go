//+build windows

package winmm

import (
	"strconv"
	"syscall"
	"unsafe"
)

var (
	modwinmm = syscall.NewLazyDLL("winmm.dll")

	procWaveInGetNumDevs      = modwinmm.NewProc("waveInGetNumDevs")
	procWaveInGetDevCaps      = modwinmm.NewProc("waveInGetDevCapsW")
	procWaveInOpen            = modwinmm.NewProc("waveInOpen")
	procWaveInClose           = modwinmm.NewProc("waveInClose")
	procWaveInPrepareHeader   = modwinmm.NewProc("waveInPrepareHeader")
	procWaveInUnprepareHeader = modwinmm.NewProc("waveInUnprepareHeader")
	procWaveInAddBuffer       = modwinmm.NewProc("waveInAddBuffer")
	procWaveInStart           = modwinmm.NewProc("waveInStart")
	procWaveInStop            = modwinmm.NewProc("waveInStop")
	procWaveInReset           = modwinmm.NewProc("waveInReset")
	procWaveInGetErrorText    = modwinmm.NewProc("waveInGetErrorTextW")

	procWaveOutGetNumDevs      = modwinmm.NewProc("waveOutGetNumDevs")
	procWaveOutGetDevCaps      = modwinmm.NewProc("waveOutGetDevCapsW")

	procWaveOutOpen            = modwinmm.NewProc("waveOutOpen")
	procWaveOutClose           = modwinmm.NewProc("waveOutClose")
	procWaveOutPrepareHeader   = modwinmm.NewProc("waveOutPrepareHeader")
	procWaveOutUnprepareHeader = modwinmm.NewProc("waveOutUnprepareHeader")
	procWaveOutWrite           = modwinmm.NewProc("waveOutWrite")
	procWaveOutPause           = modwinmm.NewProc("waveOutPause")
	procWaveOutRestart         = modwinmm.NewProc("waveOutRestart")
	procWaveOutReset           = modwinmm.NewProc("waveOutReset")
	procWaveOutGetErrorText    = modwinmm.NewProc("waveOutGetErrorTextW")
)

const (
	WAVE_MAPPER = 0xFFFFFFFF

	WAVE_FORMAT_PCM = 1

	CALLBACK_FUNCTION = 0x00030000

	WIM_OPEN   = 0x3BE
	WIM_CLOSE  = 0x3BF
	WIM_DATA   = 0x3C0

	WOM_OPEN   = 0x3BB
	WOM_CLOSE  = 0x3BC
	WOM_DONE   = 0x3BD
)

type WaveFormatEx struct {
	FormatTag     uint16
	Channels      uint16
	SamplesPerSec uint32
	BytesPerSec   uint32
	BlockAlign    uint16
	BitsPerSample uint16
	Size          uint16
}

type WaveHeader struct {
	Data          unsafe.Pointer
	BufferLength  uint32
	BytesRecorded uint32
	User          unsafe.Pointer
	Flags         uint32
	Loops         uint32
	Next          *WaveHeader
	Reserved      unsafe.Pointer
}


type WaveInCaps struct {
	Mid       uint16
	Pid       uint16
	DriverVersion uint32
	Pname     [32]uint16
	Formats   uint32
	Channels  uint16
	Reserved1 uint16
}

type WaveInHandle = syscall.Handle

type WaveInResult uint32

func (res WaveInResult) Error() string {
	buf := make([]uint16, 256)
	return waveInGetErrorText(uint32(res), buf)
}


type WaveOutCaps struct {
	Mid       uint16
	Pid       uint16
	DriverVersion uint32
	Pname     [32]uint16
	Formats   uint32
	Channels  uint16
	Reserved1 uint16
	Support   uint32
}

type WaveOutHandle = syscall.Handle

type WaveOutResult uint32

func (res WaveOutResult) Error() string {
	buf := make([]uint16, 256)
	return waveOutGetErrorText(uint32(res), buf)
}


//wave in api

func WaveInGetNumDevs() uint32 {
	num, _, _ := syscall.Syscall(procWaveInGetNumDevs.Addr(), 0, 0, 0, 0)
	return uint32(num)
}

func WaveInGetDevCaps(deviceID uint32, pwic *WaveInCaps) error {
	ret, _, _ := syscall.Syscall(procWaveInGetDevCaps.Addr(), 3, uintptr(deviceID), uintptr(unsafe.Pointer(pwic)), unsafe.Sizeof(WaveInCaps{}))
	if ret != 0 {
		return WaveInResult(ret)
	}
	return nil
}

func WaveInOpen(phwi *WaveInHandle, deviceID uint32, pwfx *WaveFormatEx, callback, instance uintptr, flag uint32) error {
	ret, _, _ := syscall.Syscall6(procWaveInOpen.Addr(), 6, uintptr(unsafe.Pointer(phwi)), uintptr(deviceID), uintptr(unsafe.Pointer(pwfx)), callback, instance, uintptr(flag))
	if ret != 0 {
		return WaveInResult(ret)
	}
	return nil
}

func WaveInClose(hwi WaveInHandle) error {
	ret, _, _ := syscall.Syscall(procWaveInClose.Addr(), 1, uintptr(hwi), 0, 0)
	if ret != 0 {
		return WaveInResult(ret)
	}
	return nil
}

func WaveInPrepareHeader(hwi WaveInHandle, hdr *WaveHeader) error {
	ret, _, _ := syscall.Syscall(procWaveInPrepareHeader.Addr(), 3, uintptr(hwi), uintptr(unsafe.Pointer(hdr)), unsafe.Sizeof(WaveHeader{}))
	if ret != 0 {
		return WaveInResult(ret)
	}
	return nil
}

func WaveInUnprepareHeader(hwi WaveInHandle, hdr *WaveHeader) error {
	ret, _, _ := syscall.Syscall(procWaveInUnprepareHeader.Addr(), 3, uintptr(hwi), uintptr(unsafe.Pointer(hdr)), unsafe.Sizeof(WaveHeader{}))
	if ret != 0 {
		return WaveInResult(ret)
	}
	return nil
}

func WaveInAddBuffer(hwi WaveInHandle, hdr *WaveHeader) error {
	ret, _, _ := syscall.Syscall(procWaveInAddBuffer.Addr(), 3, uintptr(hwi), uintptr(unsafe.Pointer(hdr)), unsafe.Sizeof(WaveHeader{}))
	if ret != 0 {
		return WaveInResult(ret)
	}
	return nil
}

func WaveInStart(hwi WaveInHandle) error {
	ret, _, _ := syscall.Syscall(procWaveInStart.Addr(), 1, uintptr(hwi), 0, 0)
	if ret != 0 {
		return WaveInResult(ret)
	}
	return nil
}

func WaveInStop(hwi WaveInHandle) error {
	ret, _, _ := syscall.Syscall(procWaveInStop.Addr(), 1, uintptr(hwi), 0, 0)
	if ret != 0 {
		return WaveInResult(ret)
	}
	return nil
}

func WaveInReset(hwi WaveInHandle) error {
	ret, _, _ := syscall.Syscall(procWaveInReset.Addr(), 1, uintptr(hwi), 0, 0)
	if ret != 0 {
		return WaveInResult(ret)
	}
	return nil
}

func waveInGetErrorText(ec uint32, buf []uint16) string {
	ret, _, _ := syscall.Syscall(procWaveInGetErrorText.Addr(), 3, uintptr(ec), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if ret != 0 {
		return "winmm wave in err: " + strconv.Itoa(int(ec))
	}
	return "winmm wave in err: " + syscall.UTF16ToString(buf)
}


//wave out api

func WaveOutGetNumDevs() uint32 {
	num, _, _ := syscall.Syscall(procWaveOutGetNumDevs.Addr(), 0, 0, 0, 0)
	return uint32(num)
}

func WaveOutGetDevCaps(deviceID uint32, pwoc *WaveOutCaps) error {
	ret, _, _ := syscall.Syscall(procWaveOutGetDevCaps.Addr(), 3, uintptr(deviceID), uintptr(unsafe.Pointer(pwoc)), unsafe.Sizeof(WaveOutCaps{}))
	if ret != 0 {
		return WaveOutResult(ret)
	}
	return nil
}

func WaveOutOpen(phwo *WaveOutHandle, deviceID uint32, pwfx *WaveFormatEx, callback, instance uintptr, flag uint32) error {
	ret, _, _ := syscall.Syscall6(procWaveOutOpen.Addr(), 6, uintptr(unsafe.Pointer(phwo)), uintptr(deviceID), uintptr(unsafe.Pointer(pwfx)), callback, instance, uintptr(flag))
	if ret != 0 {
		return WaveOutResult(ret)
	}
	return nil
}

func WaveOutClose(hwo WaveOutHandle) error {
	ret, _, _ := syscall.Syscall(procWaveOutClose.Addr(), 1, uintptr(hwo), 0, 0)
	if ret != 0 {
		return WaveOutResult(ret)
	}
	return nil
}

func WaveOutPrepareHeader(hwo WaveOutHandle, hdr *WaveHeader) error {
	ret, _, _ := syscall.Syscall(procWaveOutPrepareHeader.Addr(), 3, uintptr(hwo), uintptr(unsafe.Pointer(hdr)), unsafe.Sizeof(WaveHeader{}))
	if ret != 0 {
		return WaveOutResult(ret)
	}
	return nil
}

func WaveOutUnprepareHeader(hwo WaveOutHandle, hdr *WaveHeader) error {
	ret, _, _ := syscall.Syscall(procWaveOutUnprepareHeader.Addr(), 3, uintptr(hwo), uintptr(unsafe.Pointer(hdr)), unsafe.Sizeof(WaveHeader{}))
	if ret != 0 {
		return WaveOutResult(ret)
	}
	return nil
}

func WaveOutWrite(hwo WaveOutHandle, hdr *WaveHeader) error {
	ret, _, _ := syscall.Syscall(procWaveOutWrite.Addr(), 3, uintptr(hwo), uintptr(unsafe.Pointer(hdr)), unsafe.Sizeof(WaveHeader{}))
	if ret != 0 {
		return WaveOutResult(ret)
	}
	return nil
}

func WaveOutPause(hwo WaveOutHandle) error {
	ret, _, _ := syscall.Syscall(procWaveOutPause.Addr(), 1, uintptr(hwo), 0, 0)
	if ret != 0 {
		return WaveOutResult(ret)
	}
	return nil
}

func WaveOutRestart(hwo WaveOutHandle) error {
	ret, _, _ := syscall.Syscall(procWaveOutRestart.Addr(), 1, uintptr(hwo), 0, 0)
	if ret != 0 {
		return WaveOutResult(ret)
	}
	return nil
}

func WaveOutReset(hwo WaveOutHandle) error {
	ret, _, _ := syscall.Syscall(procWaveOutReset.Addr(), 1, uintptr(hwo), 0, 0)
	if ret != 0 {
		return WaveOutResult(ret)
	}
	return nil
}

func waveOutGetErrorText(ec uint32, buf []uint16) string {
	ret, _, _ := syscall.Syscall(procWaveOutGetErrorText.Addr(), 3, uintptr(ec), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if ret != 0 {
		return "winmm wave out err: " + strconv.Itoa(int(ec))
	}
	return "winmm wave out err: " + syscall.UTF16ToString(buf)
}
