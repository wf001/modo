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

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

type assembler struct {
	program *mTypes.Program
}

type Context struct {
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

func newStrGlobal(ctx *Context, n *mTypes.Node) *ir.InstLoad {
	strConst := constant.NewCharArrayFromString(n.Val)
	globalStr := ctx.mod.NewGlobalDef(fmt.Sprintf(".str.%d", len(ctx.mod.Globals)), strConst)
	globalStr.Linkage = enum.LinkagePrivate
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
	ctx.prog.DeclaredGlobalStr = append(ctx.prog.DeclaredGlobalStr, str)
	return str
}

func newStrHeap(ctx *Context, n *mTypes.Node) *ir.InstCall {
	strVal := n.Val
	strLen := len(strVal)

	mallocSize := constant.NewInt(types.I64, int64(strLen))
	dest := ctx.block.NewCall(ctx.prog.Internal.Cstd.Malloc, mallocSize)

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
		ctx.prog.Internal.Cstd.Memcpy,
		dest,
		srcPtr,
		mallocSize,
		constant.False,
	)
	return dest
}

// Note: remain here until it will be defined the strategy of memory lifecycle
func newVectorOld(ctx *Context, n *mTypes.Node) value.Value {
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

// Unlike strings, vectors use malloc, which means structures containing pointer types cannot be
// used directly in the global space. While a vector used as a local variable can be represented
// as a structure, and a global vector can be defined using types.ArrayType, the added complexity
// doesn't seem worth the benefit.
//
// Therefore, global vectors also use heap memory. At the point of view, this function is NOT USED
// currently.

// However, if fixed-size arrays — which are different from vectors — are implemented in the future,
// this logic may be reused. So, the decision to delete this function will be postponed until it is
// determined whether fixed-size arrays will be supported.
func newVectorGlobal(ctx *Context, n *mTypes.Node) value.Value {
	elemType, _ := mTypes.GetLLVMType(n.ElemType)

	var arrContent []constant.Constant
	arrLength := 0

	if elemType == types.I8Ptr {
		// Define each string as a global variable and get a pointer to its contents using GEP.
		for e := n.Child; e != nil; e = e.Next {
			strConst := constant.NewCharArrayFromString(e.Val)
			globalStr := ctx.mod.NewGlobalDef(
				fmt.Sprintf(".str.%d", len(ctx.mod.Globals)),
				strConst,
			)
			globalStr.Linkage = enum.LinkagePrivate
			globalStr.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
			globalStr.Immutable = true
			globalStr.Align = 1

			gep := constant.NewGetElementPtr(
				globalStr.ContentType,
				globalStr,
				constant.NewInt(types.I32, 0),
				constant.NewInt(types.I32, 0),
			)
			arrContent = append(arrContent, gep)
			arrLength++
		}
	} else {
		for e := n.Child; e != nil; e = e.Next {
			e.IRValue = ctx.gen(e)
			c, ok := e.IRValue.(constant.Constant)
			if !ok {
				log.Panic("Each array element must be constant.Constant: have %+v", e.IRValue)
			}
			arrContent = append(arrContent, c)
			arrLength++
		}
	}

	arrType := types.NewArray(uint64(arrLength), elemType)
	arrConst := constant.NewArray(arrType, arrContent...)

	vecGlobal := ctx.mod.NewGlobalDef(
		fmt.Sprintf(".vector.%d", len(ctx.mod.Globals)),
		arrConst,
	)
	vecGlobal.Align = 8

	n.IRValue = vecGlobal
	return vecGlobal
}

func newVectorHeap(ctx *Context, n *mTypes.Node) value.Value {
	var arrLength int64

	arrContent := []value.Value{}
	for e := n.Child; e != nil; e = e.Next {
		e.IRValue = ctx.gen(e)
		arrContent = append(arrContent, e.IRValue)
		arrLength++
	}
	elemType, _ := mTypes.GetLLVMType(n.ElemType)
	structedArrType, _ := mTypes.GetLLVMTypeForVector(n, ctx.prog.Prelude)

	typeSize := constant.NewInt(types.I64, int64(mTypes.GetBitWidth(elemType)))
	arrSize := constant.NewInt(types.I64, int64(arrLength))
	allocSize := ctx.block.NewMul(typeSize, arrSize)

	allocatedPtr := ctx.block.NewCall(ctx.prog.Internal.Cstd.Malloc, allocSize)
	arrayPtr := ctx.block.NewBitCast(allocatedPtr, types.NewPointer(elemType))

	for i, e := range arrContent {
		ptr := ctx.block.NewGetElementPtr(elemType, arrayPtr, constant.NewInt(types.I32, int64(i)))
		ctx.block.NewStore(e, ptr)
	}
	arrayIntAlloca := ctx.block.NewAlloca(structedArrType)

	arrElemPtr := ctx.block.NewGetElementPtr(
		structedArrType,
		arrayIntAlloca,
		newI32("0"),
		newI32("0"),
	)
	arrElemPtr.SetName(n.GetVarName("arrElemPtr", ctx.block.Insts))
	ctx.block.NewStore(arrayPtr, arrElemPtr)

	lenElemPtr := ctx.block.NewGetElementPtr(
		structedArrType,
		arrayIntAlloca,
		newI32("0"),
		newI32("1"),
	)
	lenElemPtr.SetName(n.GetVarName("lenElemPtr", ctx.block.Insts))
	ctx.block.NewStore(constant.NewInt(types.I64, arrLength), lenElemPtr)

	return arrayIntAlloca

}

func (ctx *Context) genVarDeclare(node *mTypes.Node) value.Value {
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
		retType, ok := mTypes.GetLLVMType(node.Type)

		if !ok {
			if node.Child.Kind == mTypes.ND_COLLECTION {
				structType, _ := mTypes.GetLLVMTypeForVector(node.Child, ctx.prog.Prelude)
				retType = types.NewPointer(structType)

			} else {
				structType, _ := mTypes.GetLLVMTypeForVector(node, ctx.prog.Prelude)
				retType = types.NewPointer(structType)
			}
		}

		funcName := node.GetFuncName()

		var arg []value.Value
		var argp []*ir.Param

		// define arguments type of function
		for a := node.Child.Args; a != nil; a = a.Next {
			ty, ok := mTypes.GetLLVMType(a.Type)
			if !ok {
				structTy, _ := mTypes.GetLLVMTypeForVector(a, ctx.prog.Prelude)
				ty = types.NewPointer(structTy)
			}

			arg = append(arg, ir.NewParam(a.Val, ty))
			argp = append(argp, ir.NewParam(a.Val, ty))
		}

		fnc := ctx.mod.NewFunc(
			funcName,
			retType,
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

func (ctx *Context) genVarReference(node *mTypes.Node) value.Value {
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

func (ctx *Context) genLambda(node *mTypes.Node) value.Value {
	isParentMain := ctx.function.GlobalName == "main"
	unnamedFuncName := node.GetUnnamedFuncName()
	fnEntryBlockName := "fn.entry"

	if isParentMain {
		funcFn := ctx.mod.NewFunc(
			unnamedFuncName,
			types.Void,
			ctx.function.Params...,
		)
		llBlock := funcFn.NewBlock(node.GetBlockName(fnEntryBlockName, ctx.function.Blocks))

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
		llBlock := funcFn.NewBlock(node.GetBlockName(fnEntryBlockName, ctx.function.Blocks))

		ctx.function = funcFn
		ctx.block = llBlock

		res := ctx.gen(node.Child)

		if ctx.block.Term == nil {
			ctx.block.NewRet(res)
		}
		return funcFn
	}
}

func (ctx *Context) genBranch(
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

func (ctx *Context) genCondition(node *mTypes.Node) {
	// cond
	condBlock := ctx.function.NewBlock(node.GetBlockName("if.cond", ctx.function.Blocks))
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
	exitBlock := ctx.function.NewBlock(node.GetBlockName("if.exit", ctx.function.Blocks))

	if retType.Equal(types.Void) {
		exitBlock.NewRet(nil)
	} else {
		exitBlock.NewRet(exitBlock.NewLoad(retType, node.CondRet))
	}

	// then
	thenBlock := ctx.function.NewBlock(node.GetBlockName("if.then", ctx.function.Blocks))
	ctx.genBranch(thenBlock, node.Then, node.CondRet, exitBlock)

	// else
	elseBlock := ctx.function.NewBlock(node.GetBlockName("if.else", ctx.function.Blocks))
	ctx.genBranch(elseBlock, node.Else, node.CondRet, exitBlock)

	condBlock.NewCondBr(cond, thenBlock, elseBlock)
	ctx.block = exitBlock
}

func (ctx *Context) gen(node *mTypes.Node) value.Value {
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

		libFunc := LibInsts[node.Val]
		return libFunc(ctx, node.Child)

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
				return newStrGlobal(ctx, node)
			}
			return newStrHeap(ctx, node)

		} else if node.IsType(mTypes.TY_NIL) {
			return newStrGlobal(ctx, node)

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
	prog.Internal = &mTypes.Internal{}
	prog.Internal.Cstd = &mTypes.Cstd{}
	prog.Internal.GlobalConst = &mTypes.GlobalConst{}
	declareInternal(module, prog.Internal)
	prog.Prelude = &mTypes.PreludeProps{}
	declarePrelude(module, prog.Prelude)

	for declare := prog.Declares; declare != nil; declare = declare.Next {
		c := &Context{
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
