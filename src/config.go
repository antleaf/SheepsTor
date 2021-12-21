package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Configuration struct {
	DebugLogging                   bool
	SourceRoot                     string          `yaml:"source_root"`
	WebRoot                        string          `yaml:"webroot"`
	Port                           int             `yaml:"port"`
	GitHubWebHookSecretEnvKey      string          `yaml:"github_webhook_secret_env_key"`
	AkismetApiKeyEnvKey            string          `yaml:"akismet_api_key_env_key"`
	DisableGitCommitForDevelopment bool            `yaml:"disable_git_commit_for_development"`
	WebsiteConfigs                 []WebsiteConfig `yaml:"websites"`
}

type GitRepoConfig struct {
	CloneId       string `yaml:"clone_id"`
	RepoName      string `yaml:"repo_name"`
	BranchName    string `yaml:"branch_name"`
	BranchRef     string `yaml:"-"`
	RepoLocalPath string `yaml:"-"`
}

type PathProcessorConfig struct {
	Name                  string `yaml:"name"`
	FolderMatchExpression string `yaml:"folder_match_expression"`
	URLMatchExpression    string `yaml:"url_match_expression"`
	UrlGenerationPattern  string `yaml:"url_generation_pattern"`
	FileGenerationPattern string `yaml:"file_generation_pattern"`
}

type SheepsTorProcessorConfig struct {
	BaseURL              string                `yaml:"base_url"`
	PathProcessorConfigs []PathProcessorConfig `yaml:"path_processors"`
	IndieWebConfig       IndieWebConfig        `yaml:"indieweb"`
}

type IndieWebConfig struct {
	IndieAuthTokenEndpoint          string `yaml:"indieauth_token_endpoint"`
	MicroPubMediaEndpoint           string `yaml:"micropub_media_endpoint"`
	IndieAuthId                     string `yaml:"indie_auth_id"`
	DraftPosts                      bool   `yaml:"draft_posts"`
	WebMentionIoWebhookSecretEnvKey string `yaml:"webmention_io_webhook_secret_env_key"`
}

type WebsiteConfig struct {
	ID                         string `yaml:"id"`
	Enabled                    bool   `yaml:"enabled"`
	ContentProcessor           string `yaml:"content_processor"` //either 'hugo' or nil
	ProcessorRootSubFolderPath string `yaml:"processor_root"`    //e.g. a sub-folder in the repo called 'webroot'
	ContentRootSubFolderPath   string `yaml:"content_root"`      //for hugo this is 'content' by default
	//ProcessorRoot              string                   `yaml:"-"`
	//ContentRoot                string                   `yaml:"-"`
	//WebRoot                    string                   `yaml:"-"`
	GitRepoConfig       GitRepoConfig            `yaml:"git"`
	SheepsTorProcessing SheepsTorProcessorConfig `yaml:"sheepstor"`
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
