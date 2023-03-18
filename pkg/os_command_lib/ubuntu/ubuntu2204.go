package ubuntu

import (
	"fmt"

	cl "github.com/Killer-Feature/PaaS_ServerSide/pkg/os_command_lib"
)

type Ubuntu2204CommandLib struct{}

func (_ Ubuntu2204CommandLib) RunBinaryCommand(path string) cl.Command {
	return cl.Command("./" + path)
}

func (_ Ubuntu2204CommandLib) NohupCommand(command cl.Command, output string) cl.Command {
	return cl.Command(fmt.Sprintf("nohup %s > %s 2>&1 ", command, output))
}

func (_ Ubuntu2204CommandLib) BackgroundCommand(command cl.Command) cl.Command {
	return command + " &"
}

func (u Ubuntu2204CommandLib) RunBinaryNohupBackground(path string, output string) cl.Command {
	return u.BackgroundCommand(u.NohupCommand(u.RunBinaryCommand(path), output))
}

func (_ Ubuntu2204CommandLib) IfNotEmptyOutputThen(condition cl.Command, statements cl.Command) cl.Command {
	return cl.Command(fmt.Sprintf("if [[ $(%s) ]]; then %s; fi", condition, statements))
}

func (_ Ubuntu2204CommandLib) Rmdir(path string) cl.Command {
	return cl.Command("rm -rf " + path)
}

func (_ Ubuntu2204CommandLib) Mkdir(path string) cl.Command {
	return cl.Command("mkdir " + path)
}

func (_ Ubuntu2204CommandLib) LoadWebResource(urlPath string, dstPath string) cl.Command {
	return cl.Command(fmt.Sprintf("wget -O %s %s", dstPath, urlPath))
}

func (_ Ubuntu2204CommandLib) Chmod777(path string) cl.Command {
	return cl.Command(fmt.Sprintf("chmod 777 %s", path))
}

func (_ Ubuntu2204CommandLib) GetPIDListeningPort(port string) cl.Command {
	return cl.Command("lsof -i:" + port)
}

func (_ Ubuntu2204CommandLib) ExitWithCode(code cl.Code) cl.Command {
	return cl.Command("exit " + code)
}

func (_ Ubuntu2204CommandLib) CreateFile(path string) cl.Command {
	return cl.Command("touch " + path)
}
