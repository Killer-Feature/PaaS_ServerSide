package ubuntu

import (
	"fmt"

	cl "github.com/Killer-Feature/PaaS_ServerSide/pkg/os_command_lib"
)

type Ubuntu2004CommandLib struct{}

func (_ Ubuntu2004CommandLib) RunBinaryCommand(path string) cl.Command {
	return cl.Command("./" + path)
}

func (_ Ubuntu2004CommandLib) NohupCommand(command cl.Command, output string) cl.Command {
	return cl.Command(fmt.Sprintf("nohup %s > %s 2>&1 ", command, output))
}

func (_ Ubuntu2004CommandLib) BackgroundCommand(command cl.Command) cl.Command {
	return command + " &"
}

func (u Ubuntu2004CommandLib) RunBinaryNohupBackground(path string, output string) cl.Command {
	return u.NohupCommand(u.BackgroundCommand(u.RunBinaryCommand(path)), output)
}

func (_ Ubuntu2004CommandLib) IfNotEmptyOutputThen(condition cl.Command, statements cl.Command) cl.Command {
	return cl.Command(fmt.Sprintf("if [[ $(%s) ]]; then %s; fi", condition, statements))
}

func (_ Ubuntu2004CommandLib) Rmdir(path string) cl.Command {
	return cl.Command("rmdir -rf" + path)
}

func (_ Ubuntu2004CommandLib) Mkdir(path string) cl.Command {
	return cl.Command("mkdir " + path)
}

func (_ Ubuntu2004CommandLib) LoadWebResource(urlPath string, dstPath string) cl.Command {
	return cl.Command(fmt.Sprintf("wget -O %s %s", dstPath, urlPath))
}

func (_ Ubuntu2004CommandLib) Chmod777(path string) cl.Command {
	return cl.Command(fmt.Sprintf("chmod 777 %s", path))
}

func (_ Ubuntu2004CommandLib) GetPIDListeningPort(port string) cl.Command {
	return cl.Command("lsof -i:" + port)
}

func (_ Ubuntu2004CommandLib) ExitWithCode(code cl.Code) cl.Command {
	return cl.Command("exit " + code)
}

func (_ Ubuntu2004CommandLib) CreateFile(path string) cl.Command {
	return cl.Command("touch " + path)
}
