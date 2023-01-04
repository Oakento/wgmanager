/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"wgmanager/app"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Update peer",
	Long:  "Update peer",
	Run: func(cmd *cobra.Command, args []string) {
		isSetKey := cmd.Flags().Changed("key")
		isSetAddr := cmd.Flags().Changed("address")
		if isSetKey || isSetAddr {
			app.CmdSetHandler(cmd, args)
		} else {
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	setCmd.Flags().BoolP("key", "k", false, "Generate new key pairs")
	setCmd.Flags().StringP("address", "a", "", "Set subnet ip address")
}
