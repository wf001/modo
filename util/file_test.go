package util

import (
	"fmt"
	"testing"
)

func TestPrepareWorkingFile(t *testing.T) {
	artifactDir, llName, asmName, executableName := PrepareWorkingFile("out", false)
	if llName != fmt.Sprintf("%s/out.ll", artifactDir) {
		t.Errorf("have = %s, want = %s", "out.ll", llName)
	}
	if asmName != fmt.Sprintf("%s/out.s", artifactDir) {
		t.Errorf("have = %s, want = %s", "out.s", asmName)
	}
	if executableName != "out" {
		t.Errorf("have = %s, want = %s", "out", executableName)
	}
}
