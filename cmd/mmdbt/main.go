// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
//

package main

import (
	"os"

	"github.com/ory/viper"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	viper.SetEnvPrefix("MMDBT")
	viper.AutomaticEnv()

	rootCmd.AddCommand(verifyIndexesCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error(errors.Wrap(err, "Command failed").Error())
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:           "mmdbt",
	Short:         "Tooling for reviewing Mattermost databases",
	SilenceErrors: true,
}
