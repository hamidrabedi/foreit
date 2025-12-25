package middleware

import (
	"github.com/gofiber/fiber/v2"
)

type HookType string

const (
	HookBeforeRequest  HookType = "before_request"
	HookAfterRequest   HookType = "after_request"
	HookBeforeResponse HookType = "before_response"
	HookAfterResponse  HookType = "after_response"
	HookOnError        HookType = "on_error"
)

type Hook func(ctx *Context) error

type HookRegistry struct {
	hooks map[HookType][]Hook
}

func NewHookRegistry() *HookRegistry {
	return &HookRegistry{
		hooks: make(map[HookType][]Hook),
	}
}

func (hr *HookRegistry) Register(hookType HookType, hook Hook) {
	hr.hooks[hookType] = append(hr.hooks[hookType], hook)
}

func (hr *HookRegistry) Execute(hookType HookType, ctx *Context) error {
	hooks := hr.hooks[hookType]
	for _, hook := range hooks {
		if err := hook(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (hr *HookRegistry) ExecuteAll(hookType HookType, ctx *Context) []error {
	var errors []error
	hooks := hr.hooks[hookType]
	for _, hook := range hooks {
		if err := hook(ctx); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

type Pipeline struct {
	middleware []Middleware
	hooks      *HookRegistry
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		middleware: make([]Middleware, 0),
		hooks:      NewHookRegistry(),
	}
}

func (p *Pipeline) Use(middleware ...Middleware) {
	p.middleware = append(p.middleware, middleware...)
}

func (p *Pipeline) AddHook(hookType HookType, hook Hook) {
	p.hooks.Register(hookType, hook)
}

func (p *Pipeline) Execute(fiberCtx *fiber.Ctx, handler Handler) error {
	pipelineCtx := NewContext(fiberCtx)
	
	if err := p.hooks.Execute(HookBeforeRequest, pipelineCtx); err != nil {
		return err
	}
	
	chain := buildChain(p.middleware, handler)
	
	err := chain(pipelineCtx)
	
	if err != nil {
		if hookErr := p.hooks.Execute(HookOnError, pipelineCtx); hookErr != nil {
			return hookErr
		}
		return err
	}
	
	if err := p.hooks.Execute(HookBeforeResponse, pipelineCtx); err != nil {
		return err
	}
	
	if err := p.hooks.Execute(HookAfterResponse, pipelineCtx); err != nil {
		return err
	}
	
	return nil
}

func buildChain(middleware []Middleware, handler Handler) Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	return handler
}

