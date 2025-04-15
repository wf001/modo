package codegen

import (
	"fmt"
	"os"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/codegen/lib"
	"github.com/wf001/modo/pkg/codegen/lib/core"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

type assembler struct {
	program *mTypes.Program
}

type context struct {
	mod      *ir.Module
	function *ir.Func
	block    *ir.Block
	prog     *mTypes.Program
	scope    *mTypes.Node
	argument *mTypes.Node
}

func newBool(s string) *constant.Int {
	i, err := strconv.ParseInt(s, 2, 2)
	if err != nil {
		log.Panic("fail to newBool: %s", err)
	}
	return constant.NewInt(types.I1, i)
}

func newI32(s string) *constant.Int {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		log.Panic("fail to newI32: %s", err)
	}
	return constant.NewInt(types.I32, i)
}

func newStr(ctx *context, n *mTypes.Node) *ir.InstLoad {
	strConst := constant.NewCharArrayFromString(n.Val)
	globalStr := ctx.mod.NewGlobalDef(fmt.Sprintf(".str.%d", len(ctx.mod.Globals)), strConst)
	globalStr.Linkage = enum.LinkagePrivate
	globalStr.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	globalStr.Immutable = true
	globalStr.Align = 1

	strPtr := ctx.block.NewAlloca(types.I8Ptr)
	strGEP := ctx.block.NewGetElementPtr(
		types.NewArray(strConst.Typ.Len, types.I8),
		globalStr,
		newI32("0"),
		newI32("0"),
	)
	ctx.block.NewStore(strGEP, strPtr)
	str := ctx.block.NewLoad(types.I8Ptr, strPtr)
	ctx.prog.GlobalStr = append(ctx.prog.GlobalStr, str)
	return str
}

func newStrHeap(ctx *context, n *mTypes.Node) *ir.InstCall {
	strVal := n.Val
	strLen := len(strVal)

	mallocSize := constant.NewInt(types.I64, int64(strLen))
	dest := ctx.block.NewCall(ctx.prog.BuiltinLibs.Malloc.FuncPtr, mallocSize)

	strConst := constant.NewCharArrayFromString(strVal)
	strConstType := strConst.Typ

	srcAlloca := ctx.block.NewAlloca(strConstType)

	ctx.block.NewStore(strConst, srcAlloca)

	srcPtr := ctx.block.NewGetElementPtr(
		strConstType,
		srcAlloca,
		constant.NewInt(types.I64, 0),
		constant.NewInt(types.I64, 0),
	)

	ctx.block.NewCall(
		ctx.prog.BuiltinLibs.Memcpy.FuncPtr,
		dest,
		srcPtr,
		mallocSize,
		constant.False,
	)
	return dest
}

func newVector(ctx *context, n *mTypes.Node) value.Value {
	elemType, _ := mTypes.GetLLVMType(n.ElemType)
	var arrLength uint64 = 0
	var arr value.Value

	if elemType == types.I8Ptr {
		arrContent := []value.Value{}
		for e := n.Child; e != nil; e = e.Next {
			e.IRValue = ctx.gen(e)
			arrContent = append(arrContent, e.IRValue)
			arrLength++
		}

		arrType := types.NewArray(arrLength, elemType)
		arr = ctx.block.NewAlloca(arrType)

		for i := uint64(0); i < arrLength; i++ {
			elemPtr := ctx.block.NewGetElementPtr(
				arrType,
				arr,
				constant.NewInt(types.I32, 0),
				constant.NewInt(types.I32, int64(i)),
			)
			ctx.block.NewStore(arrContent[i], elemPtr)
		}

	} else {
		arrContent := []constant.Constant{}
		for e := n.Child; e != nil; e = e.Next {
			e.IRValue = ctx.gen(e)
			c, ok := e.IRValue.(constant.Constant)
			if !ok {
				log.Panic("Each array element must be constant.Constant: have %+v", e.IRValue)
			}
			arrContent = append(arrContent, c)
			arrLength++
		}

		arrType := types.NewArray(arrLength, elemType)
		arr = ctx.block.NewAlloca(arrType)
		ctx.block.NewStore(
			constant.NewArray(
				arrType,
				arrContent...,
			),
			arr,
		)
	}
	n.IRValue = arr
	return arr
}

func newVectorGlobal(ctx *context, n *mTypes.Node) value.Value {
	elemType, _ := mTypes.GetLLVMType(n.ElemType)
	var arrLength uint64 = 0
	var arr value.Value

	if elemType == types.I8Ptr {

		strGlobals := []*ir.Global{}
		for i, e := 0, n.Child; e != nil; e, i = e.Next, i+1 {
			strConst := constant.NewCharArrayFromString(e.Val)
			elementStr := ctx.mod.NewGlobalDef(
				fmt.Sprintf(".str.%d", len(ctx.mod.Globals)),
				strConst,
			)
			elementStr.Linkage = enum.LinkagePrivate
			elementStr.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
			elementStr.Immutable = true
			elementStr.Align = 1

			strGlobals = append(strGlobals, elementStr)
			arrLength++
		}

		var strGlobalPtrs []constant.Constant
		for _, g := range strGlobals {
			gep := constant.NewGetElementPtr(
				g.ContentType, // = types.NewArray(len(str), types.I8)
				g,
				constant.NewInt(types.I32, 0),
				constant.NewInt(types.I32, 0),
			)
			strGlobalPtrs = append(strGlobalPtrs, gep)
		}

		arrType := types.NewArray(uint64(len(strGlobalPtrs)), types.NewPointer(types.I8))
		arrConst := constant.NewArray(arrType, strGlobalPtrs...)
		vecGlobal := ctx.mod.NewGlobalDef(
			fmt.Sprintf(".vector.%d", len(ctx.mod.Globals)),
			arrConst,
		)
		vecGlobal.Align = 8
		arr = vecGlobal

	} else {
		arrContent := []constant.Constant{}
		for e := n.Child; e != nil; e = e.Next {
			e.IRValue = ctx.gen(e)
			c, ok := e.IRValue.(constant.Constant)
			if !ok {
				log.Panic("Each array element must be constant.Constant: have %+v", e.IRValue)
			}
			arrContent = append(arrContent, c)
			arrLength++
		}

		arrType := types.NewArray(arrLength, elemType)

		arrConst := constant.NewArray(arrType, arrContent...)
		vecGlobal := ctx.mod.NewGlobalDef(
			fmt.Sprintf(".vector.%d", len(ctx.mod.Globals)),
			arrConst,
		)
		vecGlobal.Align = 8
		arr = vecGlobal

	}
	n.IRValue = arr
	return arr
}

func newVectorHeap(ctx *context, n *mTypes.Node) value.Value {
	log.Debug("heap")
	elemType, _ := mTypes.GetLLVMType(n.ElemType)
	bitWidth := mTypes.GetBitWidth(elemType)
	var arrLength uint64
	var elems []value.Value

	// 各要素を収集
	for e := n.Child; e != nil; e = e.Next {
		e.IRValue = ctx.gen(e)
		elems = append(elems, e.IRValue)
		arrLength++
	}

	// 配列の型: [length x elemType]

	// malloc するサイズ = 要素サイズ * 配列長
	elemSize := bitWidth * arrLength
	totalSize := constant.NewInt(types.I64, int64(elemSize*uint64(arrLength)))
	mallocPtr := ctx.block.NewCall(ctx.prog.BuiltinLibs.Malloc.FuncPtr, totalSize)

	// 要素の型ポインタにキャスト
	vecPtr := ctx.block.NewBitCast(mallocPtr, types.NewPointer(elemType))

	// 各要素を書き込む
	for i, elem := range elems {
		gep := ctx.block.NewGetElementPtr(elemType, vecPtr, constant.NewInt(types.I64, int64(i)))
		ctx.block.NewStore(elem, gep)
	}

	n.IRValue = vecPtr
	return vecPtr
}

func (ctx *context) genVarDeclare(node *mTypes.Node) value.Value {
	if node.Val == "main" {
		// means declaring main function regarded as entrypoint

		fnc := ctx.mod.NewFunc(
			"main",
			types.I32,
		)
		llBlock := fnc.NewBlock("")

		ctx.function = fnc
		ctx.block = llBlock
		res := ctx.gen(node.Child)
		llBlock.NewCall(res)
		llBlock.NewRet(newI32("0"))

	} else {
		// means declaring global variable or function named except main

		// define function return type
		// TODO: in vector need to see Child.ElemType
		varType, ok := mTypes.GetLLVMType(node.Type)

		if !ok {
			varType, ok = mTypes.GetLLVMTypeForVector(node.Child)
		}

		funcName := node.GetFuncName()

		var arg []value.Value
		var argp []*ir.Param

		// define arguments type of function
		for a := node.Child.Args; a != nil; a = a.Next {
			childType, _ := mTypes.GetLLVMType(a.Type)

			arg = append(arg, ir.NewParam(a.Val, childType))
			argp = append(argp, ir.NewParam(a.Val, childType))
		}

		fnc := ctx.mod.NewFunc(
			funcName,
			varType,
			argp...,
		)
		llBlock := fnc.NewBlock("")

		ctx.function = fnc
		ctx.argument = node.Child.Args
		ctx.block = llBlock
		child := ctx.gen(node.Child)
		node.FuncPtr = fnc

		if node.Child.IsKind(mTypes.ND_LAMBDA) {
			lambda := llBlock.NewCall(child, arg...)

			if lambda.Type().Equal(types.Void) {
				llBlock.NewRet(nil)
			} else {
				llBlock.NewRet(lambda)
			}
		} else {
			llBlock.NewRet(child)
		}

	}
	return nil
}

func (ctx *context) genVarReference(node *mTypes.Node) value.Value {
	// PERFORMANCE: too redundant
	// TODO: prohibit same name identifier between global var, binded variable and function argument

	// find in variable which is passed as function argument
	for arg := ctx.argument; arg != nil; arg = arg.Next {
		if arg.Val == node.Val {
			var param value.Value
			for i := 0; i < len(ctx.function.Params); i = i + 1 {
				if ctx.function.Params[i].LocalIdent.LocalName == node.Val {
					param = ctx.function.Params[i]
				}
			}
			return param

		}
	}

	// find in local variable which is declared with let
	for scope := ctx.scope; scope != nil; scope = scope.Next {
		if scope.Val == node.Val {
			if scope.Child.IsScalar() {
				return scope.VarPtr

			} else if scope.Child.IsType(mTypes.TY_VECTOR) {
				return scope.VarPtr

			} else {
				log.Panic("unresolved NodeType: have %+v", node)
			}
		}
	}

	// find in global variable which is declared with def
	for declare := ctx.prog.Declares; declare != nil; declare = declare.Next {
		if declare.Child.Val == node.Val {
			return ctx.block.NewCall(declare.Child.FuncPtr)
		}
	}

	log.Panic("unresolved symbol: '%s'", node.Val)

	return nil
}

func (ctx *context) genLambda(node *mTypes.Node) value.Value {
	isParentMain := ctx.function.GlobalName == "main"
	unnamedFuncName := node.GetUnnamedFuncName()
	fnEntryBlockName := "fn.entry"

	if isParentMain {
		funcFn := ctx.mod.NewFunc(
			unnamedFuncName,
			types.Void,
			ctx.function.Params...,
		)
		llBlock := funcFn.NewBlock(node.GetBlockName(fnEntryBlockName))

		ctx.function = funcFn
		ctx.block = llBlock

		ctx.gen(node.Child)

		ctx.block.NewRet(nil)

		return funcFn

	} else {
		funcFn := ctx.mod.NewFunc(
			unnamedFuncName,
			ctx.function.Sig.RetType,
			ctx.function.Params...,
		)
		llBlock := funcFn.NewBlock(node.GetBlockName(fnEntryBlockName))

		ctx.function = funcFn
		ctx.block = llBlock

		res := ctx.gen(node.Child)

		if ctx.block.Term == nil {
			ctx.block.NewRet(res)
		}
		return funcFn
	}
}

func (ctx *context) genBranch(
	block *ir.Block,
	node *mTypes.Node,
	condRet value.Value,
	exitBlock *ir.Block,
) {
	ctx.block = block
	res := ctx.gen(node)
	retType := ctx.function.Sig.RetType
	isVoid := retType.Equal(types.Void)

	if res != nil && mTypes.IsScalar(res) {
		if retType.Equal(types.Void) {
			ctx.block.NewRet(nil)
		} else {
			ctx.block.NewRet(res)
		}
	} else {
		ctx.block.NewBr(exitBlock)
	}

	if res != nil && !isVoid {
		ctx.block.NewStore(res, condRet)
	}
}

func (ctx *context) genCondition(node *mTypes.Node) {
	// cond
	condBlock := ctx.function.NewBlock(node.GetBlockName("if.cond"))
	ctx.block.NewBr(condBlock)
	ctx.block = condBlock
	cond := ctx.gen(node.Cond)

	retType := ctx.function.Sig.RetType
	isVoid := retType.Equal(types.Void)

	// NOTE: is it the type truly?
	if !isVoid {
		node.CondRet = ctx.block.NewAlloca(retType)
	}

	// exit
	// NOTE: is it the type truly?
	exitBlock := ctx.function.NewBlock(node.GetBlockName("if.exit"))

	if retType.Equal(types.Void) {
		exitBlock.NewRet(nil)
	} else {
		exitBlock.NewRet(exitBlock.NewLoad(retType, node.CondRet))
	}

	// then
	thenBlock := ctx.function.NewBlock(node.GetBlockName("if.then"))
	ctx.genBranch(thenBlock, node.Then, node.CondRet, exitBlock)

	// else
	elseBlock := ctx.function.NewBlock(node.GetBlockName("if.else"))
	ctx.genBranch(elseBlock, node.Else, node.CondRet, exitBlock)

	condBlock.NewCondBr(cond, thenBlock, elseBlock)
	ctx.block = exitBlock
}

func (ctx *context) gen(node *mTypes.Node) value.Value {
	// Note: no more need?
	// log.DebugNoLine(log.GREEN(fmt.Sprintf("%+v \"%+v\"", node.Kind, node.Val)))
	if node.IsKind(mTypes.ND_DECLARE) {
		return ctx.gen(node.Child)

	} else if node.IsKind(mTypes.ND_VAR_DECLARE) {
		ctx.genVarDeclare(node)

	} else if node.IsKind(mTypes.ND_VAR_REFERENCE) {
		return ctx.genVarReference(node)

	} else if node.IsKind(mTypes.ND_LAMBDA) {
		return ctx.genLambda(node)

	} else if node.IsKind(mTypes.ND_BIND) {
		// add node.Bind to last element of ctx.scope
		if ctx.scope == nil {
			ctx.scope = node.Bind
		} else {
			lastScope := ctx.scope.GetLastNode()
			lastScope.Next = node.Bind
		}

		for bind := node.Bind; bind != nil; bind = bind.Next {
			child := ctx.gen(bind.Child)

			if bind.IsType(mTypes.TY_INT32) {
				bind.VarPtr = child

			} else if bind.IsType(mTypes.TY_STR) {
				bind.VarPtr = child

			} else if bind.IsType(mTypes.TY_BOOL) {
				bind.VarPtr = child

			} else if bind.IsType(mTypes.TY_VECTOR) {
				bind.VarPtr = child

			} else {
				log.Panic("unresolved NodeType: have %+v", node)
			}

		}
		return ctx.gen(node.Child)

	} else if node.IsKind(mTypes.ND_EXPR) {
		var res value.Value
		for child := node.Child; child != nil; child = child.Next {
			res = ctx.gen(child)
		}
		return res

	} else if node.IsKind(mTypes.ND_IF) {
		ctx.genCondition(node)

	} else if node.IsKind(mTypes.ND_LIBCALL) {
		// means calling standard library
		arg := ctx.gen(node.Child)
		node.Child.IRValue = arg

		for n := node.Child.Next; n != nil; n = n.Next {
			arg := ctx.gen(n)
			n.IRValue = arg
		}

		libFunc := core.LibInsts[node.Val]
		return libFunc(ctx.block, ctx.prog.BuiltinLibs, node.Child)

	} else if node.IsKind(mTypes.ND_FUNCCALL) {
		var arg []value.Value
		// generate ir of their arguments
		for node := node.Child; node != nil; node = node.Next {
			arg = append(arg, ctx.gen(node))
		}

		for i := 0; i < len(ctx.mod.Funcs); i = i + 1 {
			if ctx.mod.Funcs[i].GlobalName == node.GetFuncName() {
				return ctx.block.NewCall(ctx.mod.Funcs[i], arg...)
			}

		}
		log.Panic("unresolved function name: have %+v", node)

	} else if node.IsKind(mTypes.ND_SCALAR) {
		if node.IsType(mTypes.TY_INT32) {
			return newI32(node.Val)

		} else if node.IsType(mTypes.TY_STR) {
			if node.IsGlobal {
				return newStr(ctx, node)
			}
			return newStrHeap(ctx, node)

		} else if node.IsType(mTypes.TY_NIL) {
			return newStr(ctx, node)

		} else if node.IsType(mTypes.TY_BOOL) {
			return newBool(node.Val)

		} else {
			log.Panic("unresolved Scalar: have %+v", node)
		}
	} else if node.IsKind(mTypes.ND_COLLECTION) {

		return newVectorHeap(ctx, node)
	} else {
		log.Panic("unresolved Nodekind: have %+v", node)
	}
	return nil
}

func constructModule(prog *mTypes.Program) *ir.Module {
	module := ir.NewModule()
	prog.BuiltinLibs = &mTypes.BuiltinLibProp{}
	lib.DeclareBuiltin(module, prog.BuiltinLibs)

	for declare := prog.Declares; declare != nil; declare = declare.Next {
		c := &context{
			mod:  module,
			prog: prog,
		}
		c.gen(declare)
	}

	return module
}

func Construct(program *mTypes.Program) *assembler {
	return &assembler{
		program: program,
	}
}

func (a assembler) GenIntermediates(llName string, asmName string) {
	log.DebugMessage("ir module constructing")
	module := constructModule(a.program)
	log.DebugMessage("ir module constructed")
	log.Debug("[IR]\n%s\n", module.String())

	err := os.WriteFile(llName, []byte(module.String()), 0600)
	if err != nil {
		log.Panic("fail to write ll: %+v", map[string]interface{}{"err": err, "llName": llName})
	}
	log.Debug("written ll: %s", llName)

}
