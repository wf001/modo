package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/wf001/modo/pkg/error"
	"github.com/wf001/modo/pkg/log"
)

func PrepareWorkingFile(artifactFilePrefix string, currentTime int64) (string, string, string) {
	if artifactFilePrefix == "" {
		generated := "generated"
		artifactDir := fmt.Sprintf("%s/%d", generated, currentTime)
		out, err := exec.Command("mkdir", "-p", artifactDir).CombinedOutput()

		if err != nil {
			log.Panic(
				"%s: fail to make directory: %+v",
				error.ERROR_UNDEFINE,
				map[string]interface{}{"err": err, "out": out, "artifactDir": artifactDir},
			)
		}
		log.Info("make working directory: %s", artifactDir)

		artifactFilePrefix = fmt.Sprintf("%s/out", artifactDir)
	}
	log.Info("complete to persist all of build artifacts in %s", artifactFilePrefix)

	llName := fmt.Sprintf("%s.ll", artifactFilePrefix)
	asmName := fmt.Sprintf("%s.s", artifactFilePrefix)
	executableName := fmt.Sprintf("%s", artifactFilePrefix)

	return llName, asmName, executableName
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
