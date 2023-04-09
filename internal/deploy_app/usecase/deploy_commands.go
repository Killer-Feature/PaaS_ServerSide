package usecase

import (
	cl "github.com/Killer-Feature/PaaS_ClientSide/pkg/os_command_lib"
	cl2 "github.com/Killer-Feature/PaaS_ServerSide/pkg/os_command_lib"
	"path/filepath"
)

const (
	HUGINN_LISTENING_PORT uint16 = 80
	HUGINN_DIR                   = "huginn"
	HUGINN_BINARY_NAME           = "HUGINN"
	NOHUP_OUTPUT_NAME            = "nohup_output"
)

const (
	HUGINN_BINARY_URL_U2204 = "https://github.com/Killer-Feature/PaaS_ClientSide/releases/latest/download/PaaS_22.04"
	HUGINN_BINARY_URL_U2004 = "https://github.com/Killer-Feature/PaaS_ClientSide/releases/latest/download/PaaS_20.04"
)

func getDeployCommands(release cl2.OSRelease) []cl.CommandAndParser {
	binaryPath := filepath.Join(HUGINN_DIR, HUGINN_BINARY_NAME)
	outputPath := filepath.Join(HUGINN_DIR, NOHUP_OUTPUT_NAME)
	switch release {
	case cl2.Ubuntu2204:
		{
			os := cl2.Ubuntu2204CommandLib{}
			return []cl.CommandAndParser{
				os.AssertHasProcessListeningPort(HUGINN_LISTENING_PORT),
				os.RmFile(binaryPath),
				os.Mkdir(HUGINN_DIR),
				os.LoadWebResource(HUGINN_BINARY_URL_U2204, binaryPath),
				os.Chmod777(binaryPath),
				os.CreateFile(outputPath),
				os.RunBinaryNohupBackground(binaryPath, outputPath),
			}
		}

	case cl2.Ubuntu2004:
		{
			os := cl2.Ubuntu2004CommandLib{}
			return []cl.CommandAndParser{
				os.AssertHasProcessListeningPort(HUGINN_LISTENING_PORT),
				os.RmFile(binaryPath),
				os.Mkdir(HUGINN_DIR),
				os.LoadWebResource(HUGINN_BINARY_URL_U2004, binaryPath),
				os.Chmod777(binaryPath),
				os.CreateFile(outputPath),
				os.RunBinaryNohupBackground(binaryPath, outputPath),
			}
		}
	default:
		return nil
	}
}
