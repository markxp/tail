// +build linux

package tail

import (
	"context"
	"io"
	"os/exec"
	"sync"
	"time"
)

// cmdTailf wraps build command.
func cmdTailF(ctx context.Context, filepath string, stdout, stderr io.Writer) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "tail", "--follow", filepath)
	cmdIO(cmd, nil, stdout, stderr)

	return cmd
}

func cmdTail(ctx context.Context, filepath string, n string, stdout io.Writer, stderr io.Writer) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "tail", "--lines", n, filepath)
	cmdIO(cmd, nil, stdout, stderr)

	return cmd
}

func cmdRetryTail(ctx context.Context, filepath string, stdout, stderr io.Writer) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "tail", "--follow", "--retry", filepath)
	cmdIO(cmd, nil, stdout, stderr)

	return cmd
}

func retryHelper(bgCtx context.Context, wg *sync.WaitGroup, filepath string, d time.Duration, w io.Writer, errc chan<- error) {
	periodic := int64(d) != int64(0)
	defer wg.Done()

	if !periodic {
		errr, errw := io.Pipe()
		cmd := cmdRetryTail(bgCtx, filepath, w, errw)

		wg.Add(1)
		go func(){
			sc := bufio.NewScanner(errr)
			for sc.Scan(){
				errc<-fmt.Errorf(sc.Text())
			}
		}
		if err := cmd.Run(); err != nil {
			if cerr:= bgCtx.Err(); cerr != nil {
				return
			}
			errc <- err
			return
		}
	}

	ctx, cancel = context.WithTimeout(bgCtx, d)
	defer cancel()
	r, ec := TailF(ctx, filepath)

	wg.Add(1)
	go func() {
		for e := range ec {
			errc <- e
		}
	}()

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		io.WriteString(w, fmt.Sprintln(sc.Text()))
	}
	if sc.Err() != nil {
		errc <- sc.Err()
	}

}
