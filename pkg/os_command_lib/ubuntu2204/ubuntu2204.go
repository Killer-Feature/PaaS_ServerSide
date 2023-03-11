package ubuntu2204

import (
	cl "KillerFeature/ServerSide/pkg/os_command_lib"
)

const (
	Ls        cl.Command = "ls"
	Mkdir     cl.Command = "mkdir"
	Wget      cl.Command = "wget"
	Cd        cl.Command = "cd"
	Chmod     cl.Command = "chmod"
	RunBinary cl.Command = "./"
)
