package os_command_lib

import (
	"fmt"
	cl "github.com/Killer-Feature/PaaS_ClientSide/pkg/os_command_lib"
)

const (
	ResourceAlreadyBusy uint16 = 16
)

type Ubuntu2004CommandLib struct{}

func (_ Ubuntu2004CommandLib) RunBinaryCommand(path string) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command("./" + path),
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u Ubuntu2004CommandLib) RunBinaryNohupBackground(command cl.CommandAndParser, output string) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command(fmt.Sprintf("sudo nohup %s > %s 2>&1 &", string(command.Command), output)),
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (_ Ubuntu2004CommandLib) RmFile(path string) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command("rm -f " + path),
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (_ Ubuntu2004CommandLib) Mkdir(path string) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command("mkdir -p " + path),
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (_ Ubuntu2004CommandLib) LoadWebResource(urlPath string, dstPath string) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command(fmt.Sprintf("wget -O %s %s", dstPath, urlPath)),
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (_ Ubuntu2004CommandLib) Chmod777(path string) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command(fmt.Sprintf("chmod 777 %s", path)),
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (_ Ubuntu2004CommandLib) AssertHasProcessListeningPort(port uint16) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command(fmt.Sprintf("if [[ $(sudo lsof -i:%d) ]]; then  exit %d; fi", port, ResourceAlreadyBusy)),
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (_ Ubuntu2004CommandLib) CreateFile(path string) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command("touch " + path),
		Parser:    nil,
		Condition: cl.Required,
	}
}
