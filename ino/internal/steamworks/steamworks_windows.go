package steamworks

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"os"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

//go:embed steam_api64.dll
var steamAPI64DLL []byte

var steamAPI64DLLHash string

func init() {
	hash := sha256.Sum256(steamAPI64DLL)
	steamAPI64DLLHash = hex.EncodeToString(hash[:])
}

type dll struct {
	d     *windows.LazyDLL
	procs map[string]*windows.LazyProc
}

func (d *dll) call(name string, args ...uintptr) (uintptr, error) {
	if d.procs == nil {
		d.procs = map[string]*windows.LazyProc{}
	}
	if _, ok := d.procs[name]; !ok {
		d.procs[name] = d.d.NewProc(name)
	}
	r, _, err := d.procs[name].Call(args...)
	if err != nil {
		errno, ok := err.(windows.Errno)
		if !ok {
			return r, err
		}
		if errno != 0 {
			return r, err
		}
	}
	return r, nil
}

func loadDLL() (*dll, error) {
	cachedir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(cachedir, "go-inovation")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	fn := filepath.Join(dir, steamAPI64DLLHash+".dll")
	if _, err := os.Stat(fn); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err := os.WriteFile(fn+".tmp", steamAPI64DLL, 0644); err != nil {
			return nil, err
		}
		if err := os.Rename(fn+".tmp", fn); err != nil {
			return nil, err
		}
	}

	return &dll{
		d: windows.NewLazyDLL(fn),
	}, nil
}

var theDLL *dll

func init() {
	dll, err := loadDLL()
	if err != nil {
		panic(err)
	}
	theDLL = dll
}

func RestartAppIfNecessary(appID int) bool {
	v, err := theDLL.call("SteamAPI_RestartAppIfNecessary", uintptr(appID))
	if err != nil {
		panic(err)
	}
	return v != 0
}

func Init() bool {
	v, err := theDLL.call("SteamAPI_Init")
	if err != nil {
		panic(err)
	}
	return v != 0
}

func SteamApps() ISteamApps {
	v, err := theDLL.call("SteamAPI_SteamApps_v008")
	if err != nil {
		panic(err)
	}
	return steamApps(v)
}

type steamApps uintptr

func (s steamApps) GetCurrentGameLanguage() string {
	v, err := theDLL.call("SteamAPI_ISteamApps_GetAvailableGameLanguages", uintptr(s))
	if err != nil {
		panic(err)
	}

	bs := make([]byte, 0, 256)
	for i := int32(0); ; i++ {
		b := *(*byte)(unsafe.Pointer(v))
		v += unsafe.Sizeof(byte(0))
		if b == 0 {
			break
		}
		bs = append(bs, b)
	}
	return string(bs)
}
