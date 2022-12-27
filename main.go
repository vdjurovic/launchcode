//go:generate goversioninfo -icon=icons/launchcode.ico
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bitshifted/launchcode/config"
)

const (
	updateRetryCode    = 10
	skipUpdateArgument = "--skip-update"
)

func main() {
	// get current application directory
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	appDir := filepath.Dir(exePath)
	fmt.Println(appDir)
	os.Chdir(appDir)

	jvmPath, err := config.FindJvmCommand(appDir)
	if err != nil {
		log.Println("Could not find java command")
	}

	// check if we need to run update
	arguments := os.Args[1:]
	if len(arguments) == 1 && arguments[0] == skipUpdateArgument {
		log.Println("Skipping updte check")
	} else {
		syncroArgs := config.GetSyncroCmdOptions(exePath)
		fmt.Printf("Syncro args: %v\n", syncroArgs)
		syncro := exec.Command(jvmPath, syncroArgs...)
		syncroOut, err := syncro.CombinedOutput()
		if err != nil {
			log.Printf("Failed to run syncro: %s\n", err.Error())
		}
		log.Println(string(syncroOut))
		exitCode := syncro.ProcessState.ExitCode()
		log.Printf("Syncro exit code: %d\n", exitCode)
		if exitCode == updateRetryCode {
			processRetryFiles(exePath)
		}
	}

	if launcherUpdated {
		log.Printf("Launchhing from new executable")
		newCmd := exec.Command(exePath, skipUpdateArgument)
		newCmd.Start()

	} else {
		args := config.GetCmdLineOptions()
		log.Printf("Command line: %v\n", args)

		binary := exec.Command(jvmPath, args...)

		out, execErr := binary.CombinedOutput()
		if execErr != nil {
			log.Printf("Error running Java process: %s\n", execErr.Error())
		}
		log.Println(string(out))
	}

	go cleanup(exePath + ".old")
}
