package main

import (
	"errors"
	"os/exec"
)

// Generates a 5 character code.
// These are non-determinisitc and offer no guarantees for
// uniquness on different machines.
func generateCode() (string, error) {
	command := exec.Command("dbus-uuidgen")
	commandOutput, err := command.Output()
	if err != nil {
		return "", errors.New("Couldn't generate code")
	}
	return string(commandOutput)[:5], nil
}
