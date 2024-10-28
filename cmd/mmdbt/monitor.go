// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
//

package main

import (
	"time"

	"github.com/mattermost/mattermost-dbt/internal/store"
	"github.com/mattermost/mattermost-dbt/internal/tui"
	"github.com/mattermost/mattermost-dbt/model"
	"github.com/spf13/cobra"
)

func init() {
	monitorCmd.Flags().String("pgedge-config", "", "The location of the pgEdge config file")
	monitorCmd.Flags().Bool("include-mattermost-data", false, "Whether to monitor basic Mattermost table data or not")
	monitorCmd.Flags().Duration("refresh-timer", 3*time.Second, "The amount of time between monitoring refreshes")
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Provides sync and data status for multiple Mattermost databases",
	RunE: func(command *cobra.Command, args []string) error {
		command.SilenceUsage = true

		configLocation, _ := command.Flags().GetString("pgedge-config")
		includeMattermost, _ := command.Flags().GetBool("include-mattermost-data")
		refreshTimer, _ := command.Flags().GetDuration("refresh-timer")
		config, err := model.NewClusterConfigFromFile(configLocation)
		if err != nil {
			return err
		}

		nodeStores, err := store.NewStoresForAllPgedgeNodes(config, logger)
		if err != nil {
			return err
		}

		tui.StartMonitoring(refreshTimer, includeMattermost, nodeStores, logger)

		return nil
	},
}
