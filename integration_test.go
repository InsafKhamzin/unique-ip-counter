package main

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIpCounter(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "ips_for_test.txt", "test")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Application failed to run: %v %s %s", err, stdout.String(), stderr.String())
	}
	assert.Contains(t, stdout.String(), "Unique IP count: 18")
}
