package environment

import (
	"log/slog"
	"os/exec"
	"strings"
)

func init() {
	RegisterSoftware(&Python{})
}

// Python represents the Python SDK.
type Python struct{}

func (p Python) Name() string {
	return "python"
}

func (p Python) Version() (string, bool) {
	cmd := exec.Command("python", "--version")
	out, err := cmd.Output()
	if err == nil {
		// Python 3.11.7
		version, ok := strings.CutPrefix(string(out), "Python ")

		if ok {
			return strings.TrimSpace(version), true
		} else {
			slog.Debug("failed to extract Python version", "output", string(out))
		}
	} else {
		slog.Debug("failed to get Python version", "error", err)
	}

	return "", false
}
