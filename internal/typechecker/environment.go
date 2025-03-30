package typechecker

type Context[T any] map[string]T

type Environment[T any] struct {
	contexts        []Context[T]
	currentFuncName string
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

func NewEnvironment[T any]() *Environment[T] {
	environment := Environment[T]{
		contexts: make([]Context[T], 0),
	}
	return &environment
}
