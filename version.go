package venova

import (
	"log"
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
	NAME      = "venova"
	NAME2     = "venova-stream"
	VERSION   = "main"
	REVISION  = "HEAD"
	REFERENCE = "HEAD"
	GoVers    = runtime.Version()
	BUILT     = "now"
	OS        = runtime.GOOS
	Arch      = runtime.GOARCH
	Version   VersionInfo
)

func GetVersion(customName string) VersionInfo {
	if customName == "venova" {
		Version = VersionInfo{
			Name:      NAME,
			Version:   VERSION,
			Revision:  REVISION,
			Reference: REFERENCE,
			GoVers:    GoVers,
			BuiltAt:   BUILT,
			OS:        OS,
			Arch:      Arch,
		}
	} else if customName == "venova-audio" {
		Version = VersionInfo{
			Name:      NAME2,
			Version:   VERSION,
			Revision:  REVISION,
			Reference: REFERENCE,
			GoVers:    GoVers,
			BuiltAt:   BUILT,
			OS:        OS,
			Arch:      Arch,
		}
	} else {
		log.Fatalf("error gathering version, exiting")
	}
	return Version
}
