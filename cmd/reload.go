/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"pandora-cli/pkg/api"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
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
	result, err := api.GetJsonFromFile("./config.json")
	if err != nil { // Handle errors reading the config file
		color.Red(err.Error())
		return
	}

	bind := result.Get("bind").String()
	if bind == "" {
		color.Red("bind not found")
		return
	}
	proxy_api_prefix := result.Get("proxy_api_prefix").String()
	if proxy_api_prefix == "" {
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
