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
func RunContainer(containerName string, networkName string, hostPorts []int, containerPorts []int, volumeMount string, workingDirectory string, environmentVariables []string, imageName string, startupCommand string) {
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
	for _, environmentVariable := range environmentVariables {
		runCommand += " -e " + environmentVariable
	}
	if volumeMount != "" {
		runCommand += " -v " + volumeMount
	}
	if workingDirectory != "" {
		runCommand += " -w " + workingDirectory
	}
	runCommand += " -it --rm -d " + imageName
	if startupCommand != "" {
		runCommand += " " + startupCommand
	}
	commands.Run(runCommand, "")
}

// CleanContainers - Attempts a `docker stop` and `docker rm`. Fails silently if the container does not exist or is not removed
func CleanContainers(containerNamesOrIds ...string) {
	for _, containerNameOrId := range containerNamesOrIds {
		commands.Run(str.Concat("docker stop ", containerNameOrId), "")
		commands.Run(str.Concat("docker rm ", containerNameOrId), "")
	}
}
