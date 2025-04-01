package typechk

import "fmt"

type Context[T any] map[string]T

func (c Context[T]) Has(key string) bool {
	_, ok := c[key]
	return ok
}

type Signature[T any] struct {
	params  map[string]T
	returns T
}

type Environment[T any] struct {
	contexts    []Context[T]
	signatures  map[string]Signature[T]
	currentFunc string
}

func (e *Environment[T]) PushContext() {
	e.contexts = append(e.contexts, make(map[string]T))
}

func (e *Environment[T]) PopContext() (Context[T], bool) {
	ctxLen := len(e.contexts)

	if len(e.contexts) == 0 {
		return nil, false
	}

	ctx := e.contexts[ctxLen-1]
	e.contexts = e.contexts[:ctxLen-1]

	return ctx, true
}

func (e *Environment[T]) ExtendFunc(
	funcName string,
	params map[string]T,
	returns T,
) error {
	if _, ok := e.signatures[funcName]; ok {
		return fmt.Errorf("redefinition of '%s'", funcName)
	}
	e.signatures[funcName] = Signature[T]{params: params, returns: returns}
	return nil
}

func (e *Environment[T]) Peek() (*Context[T], bool) {
	if len(e.contexts) == 0 {
		return nil, false
	}
	return &e.contexts[len(e.contexts)-1], true
}

func NewEnvironment[T any]() *Environment[T] {
	environment := Environment[T]{
		contexts:   make([]Context[T], 0),
		signatures: make(map[string]Signature[T]),
	}
	return &environment
}
