package service

import (
	"context"
	"errors"
)

type ActionService interface {
	Execute(ctx context.Context, action string, id any, input any) (any, error)
	ListActions(ctx context.Context, id any) ([]ActionDef, error)
}

type ActionDef struct {
	Name         string
	Label        string
	Description  string
	InputSchema  interface{}
	RequiresInput bool
}

type BaseActionService struct {
	actions map[string]ActionFunc
}

type ActionFunc func(ctx context.Context, id any, input any) (any, error)

func NewBaseActionService() *BaseActionService {
	return &BaseActionService{
		actions: make(map[string]ActionFunc),
	}
}

func (s *BaseActionService) RegisterAction(name, label, description string, handler ActionFunc) {
	s.actions[name] = handler
}

func (s *BaseActionService) Execute(ctx context.Context, action string, id any, input any) (any, error) {
	handler, ok := s.actions[action]
	if !ok {
		return nil, ErrActionNotFound
	}
	return handler(ctx, id, input)
}

func (s *BaseActionService) ListActions(ctx context.Context, id any) ([]ActionDef, error) {
	actions := make([]ActionDef, 0, len(s.actions))
	for name := range s.actions {
		actions = append(actions, ActionDef{Name: name})
	}
	return actions, nil
}

var (
	ErrActionNotFound = errors.New("action not found")
)

type ActionResourceService[T any] interface {
	ResourceService[T]
	ActionService
}

