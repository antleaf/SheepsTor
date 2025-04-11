package cmd

import (
	"fmt"
	. "github.com/antleaf/SheepsTor/internal"
	"github.com/joho/godotenv"
	"os"
)

func initialiseApplication() {
	_ = godotenv.Load()
	configFilePath := os.Getenv("SHEEPSTOR_CONFIG_FILE_PATH")
	err := (&config).Initialise(configFilePath)
	if err != nil {
		fmt.Print(err.Error() + "\n")
		fmt.Printf("Halting execution because config file not loaded from '%s'\n", configFilePath)
		os.Exit(1)
	}
	err = InitialiseLogger(Debug)
	if err != nil {
		fmt.Printf("Unable to initialise logging, halting: %s", err.Error())
		os.Exit(-1)
	}
	if Debug {
		Log.Infof("Debugging enabled")
	}
	InitialiseRegistry(config.SourceRoot, config.DocsRoot, os.Getenv(config.GitHubWebHookSecretEnvKey))
	for _, w := range config.WebsiteConfigs {
		website := NewWebsite(
			w.ID, w.ContentProcessor,
			w.ProcessorRootSubFolderPath,
			config.SourceRoot,
			config.DocsRoot,
			w.GitRepoConfig.CloneId,
			w.GitRepoConfig.RepoName,
			w.GitRepoConfig.Branch,
			w.IndexForSearch,
		)
		Registry.Add(&website)
	}
	Log.Infof("WebRoot folder path set to: %s", config.DocsRoot)
	Log.Infof("Source Root folder path set to: %s", config.SourceRoot)
}
