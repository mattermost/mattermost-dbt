// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
//

package main

import (
	"fmt"

	"github.com/mattermost/mattermost-dbt/model"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print mmdbt version",
	RunE: func(command *cobra.Command, args []string) error {
		command.SilenceUsage = true

		fmt.Printf("mmdbt version: %s\n", model.ShortBuildHash())

		return nil
	},
}
