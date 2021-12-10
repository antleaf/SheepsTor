package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	DebugLogging                   bool
	SourceRoot                     string     `yaml:"source_root"`
	WebRoot                        string     `yaml:"webroot"`
	Port                           int        `yaml:"port"`
	GitHubWebHookSecretEnvKey      string     `yaml:"github_webhook_secret_env_key"`
	AkismetApiKeyEnvKey            string     `yaml:"akismet_api_key_env_key"`
	WebmentionIoWebhookSecret      string     `yaml:"webmention_io_webhook_secret"`
	DisableGitCommitForDevelopment bool       `yaml:"disable_git_commit_for_development"`
	PublicAddress                  string     `yaml:"public_address"`
	Websites                       []*Website `yaml:"websites"`
}

func (config *Config) initialise(debugLogging bool, configFilePath string) error {
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

func (config *Config) configureWebsites() {
	for _, website := range config.Websites {
		logger.Debug(fmt.Sprintf("Configuring website '%s'", website.Id))
		website.Configure(config.SourceRoot, config.WebRoot)
		logger.Info(fmt.Sprintf("Website '%s' configured OK", website.Id))
	}
}

func (config *Config) getWebsiteByRepoNameAndBranchRef(repoName, branchRef string) *Website {
	for _, v := range config.Websites {
		if (v.GitRepo.RepoName == repoName) && (v.GitRepo.BranchRef == branchRef) {
			return v
		}
	}
	return nil
}

func (config *Config) getWebsiteByID(id string) *Website {
	for _, website := range config.Websites {
		if website.Id == id {
			return website
		}
	}
	return nil
}
