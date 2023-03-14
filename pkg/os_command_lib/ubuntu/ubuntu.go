package ubuntu

import (
	cl "KillerFeature/ServerSide/pkg/os_command_lib"
)

const (
	RmDir cl.Command = "rm -rf"
	Ls    cl.Command = "ls"
	Mkdir cl.Command = "mkdir"
	Wget  cl.Command = "wget"
	Cd    cl.Command = "cd"
	Chmod cl.Command = "chmod"
)

const (
	nohup      cl.Command = "nohup"
	background cl.Command = "&"
	runBinary  cl.Command = "./"
)

func RunBinaryCommand(path string) cl.Command {
	return runBinary + cl.Command(path)
}

func NohupCommand(command cl.Command) cl.Command {
	return nohup + " " + command
}

func BackgroundCommand(command cl.Command) cl.Command {
	return command + " " + background
}

func RunBinaryNohupBackground(path string) cl.Command {
	return NohupCommand(BackgroundCommand(RunBinaryCommand(path)))
}
