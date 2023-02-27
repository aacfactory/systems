package commons

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

func ExecCommand(ctx context.Context, name string, args ...string) (result []byte, err error) {
	var cancel context.CancelFunc

	_, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
	}
	defer func(cancel context.CancelFunc) {
		if cancel != nil {
			cancel()
		}
	}(cancel)

	cmd := exec.CommandContext(ctx, name, args...)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if execErr := cmd.Start(); execErr != nil {
		err = execErr
		return
	}

	if waitErr := cmd.Wait(); waitErr != nil {
		err = waitErr
		return
	}
	result = buf.Bytes()
	return
}
