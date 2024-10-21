// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
//

package main

import (
	"fmt"

	"github.com/mattermost/mattermost-dbt/model"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.Flags().String("pgedge-config", "", "The location of the pgEdge config file")
	configCmd.Flags().Bool("sanitize", true, "Sanitize sensitive config values before printing")
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Prints pgEdge cluster config",
	RunE: func(command *cobra.Command, args []string) error {
		command.SilenceUsage = true

		configLocation, _ := command.Flags().GetString("pgedge-config")
		config, err := model.NewClusterConfigFromFile(configLocation)
		if err != nil {
			return err
		}

		sanitize, _ := command.Flags().GetBool("sanitize")
		if sanitize {
			config.Pgedge.Sanitize()
		}

		fmt.Println("")
		fmt.Println("CONFIG")
		fmt.Println("================================================")
		printJSON(config)

		return nil
	},
}
