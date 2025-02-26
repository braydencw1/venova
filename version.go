package venova

import (
	"os"
	"runtime"
)

type VersionInfo struct {
	Name      string
	Version   string
	Revision  string
	Reference string
	GoVers    string
	BuiltAt   string
	OS        string
	Arch      string
}

var (
	NAME      = GetEntryPoint()
	VERSION   = "main"
	REVISION  = "HEAD"
	REFERENCE = "HEAD"
	GoVers    = runtime.Version()
	BUILT     = "now"
	OS        = runtime.GOOS
	Arch      = runtime.GOARCH
	Version   VersionInfo
)

func GetVersion(customName ...string) VersionInfo {
	// Default
	name := NAME

	if len(customName) > 0 && customName[0] != "" {
		name = customName[0]
	}

	Version = VersionInfo{
		Name:      name,
		Version:   VERSION,
		Revision:  REVISION,
		Reference: REFERENCE,
		GoVers:    GoVers,
		BuiltAt:   BUILT,
		OS:        OS,
		Arch:      Arch,
	}
	return Version
}

func GetEntryPoint() string {
	if entry := os.Getenv("ENTRYPOINT"); entry != "" {
		return entry
	}
	return "Venova"
}
