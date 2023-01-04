/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"wgmanager/app"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new peer",
	Long:  "Create a new peer",
	Run: func(cmd *cobra.Command, args []string) {

		app.CmdNewHandler(cmd, args)

	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	newCmd.Flags().StringSliceP("address", "a", []string{}, "Wireguard IP addresses.")
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
