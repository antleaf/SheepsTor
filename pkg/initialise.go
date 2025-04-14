package pkg

import (
	"github.com/antleaf/toolbox2go"
	"go.uber.org/zap"
)

var Config = SheepstorConfiguration{}
var log *zap.SugaredLogger
var Registry WebsiteRegistry

func InitialiseConfiguration(configFilePath string) error {
	config, err := toolbox2go.NewConfigurationFromYamlFile(Config, configFilePath)
	if err != nil {
		return err
	}
	Config = config
	return err
}

func SetLogger(logger *zap.SugaredLogger) {
	log = logger
}

func InitialiseRegistry(sourceRoot string, docsRoot string, githubWebHookSecretEnvKey string) {
	Registry = NewRegistry(sourceRoot, docsRoot, githubWebHookSecretEnvKey)
}
