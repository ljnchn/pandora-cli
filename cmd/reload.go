/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"pandora-cli/pkg/api"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reloadCmd represents the reload command
var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "重载当前服务的config.json、tokens.json等配置",
	Long:  `重载当前服务的config.json、tokens.json等配置`,
	Run: func(cmd *cobra.Command, args []string) {
		reload()
	},
}

func init() {
	rootCmd.AddCommand(reloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func reload() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("json")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		color.Red("config.json not found")
		return
	}
	bind := viper.GetString("bind")
	if bind == "" {
		color.Red("bind not found")
		return
	}
	setup_password := viper.GetString("setup_password")
	if setup_password == "" {
		color.Red("setup_password not found")
		return
	}
	proxy_api_prefix := viper.GetString("proxy_api_prefix")
	if setup_password == "" {
		color.Red("proxy_api_prefix not found")
		return
	}

	api.SetBaseUrl(fmt.Sprintf("http://%s/%s", bind, proxy_api_prefix))
	err = api.Reload()

	if err != nil {
		color.Red("reload fail")
		return
	}
	color.Green("reload success")
}
