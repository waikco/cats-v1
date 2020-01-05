package functional

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func StartBinary(config string) (*os.Process, bytes.Buffer, error) {
	cmd := exec.Command(
		"../../builds/cats-v1",
		"--config",
		fmt.Sprintf("../settings/%v",
			config))

	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b

	return cmd.Process, b, cmd.Start()
}
