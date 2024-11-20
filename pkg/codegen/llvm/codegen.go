package codegen

import (
	"os"
	"os/exec"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	modoTypes "github.com/wf001/modo/pkg/types"
)

type assembler struct {
	node *modoTypes.Node
}

var naryMap = map[modoTypes.NodeKind]func(*ir.Block, value.Value, value.Value) value.Value{
	modoTypes.ND_ADD: func(block *ir.Block, x, y value.Value) value.Value {
		return block.NewAdd(x, y)
	},
	modoTypes.ND_SUB: func(block *ir.Block, x, y value.Value) value.Value {
		return block.NewSub(x, y)
	},
	modoTypes.ND_MUL: func(block *ir.Block, x, y value.Value) value.Value {
		return block.NewMul(x, y)
	},
}

func newI32(s string) *constant.Int {

	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		log.Panic("fail to newInt32: %s", err)
	}
	return constant.NewInt(types.I32, i)
}

func doAsemble(llFile string, asmFile string) {
	out, err := exec.Command("llc", llFile, "-o", asmFile).CombinedOutput()
	if err != nil {
		log.Panic(
			"fail to asemble: %+v",
			map[string]interface{}{"err": err, "out": out, "llFile": llFile, "asmFile": asmFile},
		)
	}
	log.Debug("written asm: %s", asmFile)
}

func gen(mb *ir.Block, node *modoTypes.Node) value.Value {

	if node.IsInteger() {
		return newI32(node.Val)

	} else if node.IsNary() {
		// nary takes arguments more than 2
		child := node.Child
		fst := gen(mb, node.Child)

		child = child.Next
		snd := gen(mb, node.Child.Next)

		nary := naryMap[node.Kind]
		res := nary(mb, fst, snd)

		for child = child.Next; child != nil; child = child.Next {
			fst = res
			snd = gen(mb, child)
			res = nary(mb, fst, snd)
		}
		return res
	}
	return nil
}

func codegen(node *modoTypes.Node) *ir.Module {
	ir := ir.NewModule()
	funcMain := ir.NewFunc(
		"main",
		types.I32,
	)
	llBlock := funcMain.NewBlock("")

	res := gen(llBlock, node)
	llBlock.NewRet(res)
	return ir
}

func Construct(node *modoTypes.Node) *assembler {
	return &assembler{
		node: node,
	}
}

func (a assembler) Assemble(llName string, asmName string) {
	ir := codegen(a.node)

	log.DebugMessage("code generated")
	log.Debug("IR = \n %s\n", ir.String())

	err := os.WriteFile(llName, []byte(ir.String()), 0600)
	if err != nil {
		log.Panic("fail to write ll: %+v", map[string]interface{}{"err": err, "llName": llName})
	}
	log.Debug("written ll: %s", llName)

	doAsemble(llName, asmName)
}
