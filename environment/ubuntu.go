package environment

import (
	"log/slog"
	"os/exec"
	"strings"
)

func init() {
	RegisterSoftware(&Ubuntu{})
}

// Ubuntu represents the Ubuntu OS.
type Ubuntu struct{}

func (u Ubuntu) Name() string {
	return "ubuntu"
}

func (u Ubuntu) Version() (string, bool) {
	cmd := exec.Command("lsb_release", "-r")
	out, err := cmd.Output()
	if err == nil {
		// Release: (spaces) 22.04
		version, ok := strings.CutPrefix(string(out), "Release:")
		if ok {
			return strings.TrimSpace(version), true
		} else {
			slog.Debug("failed to extract Ubuntu version", "output", string(out))
		}
	} else {
		slog.Debug("failed to get Ubuntu version", "error", err)
	}

	return "", false
}
