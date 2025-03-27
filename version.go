package venova

import (
	"fmt"
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

func GetVersionInfo(names ...string) VersionInfo {
	if len(names) == 0 {
		log.Fatalf("no name provided, exiting")
	}

	// Default to the first name in the variadic input
	name := names[0]
	switch name {
	case "venova":
		name = "venova"
	case "venova-audio-stream":
		name = "venova-audio-stream"
	default:
		log.Fatalf("error gathering version, exiting")
	}

	return VersionInfo{
		Name:      name,
		Version:   "main",
		Revision:  "HEAD",
		Reference: "HEAD",
		GoVers:    runtime.Version(),
		BuiltAt:   "now",
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

func (v VersionInfo) String() string {
	return fmt.Sprintf(
		"Name:      %s\nVersion:   %s\nRevision:  %s\nReference: %s\nGo Version: %s\nBuilt At:  %s\nOS:        %s\nArchitecture: %s\n",
		v.Name, v.Version, v.Revision, v.Reference, v.GoVers, v.BuiltAt, v.OS, v.Arch,
	)
}
