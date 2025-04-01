package typechk

type Context[T any] map[string]T

func (c Context[T]) Has(key string) bool {
	_, ok := c[key]
	return ok
}

type Signature[T any] struct {
	paramNames []string
	params     map[string]T
	returns    T
}

type Environment[T any] struct {
	contexts      []Context[T]
	signatures    map[string]Signature[T]
	currentReturn T
}

func (e *Environment[T]) LookupVar(varName string) (T, bool) {
	for i := len(e.contexts) - 1; i >= 0; i-- {
		if value, exists := e.contexts[i][varName]; exists {
			return value, true
		}
	}

	var zeroValue T
	return zeroValue, false
}

func (e *Environment[T]) AssignVar(varName string, value T) bool {
	for i := len(e.contexts) - 1; i >= 0; i-- {
		if value, exists := e.contexts[i][varName]; exists {
			e.contexts[i][varName] = value
			return true
		}
	}
	return false
}

func (e *Environment[T]) ExtendVar(varName string, value T) bool {
	if len(e.contexts) == 0 {
		return false
	}
	ctx, ok := e.Peek()
	if !ok {
		return false
	}
	if ctx.Has(varName) {
		return false
	}
	e.contexts[len(e.contexts)-1][varName] = value
	return true
}

func (e Environment[T]) ReturnType() T {
	return e.currentReturn // possibly check for zero value
}

func (e *Environment[T]) SetReturnType(returnType T) {
	e.currentReturn = returnType
}

func (e *Environment[T]) EnterContext() {
	e.contexts = append(e.contexts, make(map[string]T))
}

func (e *Environment[T]) ExitContext() (Context[T], bool) {
	ctxLen := len(e.contexts)

	if len(e.contexts) == 0 {
		return nil, false
	}

	ctx := e.contexts[ctxLen-1]
	e.contexts = e.contexts[:ctxLen-1]

	return ctx, true
}

func (e *Environment[T]) AddStdFunc(
	funcName string,
	returns T,
	input T,
) {
	params := map[string]T{"input": input}
	e.ExtendFunc(
		funcName,
		[]string{"input"},
		params,
		returns,
	)
}

func (e *Environment[T]) ExtendFunc(
	funcName string,
	paramNames []string,
	params map[string]T,
	returns T,
) bool {
	if _, ok := e.signatures[funcName]; ok {
		return false
	}
	e.signatures[funcName] = Signature[T]{
		paramNames: paramNames,
		params:     params,
		returns:    returns,
	}
	return true
}

func (e *Environment[T]) LookupFunc(funcName string) (Signature[T], bool) {
	signatue, exists := e.signatures[funcName]
	if exists {
		return signatue, true
	}
	var zeroSignature Signature[T]
	return zeroSignature, false
}

func (e *Environment[T]) Peek() (*Context[T], bool) {
	if len(e.contexts) == 0 {
		return nil, false
	}
	return &e.contexts[len(e.contexts)-1], true
}

func NewEnvironment[T any]() *Environment[T] {
	var zeroValue T
	environment := Environment[T]{
		contexts:      make([]Context[T], 0),
		signatures:    make(map[string]Signature[T]),
		currentReturn: zeroValue,
	}
	return &environment
}
