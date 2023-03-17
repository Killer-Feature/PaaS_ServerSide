package service

import (
	command_lib "KillerFeature/ServerSide/pkg/os_command_lib"
	ubuntu "KillerFeature/ServerSide/pkg/os_command_lib/ubuntu"
	"path/filepath"
)

const (
	HUGGIN_LISTENING_PORT = "8090"
	HUGGIN_DIR            = "huggin"
	HUGGIN_BINARY_NAME    = "HUGGIN"
	NOHUP_OUTPUT_NAME     = "nohup_output"
)

const (
	HUGGIN_BINARY_URL_U2204 = "https://github.com/Killer-Feature/PaaS_ClientSide/releases/download/v0.0.1/PaaS_22.04"
	HUGGIN_BINARY_URL_U2004 = "https://github.com/Killer-Feature/PaaS_ClientSide/releases/download/v0.0.1/PaaS_20.04"
)

func getDeployCommands(commandBook command_lib.CommandLib) []command_lib.Command {
	binaryPath := filepath.Join(HUGGIN_DIR, HUGGIN_BINARY_NAME)
	nohupOutputPath := filepath.Join(HUGGIN_DIR, NOHUP_OUTPUT_NAME)
	switch commandBook.(type) {
	case ubuntu.Ubuntu2204CommandLib:
		{
			return []command_lib.Command{
				commandBook.IfNotEmptyOutputThen(commandBook.GetPIDListeningPort(HUGGIN_LISTENING_PORT), commandBook.ExitWithCode(command_lib.ResourceAlreadyBusy)),
				commandBook.Rmdir(HUGGIN_DIR),
				commandBook.Mkdir(HUGGIN_DIR),
				commandBook.LoadWebResource(HUGGIN_BINARY_URL_U2204, binaryPath),
				commandBook.Chmod777(binaryPath),
				commandBook.CreateFile(nohupOutputPath),
				commandBook.RunBinaryNohupBackground(binaryPath, nohupOutputPath),
			}
		}
	case ubuntu.Ubuntu2004CommandLib:
		{
			return []command_lib.Command{
				commandBook.IfNotEmptyOutputThen(commandBook.GetPIDListeningPort(HUGGIN_LISTENING_PORT), commandBook.ExitWithCode(command_lib.ResourceAlreadyBusy)),
				commandBook.Rmdir(HUGGIN_DIR),
				commandBook.Mkdir(HUGGIN_DIR),
				commandBook.LoadWebResource(HUGGIN_BINARY_URL_U2004, binaryPath),
				commandBook.Chmod777(binaryPath),
				commandBook.CreateFile(nohupOutputPath),
				commandBook.RunBinaryNohupBackground(binaryPath, nohupOutputPath),
			}
		}
	default:
		return nil
	}
}
