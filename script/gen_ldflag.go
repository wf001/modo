package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type VersionInfo struct {
	Version string `yaml:"version"`
	Arch    string `yaml:"arch"`
	Commit  string `yaml:"commit"`
}

func main() {
	data, err := os.ReadFile("cmd/modo/version-info.yaml")
	if err != nil {
		panic(err)
	}

	var v VersionInfo
	if err := yaml.Unmarshal(data, &v); err != nil {
		panic(err)
	}

	// go build用 ldflags を出力
	flags := []string{
		"-X 'main.VERSION=" + v.Version + "'",
		"-X 'main.ARCH=" + v.Arch + "'",
		"-X 'main.COMMIT=" + v.Commit + "'",
	}
	fmt.Println(strings.Join(flags, " "))
}
