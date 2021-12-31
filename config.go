package sheepstor

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Configuration struct {
	DebugLogging              bool
	SourceRoot                string          `yaml:"source_root"`
	WebRoot                   string          `yaml:"webroot"`
	Port                      int             `yaml:"port"`
	GitHubWebHookSecretEnvKey string          `yaml:"github_webhook_secret_env_key"`
	WebsiteConfigs            []WebsiteConfig `yaml:"websites"`
}

type GitRepoConfig struct {
	CloneId    string `yaml:"clone_id"`
	RepoName   string `yaml:"repo_name"`
	BranchName string `yaml:"branch_name"`
}

type WebsiteConfig struct {
	ID                         string        `yaml:"id"`
	Enabled                    bool          `yaml:"enabled"`
	ContentProcessor           string        `yaml:"content_processor"` //either 'hugo' or nil
	ProcessorRootSubFolderPath string        `yaml:"processor_root"`    //e.g. a sub-folder in the repo called 'webroot'
	ContentRootSubFolderPath   string        `yaml:"content_root"`      //for hugo this is 'content' by default
	GitRepoConfig              GitRepoConfig `yaml:"git"`
}

func (config *Configuration) initialise(debugLogging bool, configFilePath string) error {
	config.DebugLogging = debugLogging
	configData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(configData, config)
	if err != nil {
		return err
	}
	return err
}
