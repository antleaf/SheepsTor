package cmd

import (
	"fmt"
	. "github.com/antleaf/SheepsTor/internal"
	"github.com/spf13/cobra"
	"sync"
)

var sites string

func init() {
	rootCmd.AddCommand(updateWebsiteCmd)
	updateWebsiteCmd.Flags().StringVarP(&sites, "sites", "", "", "--sites all|<some_id>")
}

var updateWebsiteCmd = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {
		initialiseApplication()
		updateWebsites(sites)
	},
}

func updateWebsites(sites string) {
	Log.Info(fmt.Sprintf("Running as CLI Process, updating website(s): '%s'...", sites))
	if sites == "all" {
		processAllWebsites()
	} else {
		website := *registry.GetWebsiteByID(sites)
		err := website.ProvisionSources()
		if err != nil {
			Log.Error(err.Error())
			return
		}
		err = website.Build()
		if err != nil {
			Log.Error(err.Error())
		}
	}
}

func processWebsiteInSynchronousWorker(websitePtr *Website, wg *sync.WaitGroup) {
	website := *websitePtr
	err := website.ProvisionSources()
	if err != nil {
		Log.Error(err.Error())
	} else {
		Log.Infof("Provisioned sources for website: '%s'", website.ID)
		err = website.Build()
		if err != nil {
			Log.Error(err.Error())
		} else {
			Log.Infof("Built website: '%s'", website.ID)
		}
	}
	wg.Done()
}

func processAllWebsites() {
	var wg sync.WaitGroup
	for _, website := range registry.WebSites {
		wg.Add(1)
		go processWebsiteInSynchronousWorker(website, &wg)
	}
	wg.Wait()
}
