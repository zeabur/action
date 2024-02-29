package environment

import (
	"log/slog"
	"os/exec"
	"strings"
)

func init() {
	RegisterSoftware(&Git{})
}

// Git represents the Git utility.
type Git struct{}

func (g Git) Name() string {
	return "git"
}

func (g Git) Version() (string, bool) {
	cmd := exec.Command("git", "--version")
	out, err := cmd.Output()
	if err == nil {
		// git version 2.33.0
		version, ok := strings.CutPrefix(string(out), "git version ")
		if !ok {
			return strings.TrimSpace(version), true
		}

		slog.Debug("failed to extract Git version", "output", string(out))
	} else {
		slog.Debug("failed to get Git version", "error", err)
	}

	return "", false
}
