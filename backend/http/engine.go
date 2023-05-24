package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Engine struct {
	*gin.Engine
}

func NewEngine() *Engine {
	return &Engine{
		Engine: gin.Default(),
	}
}

func (e *Engine) GET(path string, handlerArgs ...any) {
	e.Engine.GET(path, flattenHandlers(handlerArgs)...)
}

func (e *Engine) POST(path string, handlerArgs ...any) {
	e.Engine.POST(path, flattenHandlers(handlerArgs)...)
}

func (e *Engine) PATCH(path string, handlerArgs ...any) {
	e.Engine.PATCH(path, flattenHandlers(handlerArgs)...)
}

func (e *Engine) DELETE(path string, handlerArgs ...any) {
	e.Engine.DELETE(path, flattenHandlers(handlerArgs)...)
}

func flattenHandlers(handlerArgs []any) []gin.HandlerFunc {
	handlers := []gin.HandlerFunc{}
	for _, handlerArg := range handlerArgs {
		switch h := handlerArg.(type) {
		case func(*gin.Context):
			handlers = append(handlers, h)
		case gin.HandlerFunc:
			handlers = append(handlers, h)
		case []func(*gin.Context):
			for _, handlerFunc := range h {
				handlers = append(handlers, handlerFunc)
			}
		case []gin.HandlerFunc:
			handlers = append(handlers, h...)
		default:
			panic(fmt.Sprintf("unknown type in switch: %T", h))
		}
	}
	return handlers
}
