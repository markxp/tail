package tail

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Tail is equivalent to `$ tail -n lines filepath` linux command.
func Tail(ctx context.Context, filepath string, lines string) ([]string, error) {
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd := cmdTail(ctx, filepath, lines, stdout, stderr)
	err := cmd.Run()
	if err != nil {
		// is the error raised by cancelation?
		// wrap the error for more exact reason.
		cErr := ctx.Err()
		if errors.Is(cErr, context.Canceled) || errors.Is(cErr, context.DeadlineExceeded) {
			err = cErr
		} else {
			err = fmt.Errorf("%w: %s", err, stderr.String())
		}
	}
	ret := strings.Split(stdout.String(), "\n")

	return ret, err
}

// Tailf is equivalent to `$ tail -f` command.
func TailF(ctx context.Context, filepath string) (r io.Reader, errch <-chan error) {
	if ctx.Err() != nil {
		return nil, nil
	}

	pr, pw := io.Pipe()
	errc := make(chan error, 1)
	buf := new(bytes.Buffer)

	cmd := cmdTailF(ctx, filepath, pw, buf)
	go func(ctx context.Context, w *io.PipeWriter) {
		var err error
		defer close(errc)
		defer w.CloseWithError(err)

		if err = cmd.Start(); err != nil {
			errc <- err
			return
		}

		// $ tail -f always got interrupt.
		// If the command is cancel by ctx, replace the exit value as a success.
		err = cmd.Wait()
		if err != nil {
			if cerr := ctx.Err(); cerr != nil && buf.String() == "" {
				return
			}
			errc <- fmt.Errorf("%w: %s", err, buf.String())
		}
	}(ctx, pw)

	return pr, errc
}

// RetryTail acts as `$ tail -F` (equivalent to `$ tail --follow --retry`) command. If the command is finished or failed, it create a new process.
// Every read attempt is recorded and sent as an epoch.
//
// If an error in the epoch is considered not temperoray or not retriable, you should canncel RetryTail with its ctx.
func RetryTailf(ctx context.Context, filepath string, retryInteval time.Duration) (r io.Reader, errch <-chan error) {
	return retryFunc(ctx, filepath, retryInteval)
}

func retryFunc(ctx context.Context, filepath string, retryInteval time.Duration) (r io.Reader, errch <-chan error) {
	if ctx.Err() != nil {
		return nil, nil
	}

	errc := make(chan error, 1)
	tr, tw := io.Pipe()

	go func(root context.Context, filepath string, d time.Duration) {
		wg := &sync.WaitGroup{}

		for {
			select {
			case <-root.Done():
				defer tw.Close()
				defer close(errc)

				wg.Wait()
				return
			default:
				retryHelper(root, wg, filepath, d, tw, errc)
			}
		}
	}(ctx, filepath, retryInteval)
	return tr, errc
}

func cmdIO(cmd *exec.Cmd, stdin io.Reader, stdout, stderr io.Writer) {
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
}
