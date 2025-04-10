package internal

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Configuration struct {
	//LogLevel                  string          `yaml:"log_level"`
	SourceRoot                string          `yaml:"source_root"`
	DocsRoot                  string          `yaml:"docs_root"`
	Port                      int             `yaml:"port"`
	GitHubWebHookSecretEnvKey string          `yaml:"github_webhook_secret_env_key"`
	WebsiteConfigs            []WebsiteConfig `yaml:"websites"`
}

type GitRepoConfig struct {
	CloneId  string `yaml:"clone_id"`
	RepoName string `yaml:"repo_name"`
	Branch   string `yaml:"branch"`
}

type WebsiteConfig struct {
	ID                         string        `yaml:"id"`
	ContentProcessor           string        `yaml:"content_processor"` //either 'hugo' or nil
	ProcessorRootSubFolderPath string        `yaml:"processor_root"`    //e.g. a sub-folder in the repo called 'webroot'
	IndexForSearch             bool          `yaml:"index"`             //run the pagefind executable to create a search index
	GitRepoConfig              GitRepoConfig `yaml:"git"`
}

func NewConfiguration(configFilePath string) (Configuration, error) {
	config := Configuration{}
	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return config, err
	}
	return config, err
}
