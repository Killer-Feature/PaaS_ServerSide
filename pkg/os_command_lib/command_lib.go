package osCommandLib

import (
	"strings"
)

type Code string

const (
	ResourceAlreadyBusy Code = "16"
)

type Command string

func (c Command) WithArgs(args ...string) Command {
	return c + " " + Command(strings.Join(args, " "))
}

func (c Command) String() string {
	return string(c)
}

func (c Command) WithEnv(env, val, sep string) Command {
	return Command(env + " = " + val + sep + c.String())
}

func (c Command) Pipe(cmds ...Command) Command {
	strCmds := make([]string, 0, 1+len(cmds))
	strCmds = append(strCmds, c.String())
	for i := range cmds {
		strCmds = append(strCmds, cmds[i].String())
	}
	return Command(strings.Join(strCmds, " | "))
}

type CommandLib interface {
	RunBinaryCommand(path string) Command
	NohupCommand(command Command, output string) Command
	BackgroundCommand(command Command) Command
	RunBinaryNohupBackground(path string, output string) Command
	IfNotEmptyOutputThen(condition Command, statements Command) Command
	Rmdir(path string) Command
	Mkdir(path string) Command
	LoadWebResource(urlPath string, dstPath string) Command
	Chmod777(path string) Command
	GetPIDListeningPort(port string) Command
	ExitWithCode(code Code) Command
	CreateFile(path string) Command
}
