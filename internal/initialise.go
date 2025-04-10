package internal

import (
	"github.com/go-chi/chi/v5"
	"github.com/unrolled/render"
	"go.uber.org/zap"
)

var Config = Configuration{}
var Log *zap.SugaredLogger
var Router chi.Router
var Registry WebsiteRegistry
var Renderer *render.Render

func InitialiseConfiguration(configFilePath string) error {
	config, err := NewConfiguration(configFilePath)
	Config = config
	return err
}

func InitialiseLogger(debug bool) error {
	var err error
	Log, err = ConfigureZapSugarLogger(debug)
	return err
}

func InitialiseServer() {
	Router = NewRouter()
	Renderer = NewRenderer()
}

func InitialiseRegistry(sourceRoot string, docsRoot string, githubWebHookSecretEnvKey string) {
	Registry = NewRegistry(sourceRoot, docsRoot, githubWebHookSecretEnvKey)
}
