package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/sirupsen/logrus"

	"github.com/wf001/modo/pkg/codegen"
	e "github.com/wf001/modo/pkg/error"
	"github.com/wf001/modo/pkg/lexer"
	"github.com/wf001/modo/pkg/log"
	"github.com/wf001/modo/pkg/parser"
	"github.com/wf001/modo/util"
)

const (
	VERSION = "modo version modo0.0.1"
)

var (
	app = kingpin.
		New("modo", "Compiler for the modo programming language.").
		Version(VERSION)
	appVerboseEnabled = app.Flag("verbose", "Show verbose log").Bool()
	appDebugEnabled   = app.Flag("debug", "Show debug log (more detailed than verbose)").Bool()
	appWorkEnabled    = app.Flag("work", "run the program by executable and do not delete temporary work directory when exiting").
				Bool()

	buildCmd       = app.Command("build", "Build an executable")
	buildOutput    = buildCmd.Flag("output", "Write output to <OUTPUT>").Short('o').String()
	buildInputFile = buildCmd.Arg("file", "source file").String()

	runCmd       = app.Command("run", "Build and run a program")
	runExec      = runCmd.Flag("exec", "evaluate <EXEC>").String()
	runInputFile = runCmd.Arg("file", "source file").String()
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
	m["buildOutput"] = *buildOutput
	m["appWorkEnabled"] = *appWorkEnabled
	m["exec"] = *runExec
	m["debug"] = *appDebugEnabled
	m["verbose"] = *appVerboseEnabled
	m["buildInputFile"] = *buildInputFile
	m["runInputFile"] = *runInputFile
	log.Debug("options = %#+v", m)
}

func setLogLevel() {
	if *appVerboseEnabled {
		log.SetLevelInfo()
	}
	if *appDebugEnabled {
		log.SetLevelDebug()
	}
}

type context struct {
	llFlleName            string
	asmFileName           string
	executableFileName    string
	artifactDirectory     string
	isDebug               bool
	storeExecutableInTemp bool
}

func (ctx *context) compile() {
	_, err, errMsg := util.RunCommand("clang", ctx.asmFileName, "-o", ctx.executableFileName)
	if err != nil {
		log.Debug("executableFileNamne: %s", ctx.executableFileName)
		log.Panic("fail to run: err %+v, message %+v", err, errMsg)
	}
	log.Info("successfully wrote executable file: %s", ctx.executableFileName)
}

func (ctx *context) assemble() {
	// TODO: work it?
	out, err, errMsg := util.RunCommand("llc", ctx.llFlleName, "-o", ctx.asmFileName)
	if err != nil {
		log.Debug("llFlleName: %s, asmFileName: %s", ctx.llFlleName, ctx.asmFileName)
		log.Panic("fail to asemble: out %+v, err %+v, message %+v", out, err, errMsg)
	}
	log.Info("successfully wrote assembly file: %s", ctx.asmFileName)
}

func (ctx *context) genFrontend(sourceText string) {
	// string -> Token
	token := lexer.Lex(sourceText)

	// Token -> Node
	node := parser.Parse(token)

	// Node -> write intermediate representation(IR)
	codegen.Construct(node).GenIntermediates(ctx.llFlleName, ctx.asmFileName)
}

func (ctx *context) doRunLLI(sourceText string) {
	ctx.genFrontend(sourceText)

	out, err, errMsg := util.RunCommand("lli", ctx.llFlleName)
	// TODO: it works, but correctly?
	if err != nil {
		log.Panic("%s: fail to run: err %+v, message %+v", e.ERROR_RUNTINME, err, errMsg)
	}
	log.Info("successfully executed: %s", ctx.llFlleName)
	fmt.Println(out)

}

func (ctx *context) doBuild(sourceText string) {
	ctx.genFrontend(sourceText)

	// IR -> write assembly
	ctx.assemble()

	// assembly file -> write executable
	ctx.compile()
}

// HACK: It might be better if the return type matches that of doBuild
func (ctx *context) doRunExecutable(sourceText string) {
	ctx.doBuild(sourceText)

	out, err, errMsg := util.RunCommand(ctx.executableFileName)
	// TODO: it works, but correctly?
	if err != nil {
		log.Error("%s: fail to run: err %+v, message %+v", e.ERROR_RUNTINME, err, errMsg)
	}
	log.Info("successfully executed: %s", ctx.executableFileName)
	fmt.Println(out)
}

func wrapException(fn func()) (err error) {
	// Note: Use `panic` over returning `error` in deeply nested function calls:
	//
	//   - When functions are deeply nested (e.g., main → f → g → h → ...),
	//     returning `error` from every level becomes noisy and repetitive.
	//   - Using `panic` allows us to abort execution immediately from any depth
	//     without having to thread `error` through all function signatures.
	//   - A `panic` can be caught with `recover`, converted back to an `error`,
	//     or handled gracefully depending on debug flags or runtime configuration.
	defer func() {
		if logrus.GetLevel() != logrus.DebugLevel {
			if r := recover(); r != nil {
				log.Error("%s", r)
			}
		}
	}()

	fn()

	return nil
}
func (ctx *context) constructContext() {
	artifactDirectory, llName, asmName, executableName := util.PrepareWorkingFile(
		*buildOutput,
		ctx.storeExecutableInTemp,
	)
	ctx.artifactDirectory = artifactDirectory
	ctx.llFlleName = llName
	ctx.asmFileName = asmName
	ctx.executableFileName = executableName
}

func (ctx *context) runRunCmd() {

	var sourceText string

	if *runExec == "" {
		if runInputFile == nil {
			log.Panic("%s: missing input file", e.ERROR_RUNTINME)
		}
		sourceText = util.ReadFile(runInputFile)

	} else {
		sourceText = *runExec
	}

	if *appWorkEnabled {
		ctx.doRunExecutable(sourceText)

	} else {
		ctx.doRunLLI(sourceText)
	}

}

func (ctx *context) runBuildCmd() {
	var sourceText string
	if buildInputFile == nil {
		log.Panic("%s: missing input file", e.ERROR_RUNTINME)
	}
	sourceText = util.ReadFile(buildInputFile)

	ctx.doBuild(sourceText)

}

func main() {
	_ = wrapException(func() {
		cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

		setLogLevel()
		showOpts(cmd)

		ctx := &context{
			isDebug: *appDebugEnabled,
		}

		defer func() {
			if !*appWorkEnabled {
				err := os.RemoveAll(ctx.artifactDirectory)
				if err != nil {
					log.Error("%s: failed to remove temporary work directory: %v", err)
				} else {
					log.Info("successfully removed temporary work directory(%s)", ctx.artifactDirectory)
				}
			}
		}()

		switch cmd {

		case runCmd.FullCommand():
			ctx.storeExecutableInTemp = true
			ctx.constructContext()

			ctx.runRunCmd()

		case buildCmd.FullCommand():
			ctx.storeExecutableInTemp = false
			ctx.constructContext()

			ctx.runBuildCmd()

		default:
		}
	})
}
