package functional

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
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

	err := cmd.Start()
	time.Sleep(time.Second * 3)
	log.Info().Msg(string(b.Bytes()))

	return cmd.Process, b, err
}
