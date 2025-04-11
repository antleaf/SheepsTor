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
	err := InitialiseConfiguration(configFilePath)
	if err != nil {
		fmt.Print(err.Error() + "\n")
		fmt.Printf("Halting execution because Config file not loaded from '%s'\n", configFilePath)
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
	InitialiseRegistry(Config.SourceRoot, Config.DocsRoot, os.Getenv(Config.GitHubWebHookSecretEnvKey))
	for _, w := range Config.WebsiteConfigs {
		website := NewWebsite(
			w.ID, w.ContentProcessor,
			w.ProcessorRootSubFolderPath,
			Config.SourceRoot,
			Config.DocsRoot,
			w.GitRepoConfig.CloneId,
			w.GitRepoConfig.RepoName,
			w.GitRepoConfig.Branch,
			w.IndexForSearch,
		)
		Registry.Add(&website)
	}
	Log.Infof("WebRoot folder path set to: %s", Config.DocsRoot)
	Log.Infof("Source Root folder path set to: %s", Config.SourceRoot)
}
