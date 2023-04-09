package os_command_lib

import (
	"fmt"
	cl "github.com/Killer-Feature/PaaS_ClientSide/pkg/os_command_lib"
)

type Ubuntu2204CommandLib struct{}

func (_ Ubuntu2204CommandLib) RunBinaryCommand(path string) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command("./" + path),
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u Ubuntu2204CommandLib) RunBinaryNohupBackground(path string, output string) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command(fmt.Sprintf("sudo nohup ./%s > %s 2>&1 &", path, output)),
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (_ Ubuntu2204CommandLib) RmFile(path string) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command("rm -f " + path),
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (_ Ubuntu2204CommandLib) Mkdir(path string) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command("mkdir -p " + path),
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (_ Ubuntu2204CommandLib) LoadWebResource(urlPath string, dstPath string) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command(fmt.Sprintf("wget -O %s %s", dstPath, urlPath)),
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (_ Ubuntu2204CommandLib) Chmod777(path string) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command(fmt.Sprintf("chmod 777 %s", path)),
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (_ Ubuntu2204CommandLib) AssertHasProcessListeningPort(port uint16) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command(fmt.Sprintf("if [[ $(sudo lsof -i:%d) ]]; then  exit %d; fi", port, ResourceAlreadyBusy)),
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (_ Ubuntu2204CommandLib) CreateFile(path string) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command("touch " + path),
		Parser:    nil,
		Condition: cl.Required,
	}
}
