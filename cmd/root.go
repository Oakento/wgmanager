/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"wgmanager/app"
	"wgmanager/infra"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wgm",
	Short: "Wireguard Manager",
	Long:  "Wireguard Manager\nA CLI tool for Wireguard peer managing.",

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		app.Init(
			// cmd.PersistentFlags().Changed("config"),
			cmd.PersistentFlags().Changed("db"),
			cmd.PersistentFlags().Changed("wgconf"),
			cmd.PersistentFlags().Changed("subnet"),
		)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&app.AppConfigFile, "config", "c", infra.DEFAULT_CONFIG, "config file")
	rootCmd.PersistentFlags().StringVarP(&app.DbFile, "db", "d", infra.DEFAULT_DB, "database file")
	rootCmd.PersistentFlags().StringVarP(&app.WgConfigFile, "wgconf", "w", infra.DEFAULT_WG_CONFIG, "wireguard config file")
	rootCmd.PersistentFlags().StringVarP(&app.Subnet, "subnet", "N", infra.DEFAULT_SUBNET, "wireguard subnet")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
