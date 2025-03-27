package venova

import (
	"fmt"
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

func GetVersionInfo(CustomName string) (VersionInfo, error) {
	if CustomName == "" {
		return VersionInfo{}, fmt.Errorf("no name provided, exiting")
	}
	return VersionInfo{
		Name:      CustomName,
		Version:   VERSION,
		Revision:  REVISION,
		Reference: REFERENCE,
		GoVers:    GoVers,
		BuiltAt:   BUILT,
		OS:        OS,
		Arch:      Arch,
	}, nil
}

func (v VersionInfo) String() string {
	return fmt.Sprintf(
		"Name:      %s\nVersion:   %s\nRevision:  %s\nReference: %s\nGo Version: %s\nBuilt At:  %s\nOS:        %s\nArchitecture: %s\n",
		v.Name, v.Version, v.Revision, v.Reference, v.GoVers, v.BuiltAt, v.OS, v.Arch,
	)
}

var (
	NAME      = "venova"
	VERSION   = "main"
	REVISION  = "HEAD"
	REFERENCE = "HEAD"
	GoVers    = runtime.Version()
	BUILT     = "now"
	OS        = runtime.GOOS
	Arch      = runtime.GOARCH
	Version   VersionInfo
)
