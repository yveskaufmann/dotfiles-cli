package sh

import (
	"os/exec"
)

func IsBinaryOnPath(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
