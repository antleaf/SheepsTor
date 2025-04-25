package cmd

import (
	"fmt"
	. "github.com/antleaf/sheepstor/pkg"
	"github.com/antleaf/toolbox2go"
	"github.com/joho/godotenv"
	"os"
)

var registry *SheepstorWebsiteRegistry

func initialiseApplication() {
	var config = SheepstorConfiguration{}
	_ = godotenv.Load()
	configFilePath := os.Getenv("SHEEPSTOR_CONFIG_FILE_PATH")
	config, err := toolbox2go.NewConfigurationFromYamlFile(config, configFilePath)
	//err := InitialiseConfiguration(configFilePath)
	if err != nil {
		fmt.Print(err.Error() + "\n")
		fmt.Printf("Halting execution because Config file not loaded from '%s'\n", configFilePath)
		os.Exit(1)
	}
	log, err = toolbox2go.NewZapSugarLogger(Debug)
	if err != nil {
		fmt.Printf("Unable to initialise logging, halting: %s", err.Error())
		os.Exit(-1)
	}
	SetLogger(log)
	if Debug {
		log.Infof("Debugging enabled")
	}
	reg := NewSheepstorWebsiteRegistry(config.SourceRoot, config.DocsRoot, os.Getenv(config.GitHubWebHookSecretEnvKey))
	registry = &reg
	SetSheepstorRegistry(registry)
	for _, w := range config.WebsiteConfigs {
		website := NewSheepstorWebsite(
			w.ID, w.ContentProcessor,
			w.ProcessorRootSubFolderPath,
			config.SourceRoot,
			config.DocsRoot,
			w.GitRepoConfig.CloneId,
			w.GitRepoConfig.RepoName,
			w.GitRepoConfig.Branch,
			w.IndexForSearch,
		)
		registry.Add(&website)
	}
	log.Infof("WebRoot folder path set to: %s", config.DocsRoot)
	log.Infof("Source Root folder path set to: %s", config.SourceRoot)
}
