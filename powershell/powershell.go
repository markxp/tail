package powershell

import (
	"context"
	"os/exec"
)

type PSFlags struct {
	UseProfile    bool
	Interactive   bool
	LocalEncoding bool
	ShellPath     string
}

var (
	shellPath         = "powershell.exe"
	flagNoProfile     = "-NoProfile"
	flagNoInteractive = "-NonInteractive"
	cmdUTF8Endcoding  = `[console]::InputEncoding = [Console]::OutputEncoding = [Text.Encoding]::UTF8;`
)

func (fl PSFlags) flags() []string {
	f := []string{}

	if fl.UseProfile == false {
		f = append(f, flagNoProfile)
	}
	if fl.Interactive == false {
		f = append(f, flagNoInteractive)
	}
	if fl.LocalEncoding == false {
		f = append(f, cmdUTF8Endcoding)
	}

	return f
}

var defaultFlag = PSFlags{}

// PSCommand wraps powershell and returns a exec.Cmd
func PSCommand(name string, arg ...string) *exec.Cmd {

	return defaultFlag.Command(name, arg...)
}

// PSCommandContext wraps powershell and returns a exec.Cmd
func PSCommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	return defaultFlag.CommandContext(ctx, name, arg...)
}

// Command wraps powershell with powershell flag settings fl and returns a exec.Cmd
func (fl PSFlags) Command(name string, arg ...string) *exec.Cmd {
	f := fl.flags()
	f = append(f, name)
	f = append(f, arg...)
	if fl.ShellPath != "" {
		return exec.Command(fl.ShellPath, f...)
	}
	return exec.Command(shellPath, f...)
}

// CommandContext wraps powershell with powershell flag settings fl and returns a exec.Cmd
func (fl PSFlags) CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	f := fl.flags()
	f = append(f, name)
	f = append(f, arg...)
	if fl.ShellPath != "" {
		return exec.CommandContext(ctx, shellPath, f...)
	}
	return exec.CommandContext(ctx, shellPath, f...)
}
