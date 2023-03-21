package usecase

import (
	cl "github.com/Killer-Feature/PaaS_ClientSide/pkg/os_command_lib"
	cl2 "github.com/Killer-Feature/PaaS_ServerSide/pkg/os_command_lib"
	"path/filepath"
)

const (
	HUGGIN_LISTENING_PORT uint16 = 8090
	HUGGIN_DIR                   = "huggin"
	HUGGIN_BINARY_NAME           = "HUGGIN"
	NOHUP_OUTPUT_NAME            = "nohup_output"
)

const (
	HUGGIN_BINARY_URL_U2204 = "https://github.com/Killer-Feature/PaaS_ClientSide/releases/download/v0.0.1/PaaS_22.04"
	HUGGIN_BINARY_URL_U2004 = "https://github.com/Killer-Feature/PaaS_ClientSide/releases/download/v0.0.1/PaaS_20.04"
)

func getDeployCommands(release cl2.OSRelease) []cl.CommandAndParser {
	binaryPath := filepath.Join(HUGGIN_DIR, HUGGIN_BINARY_NAME)
	outputPath := filepath.Join(HUGGIN_DIR, NOHUP_OUTPUT_NAME)
	switch release {
	case cl2.Ubuntu2204:
		{
			os := cl2.Ubuntu2204CommandLib{}
			return []cl.CommandAndParser{
				os.AssertHasProcessListeningPort(HUGGIN_LISTENING_PORT),
				os.Rmdir(HUGGIN_DIR),
				os.Mkdir(HUGGIN_DIR),
				os.LoadWebResource(HUGGIN_BINARY_URL_U2204, binaryPath),
				os.Chmod777(binaryPath),
				os.CreateFile(outputPath),
				os.RunBinaryNohupBackground(binaryPath, outputPath),
			}
		}

	case cl2.Ubuntu2004:
		{
			os := cl2.Ubuntu2004CommandLib{}
			return []cl.CommandAndParser{
				os.AssertHasProcessListeningPort(HUGGIN_LISTENING_PORT),
				os.Rmdir(HUGGIN_DIR),
				os.Mkdir(HUGGIN_DIR),
				os.LoadWebResource(HUGGIN_BINARY_URL_U2004, binaryPath),
				os.Chmod777(binaryPath),
				os.CreateFile(outputPath),
				os.RunBinaryNohupBackground(binaryPath, outputPath),
			}
		}
	default:
		return nil
	}
}
