package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/wf001/modo/pkg/codegen"
	"github.com/wf001/modo/pkg/lexer"
	"github.com/wf001/modo/pkg/log"
	"github.com/wf001/modo/pkg/parser"
	"github.com/wf001/modo/util"
)

const (
	VERSION = "0.0.1"
	AUTHOR  = "wf001"
)

var (
	app = kingpin.
		New("modo", "Compiler for the modo programming language.").
		Version(VERSION).
		Author(AUTHOR)
	appVerbose = app.Flag("verbose", "Use verbose log").Bool()
	appDebug   = app.Flag("debug", "Use debug log").Bool()
	appOutput  = app.Flag("output", "Write output to <OUTPUT>").Short('o').String()
	appLLI     = app.Flag("lli", "run with lli").Bool()

	buildCmd = app.Command("build", "Build an executable.")

	runCmd    = app.Command("run", "Build and run an executable.")
	runExec   = runCmd.Flag("exec", "evaluate <EXEC>").String()
	inputFile = runCmd.Arg("file", "source file").String()
)

type IAssebler interface {
	GenFrontend()
}

func GenFrontend(a IAssebler) {
	a.GenFrontend()
}

// HACK: It might be better to move it to different package?
func showOpts(cmd string) {
	m := map[string]interface{}{}
	m["cmd"] = cmd
	m["appOutput"] = *appOutput
	m["appLLI"] = *appLLI
	m["exec"] = *runExec
	m["debug"] = *appDebug
	m["verbose"] = *appVerbose
	m["inputFile"] = *inputFile
	log.Debug("options = %#+v", m)
}

func setLogLevel() {
	if *appVerbose {
		log.SetLevelInfo()
	}
	if *appDebug {
		log.SetLevelDebug()
	}
}

func compile(asmFile string, executableFile string) {
	_, err, errMsg := util.RunCommand("clang", asmFile, "-o", executableFile)
	if err != nil {
		log.Debug("artifactDir: %s", executableFile)
		log.Panic("fail to run: err %+v, message %+v", err, errMsg)
	}
	log.Debug("written executable: %s", executableFile)
}

func assemble(llFile string, asmFile string) {
	// TODO: work it?
	out, err, errMsg := util.RunCommand("llc", llFile, "-o", asmFile)
	if err != nil {
		log.Debug("llFile: %s, asmFile: %s", llFile, asmFile)
		log.Panic("fail to asemble: out %+v, err %+v, message %+v", out, err, errMsg)
	}
	log.Debug("written asm: %s", asmFile)
}

func genFrontend(workingDirPrefix string, evaluatee string) (string, string, string) {
	currentTime := time.Now().Unix()
	// string -> Token
	token := lexer.Lex(evaluatee)

	// Token -> Node
	node := parser.Parse(token)

	llName, asmName, executableName := util.PrepareWorkingFile(workingDirPrefix, currentTime)

	// Node -> write intermediate representation(IR)
	codegen.Construct(node).GenIntermediates(llName, asmName)

	return llName, asmName, executableName
}

func doBuild(workingDirPrefix string, evaluatee string) (error, string) {
	llName, asmName, executableName := genFrontend(workingDirPrefix, evaluatee)

	// IR -> write assembly
	assemble(llName, asmName)

	// assembly file -> write executable
	compile(asmName, executableName)

	return nil, executableName
}

func doRunLLI(workingDirPrefix string, evaluatee string) int {
	llName, _, _ := genFrontend(workingDirPrefix, evaluatee)

	out, err, errMsg := util.RunCommand("lli", llName)
	// TODO: it works, but correctly?
	log.Debug("executed: %s", llName)
	if err != nil {
		log.Debug("artifactDir: %s", llName)
		log.Error("fail to run: err %+v, message %+v", err, errMsg)
		return 1
	}
	fmt.Println(out)

	return 0
}

// HACK: It might be better if the return type matches that of doBuild
func doRunExecutable(workingDirPrefix string, evaluatee string) int {
	err, executableName := doBuild(workingDirPrefix, evaluatee)
	if err != nil {
		log.Panic("fail to run: err %+v, executable %+v", err, executableName)
	}

	out, err, errMsg := util.RunCommand(executableName)
	// TODO: it works, but correctly?
	log.Debug("executed: %s", executableName)
	if err != nil {
		log.Debug("artifactDir: %s", executableName)
		log.Error("fail to run: err %+v, message %+v", err, errMsg)
		return 1
	}
	fmt.Println(out)

	return 0
}

func main() {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	setLogLevel()
	showOpts(cmd)
	switch cmd {

	case runCmd.FullCommand():
		if *runExec == "" {
			if inputFile != nil {
				arg := util.ReadFile(inputFile)
				if *appLLI {
					os.Exit(doRunLLI(*appOutput, arg))
				} else {
					os.Exit(doRunExecutable(*appOutput, arg))
				}
			} else {
				log.Panic("fail to run, input must be specified")
			}
		} else {
			os.Exit(doRunExecutable(*appOutput, *runExec))
		}
	}
}
