package ecletus

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	EmptyCommit  = strings.Repeat("-", 20)
	EmptyVersion = "0.0.0"
)

type BasicSystemInfo struct {
	ProjectName, Version,
	Commit, BuildDate, GoPath string
}

func (this *BasicSystemInfo) SetAllOrDefault(projectName, version, commit, buildDate, goPath string) {
	if projectName == "" {
		projectName = filepath.Base(os.Args[0])
	}
	if buildDate == "" {
		buildDate = time.Now().Format(time.RFC3339)
	}
	if commit == "" {
		commit = EmptyCommit
	}
	if version == "" {
		version = EmptyVersion
	}
	this.ProjectName, this.Version, this.Commit, this.BuildDate, this.GoPath = projectName,
		version, commit, buildDate, goPath
}
