package environment

import (
	"log/slog"
	"os/exec"
	"strings"
)

func init() {
	RegisterSoftware(&Bash{})
}

// Bash represents the bash shell.
type Bash struct{}

func (b Bash) Name() string {
	return "bash"
}

func (b Bash) Version() (string, bool) {
	cmd := exec.Command("bash", "-c", "echo $BASH_VERSION")
	out, err := cmd.Output()
	if err == nil {
		// the default bash version does not match semver,
		// so we need some editing.
		//
		// 5.1.8(1)-release
		//   -> 5.1.8

		semver, _, ok := strings.Cut(string(out), "(")
		if ok {
			return semver, true
		}

		return strings.TrimSpace(string(out)), true
	}

	slog.Debug("failed to get bash version", "error", err)
	return "", false
}
