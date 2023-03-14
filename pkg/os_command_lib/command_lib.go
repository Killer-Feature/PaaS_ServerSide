package osCommandLib

import "strings"

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
