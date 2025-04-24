package util

import (
	"bufio"
	"fmt"
	"os"

	"github.com/wf001/modo/pkg/error"
	"github.com/wf001/modo/pkg/log"
)

func PrepareWorkingFile(
	executableName string,
	useExecutable bool,
) (string, string, string, string) {
	artifactDir, err := os.MkdirTemp("", "modo-build-")

	if err != nil {
		log.Panic(
			"%s: fail to make directory: %+v",
			error.ERROR_UNDEFINE,
			map[string]interface{}{"err": err, "artifactDir": artifactDir},
		)
	}
	log.Info("made working directory: %s", artifactDir)

	workingDirPrefix := fmt.Sprintf("%s/out", artifactDir)

	llName := fmt.Sprintf("%s.ll", workingDirPrefix)
	asmName := fmt.Sprintf("%s.s", workingDirPrefix)

	if useExecutable {
		if executableName == "" {
			// TODO: get current dir correctly
			executableName = "./main"
		}
	} else {
		executableName = fmt.Sprintf("%s", workingDirPrefix)
	}

	return artifactDir, llName, asmName, executableName
}

func ReadFile(inputFile *string) string {
	data, err := os.Open(*inputFile)
	defer data.Close()
	// NOTE: Is it better to return err, not to do panic?
	if err != nil {
		log.Panic("file not found: have %s", *inputFile)
	}

	scanner := bufio.NewScanner(data)
	var input_arr = ""
	for scanner.Scan() {
		input_arr += scanner.Text()
	}
	return input_arr
}
