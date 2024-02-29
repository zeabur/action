package environment

import (
	"bytes"
	"log/slog"
	"os/exec"
)

func init() {
	RegisterSoftware(&MacOS{})
}

// MacOS represents the macOS.
type MacOS struct{}

func (m MacOS) Name() string {
	return "macos"
}

func (m MacOS) Version() (string, bool) {
	cmd := exec.Command("sw_vers", "-productVersion")
	out, err := cmd.Output()
	if err == nil {
		return string(bytes.TrimSpace(out)), true
	}

	slog.Debug("failed to get macOS version", "error", err)
	return "", false
}
