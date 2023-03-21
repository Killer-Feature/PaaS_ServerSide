package os_command_lib

import (
	"errors"
	cl "github.com/Killer-Feature/PaaS_ClientSide/pkg/os_command_lib"
	"regexp"
)

var (
	ErrAssertType = errors.New("error assertion target type")
)

type OSRelease int

const (
	UnknownOS OSRelease = iota
	Ubuntu2204
	Ubuntu2004
)

var (
	OSPrettyNameRegexp = regexp.MustCompile(`PRETTY_NAME="(.*)(".*)`)
	ubuntu2204Regexp   = regexp.MustCompile(`Ubuntu 22\.04\..*`)
	ubuntu2004Regexp   = regexp.MustCompile(`Ubuntu 20\.04\..*`)
)

var (
	OSRegexp = map[OSRelease]*regexp.Regexp{
		Ubuntu2204: ubuntu2204Regexp,
		Ubuntu2004: ubuntu2004Regexp,
	}
)

func GetOSRelease() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command("cat /etc/os-release"),
		Parser:    getOSReleaseParser,
		Condition: cl.Required,
	}
}

func getOSReleaseParser(output []byte, target interface{}) error {
	if len(output) == 0 {
		return nil
	}

	targetOSRelease, ok := target.(*OSRelease)
	if !ok {
		return ErrAssertType
	}

	osPrettyName := OSPrettyNameRegexp.FindSubmatch(output)[1]

	for osRelease, osRegexp := range OSRegexp {
		if osRegexp.Match(osPrettyName) {
			*targetOSRelease = osRelease
			return nil
		}
	}
	*targetOSRelease = UnknownOS
	return nil
}
