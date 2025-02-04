// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
//

package main

import (
	"slices"
	"strings"

	"github.com/mattermost/mattermost-dbt/internal/store"
	"github.com/mattermost/mattermost-dbt/model"
	"github.com/spf13/cobra"
)

func init() {
	verifyIndexesCmd.Flags().String("pgedge-config", "", "The location of the pgEdge config file")
	verifyIndexesCmd.Flags().String("schema-name", "public", "The database schema to compare indexes on")
	verifyIndexesCmd.Flags().Bool("fix-missing", true, "Create missing indexes")
}

var verifyIndexesCmd = &cobra.Command{
	Use:   "verify-index",
	Short: "Verify that database nodes have matching indexes",
	RunE: func(command *cobra.Command, args []string) error {
		command.SilenceUsage = true

		configLocation, _ := command.Flags().GetString("pgedge-config")
		config, err := model.NewClusterConfigFromFile(configLocation)
		if err != nil {
			return err
		}

		fixMissing, _ := command.Flags().GetBool("fix-missing")
		schemaName, _ := command.Flags().GetString("schema-name")
		indexFilter := &model.PgIndexFilter{SchemaName: schemaName}

		nodeStores, err := store.NewStoresForAllPgedgeNodes(config, logger)
		if err != nil {
			return err
		}

		allIndexes := make(map[string]model.PgIndexes, len(nodeStores))
		for _, nodeStore := range nodeStores {
			indexes, err := nodeStore.Store.GetIndexes(indexFilter)
			if err != nil {
				return err
			}

			logger.Infof("%s Index Count: %d", nodeStore.Node.Name, len(indexes))
			allIndexes[nodeStore.Node.Name] = indexes
		}

		// Perform a comparison of all nodes against all other nodes to find
		// missing indexes. Each comparison pass only looks for missing indexes
		// against a target so it must be performed in reverse to identify all
		// missing indexes.
		for sourceName, sourceIndexes := range allIndexes {
			for targetName, targetIndexes := range allIndexes {
				if sourceName == targetName {
					continue
				}

				logger.Infof("Comparing indexes of %s against %s", sourceName, targetName)
				missing := compareIndexes(sourceIndexes, targetIndexes)
				if len(missing) == 0 {
					logger.Info("All indexes found")
				} else {
					logger.Warnf("%d indexes missing: %s", len(missing), strings.Join(missing, ", "))
					if fixMissing {
						for _, m := range missing {
							logger.Infof("Creating missing index %s", m)
							store, err := nodeStores.GetStoreForNode(targetName)
							if err != nil {
								return err
							}
							err = store.CreateIndex(sourceIndexes[m])
							if err != nil {
								return err
							}
						}
					}
				}
			}
		}

		return nil
	},
}

// compareIndexes performs a comparison of indexes of a source against a target
// to identify indexes that are missing. This can be performed in reverse to
// identify all missing indexes.
func compareIndexes(source, target model.PgIndexes) []string {
	return findMissing(source.GetNames(), target.GetNames())
}

func findMissing(base, checkAgainst []string) []string {
	var missing []string
	for _, b := range base {
		if !slices.Contains(checkAgainst, b) {
			missing = append(missing, b)
		}
	}

	return missing
}
