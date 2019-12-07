package utils

import (
	"os"
	"os/exec"
)

// Compile ...
func Compile() error {
	cmd := exec.Command("make")

	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		Error("Compilation failed")
		return err
	}

	return nil
}
