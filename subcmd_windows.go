// +build windows

package tail

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/markxp/tail/powershell"
)

// cmdTailf wraps build command.
func cmdTailF(ctx context.Context, filepath string, stdout, stderr io.Writer) *exec.Cmd {
	cmd := powershell.PSCommandContext(ctx, "Get-Content", "-LiteralPath", filepath, "-Wait", "-Last", "0")
	cmdIO(cmd, nil, stdout, stderr)
	return cmd
}

func cmdTail(ctx context.Context, filepath string, n string, stdout io.Writer, stderr io.Writer) *exec.Cmd {
	// cmdTail build a Powershell command. It acts as `$ tail -n <n> <filepath>` in linux shell command.
	// The default value of n is "10", which means show the last 10 lines.
	// If n is in the form of "+x", it means show the first x lines.
	//
	// As a dirty fix, all invalid <n> are subsituted with defaul value "10" (compatiable to `$tail`)

	r := regexp.MustCompile(`^([+]?)(0|(?:[1-9][0-9]+))$`)

	var fromHead bool
	var lineNum int64

	m := r.FindStringSubmatch(n)
	if m == nil || len(m) != 3 {
		fromHead = false
		lineNum = 10
	} else {
		lineNum, _ = strconv.ParseInt(m[2], 10, 64)
		fromHead = (m[1] != "")
	}

	line := strconv.FormatInt(lineNum, 10)

	var cmd *exec.Cmd
	var direction string
	if fromHead {
		direction = "-Head"
	} else {
		direction = "-Last"
	}
	cmd = powershell.PSCommandContext(ctx, "Get-Content", "-LiteralPath", filepath, direction, line)
	cmdIO(cmd, nil, stdout, stderr)

	return cmd
}

func retryHelper(bgCtx context.Context, wg *sync.WaitGroup, filepath string, d time.Duration, w io.Writer, errc chan<- error) {
	var ctx context.Context
	var cancel context.CancelFunc
	periodic := int64(d) != int64(0)

	if !periodic {
		ctx = bgCtx
		cancel = func() {}
	} else {
		ctx, cancel = context.WithTimeout(bgCtx, d)
	}

	defer cancel()
	defer wg.Done()

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
