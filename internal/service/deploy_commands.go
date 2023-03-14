package service

import (
	ucase "KillerFeature/ServerSide/internal"
	command_lib "KillerFeature/ServerSide/pkg/os_command_lib"
	ubuntu_commands "KillerFeature/ServerSide/pkg/os_command_lib/ubuntu"
	"path/filepath"
)

const (
	HUGGIN_DIR         = "huggin"
	HUGGIN_BINARY_URL  = "https://github.com/Killer-Feature/PaaS_ClientSide/releases/download/v0.0.1/PaaS_22.04"
	HUGGIN_BINARY_NAME = "HUGGIN"
)

type os string

const (
	ubuntu2204 os = "UBUNTU 22.04"
	ubuntu2004 os = "UBUNTU 20.04"
)

func getDeployCommands(os os) ([]command_lib.Command, error) {
	binaryPath := filepath.Join(HUGGIN_DIR, HUGGIN_BINARY_NAME)
	switch os {
	case ubuntu2204, ubuntu2004:
		return []command_lib.Command{
			ubuntu_commands.RmDir.WithArgs(HUGGIN_DIR),
			ubuntu_commands.Mkdir.WithArgs(HUGGIN_DIR),
			ubuntu_commands.Wget.WithArgs("-O", binaryPath, HUGGIN_BINARY_URL),
			ubuntu_commands.Chmod.WithArgs("777", binaryPath),
			ubuntu_commands.RunBinaryNohupBackground(binaryPath),
		}, nil
	default:
		{
			return nil, ucase.ErrorUnsupportedOS
		}
	}
}
