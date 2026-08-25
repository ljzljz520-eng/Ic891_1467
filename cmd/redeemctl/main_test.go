package main

import (
	"os/exec"
	"testing"
)

func TestCommandHelp(t *testing.T) {
	c := exec.Command("go", "run", ".", "-h")
	if e := c.Run(); e != nil {
		t.Fatal(e)
	}
}
