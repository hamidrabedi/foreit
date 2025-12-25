package api

import (
	"github.com/gofiber/fiber/v2"
)

type ActionConfig struct {
	Detail    bool
	Methods   []string
	URLPath   string
	URLName   string
}

type Action struct {
	Config ActionConfig
	Handler func(c *fiber.Ctx) error
}

type ActionMap map[string]Action

func NewAction(handler func(c *fiber.Ctx) error, config ...ActionConfig) Action {
	actionConfig := ActionConfig{
		Detail:  false,
		Methods: []string{"GET", "POST"},
		URLPath: "",
		URLName: "",
	}
	if len(config) > 0 {
		actionConfig = config[0]
		if actionConfig.Methods == nil {
			actionConfig.Methods = []string{"GET", "POST"}
		}
	}
	return Action{
		Config:  actionConfig,
		Handler: handler,
	}
}

func ActionDetail(handler func(c *fiber.Ctx) error, config ...ActionConfig) Action {
	actionConfig := ActionConfig{
		Detail:  true,
		Methods: []string{"GET", "POST"},
		URLPath: "",
		URLName: "",
	}
	if len(config) > 0 {
		actionConfig = config[0]
		actionConfig.Detail = true
		if actionConfig.Methods == nil {
			actionConfig.Methods = []string{"GET", "POST"}
		}
	}
	return Action{
		Config:  actionConfig,
		Handler: handler,
	}
}

type ActionViewSet[T any] interface {
	ViewSet[T]
	GetActions() ActionMap
}

func RegisterViewSetWithActions[T any](app *fiber.App, path string, vs ViewSet[T]) {
	RegisterViewSet[T](app, path, vs)

	if actionVS, ok := vs.(ActionViewSet[T]); ok {
		actions := actionVS.GetActions()
		for actionName, action := range actions {
			actionPath := action.Config.URLPath
			if actionPath == "" {
				actionPath = actionName + "/"
			}

			if action.Config.Detail {
				fullPath := path + "/:id/" + actionPath
				for _, method := range action.Config.Methods {
					switch method {
					case "GET":
						app.Get(fullPath, action.Handler)
					case "POST":
						app.Post(fullPath, action.Handler)
					case "PUT":
						app.Put(fullPath, action.Handler)
					case "PATCH":
						app.Patch(fullPath, action.Handler)
					case "DELETE":
						app.Delete(fullPath, action.Handler)
					}
				}
			} else {
				fullPath := path + "/" + actionPath
				for _, method := range action.Config.Methods {
					switch method {
					case "GET":
						app.Get(fullPath, action.Handler)
					case "POST":
						app.Post(fullPath, action.Handler)
					case "PUT":
						app.Put(fullPath, action.Handler)
					case "PATCH":
						app.Patch(fullPath, action.Handler)
					case "DELETE":
						app.Delete(fullPath, action.Handler)
					}
				}
			}
		}
	}
}
