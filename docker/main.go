package docker

import (
	"strconv"
  "github.com/PyramidSystemsInc/go/commands"
  "github.com/PyramidSystemsInc/go/errors"
  "github.com/PyramidSystemsInc/go/str"
)

// CreateNetwork - Creates a Docker network
func CreateNetwork(networkName string) {
	_, err := commands.Run(str.Concat("docker network create ", networkName), "")
	errors.QuitIfError(err)
}

// DoesNetworkExist - Checks if a Docker network exists
func DoesNetworkExist(networkName string) bool {
	_, err := commands.Run(str.Concat("docker network inspect ", networkName), "")
	return err == nil
}

// TODO: Fails silently (quickest way to get this done)
// RunContainer - Performs a `docker run` command
func RunContainer(containerName string, networkName string, hostPorts []int, containerPorts []int, imageName string, startupCommand string) {
	runCommand := "docker run"
	runCommand += " --name " + containerName
	if networkName != "" {
		runCommand += " --network " + networkName
	}
	if len(hostPorts) > 0 && len(hostPorts) == len(containerPorts) {
		for i := 0; i < len(hostPorts); i++ {
			runCommand += " -p " + strconv.Itoa(hostPorts[i]) + ":" + strconv.Itoa(containerPorts[i])
		}
	}
	runCommand += " -d " + imageName
	if startupCommand != "" {
		runCommand += " " + startupCommand
	}
	commands.Run(runCommand, "")
}
