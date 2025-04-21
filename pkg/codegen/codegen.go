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
	program  *mTypes.Program
	internal *mTypes.Internal
}

type Context struct {
	mod      *ir.Module
	function *ir.Func
	block    *ir.Block
	prog     *mTypes.Program
	scope    *mTypes.Node
	argument *mTypes.Node
	internal *mTypes.Internal
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
	ctx.prog.Declare.GlobalVar = append(ctx.prog.Declare.GlobalVar, str)
	return str
}

func newStrHeap(ctx *Context, n *mTypes.Node) *ir.InstCall {
	strVal := n.Val
	strLen := len(strVal)

	mallocSize := constant.NewInt(types.I64, int64(strLen))
	dest := ctx.block.NewCall(ctx.internal.Cstd.Malloc, mallocSize)

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
		ctx.internal.Cstd.Memcpy,
		dest,
		srcPtr,
		mallocSize,
		constant.False,
	)
	return dest
}

// Note: remain here until it will be defined the strategy of memory lifecycle
func newVectorOld(ctx *Context, n *mTypes.Node) value.Value {
	_, elemType, _ := mTypes.GetLLVMType(n, ctx.prog.Prelude)

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
	_, elemType, _ := mTypes.GetLLVMType(n, ctx.prog.Prelude)

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
	var vecLength int64

	vecContent := []value.Value{}
	for e := n.Child; e != nil; e = e.Next {
		e.IRValue = ctx.gen(e)
		vecContent = append(vecContent, e.IRValue)
		vecLength++
	}
	structedArrType, elemType, _ := mTypes.GetLLVMType(n, ctx.prog.Prelude)

	typeSize := constant.NewInt(types.I64, int64(mTypes.GetBitWidth(elemType)))
	vecSize := constant.NewInt(types.I64, int64(vecLength))
	allocSize := ctx.block.NewMul(typeSize, vecSize)

	allocatedPtr := ctx.block.NewCall(ctx.internal.Cstd.Malloc, allocSize)
	vecPtr := ctx.block.NewBitCast(allocatedPtr, types.NewPointer(elemType))

	for i, e := range vecContent {
		ptr := ctx.block.NewGetElementPtr(elemType, vecPtr, constant.NewInt(types.I32, int64(i)))
		ctx.block.NewStore(e, ptr)
	}
	vecIntAlloca := ctx.block.NewAlloca(structedArrType)

	vecElemPtr := ctx.block.NewGetElementPtr(
		structedArrType,
		vecIntAlloca,
		newI32("0"),
		newI32("0"),
	)
	vecElemPtr.SetName(n.GetVarName("vecElemPtr", ctx.block.Insts))
	ctx.block.NewStore(vecPtr, vecElemPtr)

	lenElemPtr := ctx.block.NewGetElementPtr(
		structedArrType,
		vecIntAlloca,
		newI32("0"),
		newI32("1"),
	)
	lenElemPtr.SetName(n.GetVarName("lenElemPtr", ctx.block.Insts))
	ctx.block.NewStore(constant.NewInt(types.I64, vecLength), lenElemPtr)

	return vecIntAlloca

}

func newStruct(
	ctx *Context,
	node *mTypes.Node,
) value.Value {

	structCtx := getExtendedType(ctx, node)

	// null pointer to struct: %struct* null
	nullStructPtr := constant.NewNull(types.NewPointer(structCtx.Types))

	// gep: getelementptr %struct, %struct* null, 1
	gepEndPtr := ctx.block.NewGetElementPtr(
		structCtx.Types,
		nullStructPtr,
		constant.NewInt(types.I32, 1),
	)

	// ptrtoint: i64 (size in bytes)
	mallocSize := ctx.block.NewPtrToInt(gepEndPtr, types.I64)

	// malloc 呼び出し（事前に @malloc を宣言しておくこと）

	rawPtr := ctx.block.NewCall(ctx.internal.Cstd.Malloc, mallocSize)

	// bitcast i8* → %struct*
	structPtr := ctx.block.NewBitCast(rawPtr, types.NewPointer(structCtx.Types))

	for n := node.Child; n != nil; n = n.Next.Next {
		field := n
		value := n.Next

		namePtr := ctx.block.NewGetElementPtr(
			structCtx.Types,
			structPtr,
			constant.NewInt(types.I32, 0), // first element
			constant.NewInt(types.I32, int64(structCtx.Field[field.Val].Pos)), // name field
		)
		v := ctx.gen(value)
		ctx.block.NewStore(v, namePtr)

	}
	return structPtr

}

func getExtendedType(ctx *Context, node *mTypes.Node) *mTypes.StructType {
	for k, v := range ctx.prog.Declare.Type.Struct {
		if k == node.Type.ExtendName {
			return v
		}
	}
	return nil
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

		var retType types.Type
		var n *mTypes.Node

		if node.Child.Kind == mTypes.ND_COLLECTION {
			n = node.Child
		} else {
			n = node
		}

		collType, scalarType, _ := mTypes.GetLLVMType(n, ctx.prog.Prelude)

		if collType != nil {
			retType = types.NewPointer(collType)
		} else if scalarType != nil {
			retType = scalarType
		} else if ctx.prog.Declare.Type != nil && ctx.prog.Declare.Type.Struct != nil {
			// Note: also get struct type from getLLVMType?
			structType := getExtendedType(ctx, node)
			retType = types.NewPointer(structType.Types)
		} else {
			log.Panic(":have %#+v", node)
		}

		funcName := node.GetFuncName()

		var argValue []value.Value
		var argParam []*ir.Param

		// define arguments type of function
		for a := node.Child.Args; a != nil; a = a.Next {
			var ty types.Type
			collType, scalarType, _ := mTypes.GetLLVMType(a, ctx.prog.Prelude)

			if collType != nil {
				ty = types.NewPointer(collType)
			} else if scalarType != nil {
				ty = scalarType
			} else if ctx.prog.Declare.Type != nil && ctx.prog.Declare.Type.Struct != nil {
				structType := getExtendedType(ctx, a)
				ty = structType.Types
			} else {
				log.Panic(":have %#+v", a)
			}

			p := ir.NewParam(a.Val, ty)
			argValue = append(argValue, p)
			argParam = append(argParam, p)
		}

		fnc := ctx.mod.NewFunc(
			funcName,
			retType,
			argParam...,
		)
		entryBlock := fnc.NewBlock("")

		ctx.function = fnc
		ctx.argument = node.Child.Args
		ctx.block = entryBlock
		child := ctx.gen(node.Child)
		node.FuncPtr = fnc

		if node.Child.IsKind(mTypes.ND_LAMBDA) {
			lambda := entryBlock.NewCall(child, argValue...)

			if lambda.Type().Equal(types.Void) {
				entryBlock.NewRet(nil)
			} else {
				entryBlock.NewRet(lambda)
			}
		} else {
			entryBlock.NewRet(child)
		}

	}
	return nil
}

func (ctx *Context) genStructTypeDeclare(node *mTypes.Node) {
	var typsArr []types.Type
	structField := map[string]mTypes.StructTypeField{}
	var pos uint64 = 0
	structType := types.NewStruct()
	structType.SetName(node.Val)

	for n := node.Child; n != nil; n = n.Next {
		// Note: is NOT TRUE
		_, scalarTy, _ := mTypes.GetLLVMType(n, ctx.prog.Prelude)
		typsArr = append(typsArr, scalarTy)
		f := structField[n.Val]
		f.Pos = pos
		f.Type = scalarTy
		structField[n.Val] = f
		pos++
	}

	structType.Fields = typsArr
	ctx.mod.NewTypeDef(node.Val, structType)

	if ctx.prog.Declare.Type == nil {
		ctx.prog.Declare.Type = &mTypes.ExtendedTypes{}
	}

	if ctx.prog.Declare.Type.Struct == nil {
		ctx.prog.Declare.Type.Struct = map[string]*mTypes.StructType{}
	}

	ctx.prog.Declare.Type.Struct[node.Val] = &mTypes.StructType{
		Name:  node.Val,
		Field: structField,
		Types: structType,
	}

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
				return scope.IRValue
			} else if scope.Child.IsKind(mTypes.ND_LIBCALL) {
				return scope.IRValue
			} else if scope.Child.IsType(mTypes.TY_VECTOR) {
				return scope.IRValue

			} else if scope.Child.IsType(mTypes.TY_EXTENDED) {
				return scope.IRValue

			} else {
				log.Panic("unresolved NodeType: have %+v", node)
			}
		}
	}

	// find in global variable which is declared with def
	for declare := ctx.prog.Declare.Func; declare != nil; declare = declare.Next {
		if declare.Child.Val == node.Val {
			return ctx.block.NewCall(declare.Child.FuncPtr)
		}
	}

	log.Debug("unresolved symbol, treated as struct field: '%s'", node.Val)

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
		entryBlock := funcFn.NewBlock(node.GetBlockName(fnEntryBlockName, ctx.function.Blocks))

		ctx.function = funcFn
		ctx.block = entryBlock

		ctx.gen(node.Child)

		ctx.block.NewRet(nil)

		return funcFn

	} else {
		funcFn := ctx.mod.NewFunc(
			unnamedFuncName,
			ctx.function.Sig.RetType,
			ctx.function.Params...,
		)
		entryBlock := funcFn.NewBlock(node.GetBlockName(fnEntryBlockName, ctx.function.Blocks))

		ctx.function = funcFn
		ctx.block = entryBlock

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
	if node.IsKind(mTypes.ND_DECLARE) {
		return ctx.gen(node.Child)

	} else if node.IsKind(mTypes.ND_VAR_DECLARE) {
		return ctx.genVarDeclare(node)

	} else if node.IsKind(mTypes.ND_TYPE_DECLARE) {
		ctx.genStructTypeDeclare(node)

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
				bind.IRValue = child

			} else if bind.IsType(mTypes.TY_STR) {
				bind.IRValue = child

			} else if bind.IsType(mTypes.TY_BOOL) {
				bind.IRValue = child

			} else if bind.IsType(mTypes.TY_VECTOR) {
				bind.IRValue = child

			} else if bind.IsType(mTypes.TY_EXTENDED) {
				bind.IRValue = child

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

		preludeFunc := PreludeFunction[node.Val]
		return preludeFunc(ctx, node.Child)

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
	} else if node.IsKind(mTypes.ND_COLLECTION) && node.IsType(mTypes.TY_VECTOR) {
		return newVectorHeap(ctx, node)

	} else if node.IsKind(mTypes.ND_COLLECTION) && node.IsType(mTypes.TY_EXTENDED) {
		return newStruct(ctx, node)

	} else {
		log.Panic("unresolved Nodekind: have %+v", node)
	}
	return nil
}

func constructModule(prog *mTypes.Program, internal *mTypes.Internal) *ir.Module {
	module := ir.NewModule()
	internal.Cstd = &mTypes.Cstd{}
	internal.GlobalConst = &mTypes.GlobalConst{}
	declareInternal(module, internal)
	prog.Prelude = &mTypes.PreludeProps{}
	declarePrelude(module, prog.Prelude)
	prog.Declare.Type = &mTypes.ExtendedTypes{}
	prog.Declare.Type.Struct = map[string]*mTypes.StructType{}
	prog.Declare.Type.LLVM = map[string]*types.Type{}

	for declare := prog.Declare.Func; declare != nil; declare = declare.Next {
		c := &Context{
			mod:      module,
			prog:     prog,
			internal: internal,
		}
		c.gen(declare)
	}

	return module
}

func Construct(program *mTypes.Program) *assembler {
	return &assembler{
		program:  program,
		internal: &mTypes.Internal{},
	}
}

func (a assembler) GenIntermediates(llName string, asmName string) {
	log.DebugMessage("ir module constructing")
	module := constructModule(a.program, a.internal)
	log.DebugMessage("ir module constructed")
	log.Debug("[IR]\n%s\n", module.String())

	err := os.WriteFile(llName, []byte(module.String()), 0600)
	if err != nil {
		log.Panic("fail to write ll: %+v", map[string]interface{}{"err": err, "llName": llName})
	}
	log.Debug("written ll: %s", llName)

}
