/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"pandora-cli/pkg/api"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "获取 Pandora Next 的服务状态",
	Long:  `获取 Pandora Next 的服务状态`,
	Run: func(cmd *cobra.Command, args []string) {
		getConfig()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getConfig() {
	result, err := api.GetJsonFromFile("config.json")
	if err != nil {
		color.Red(err.Error())
		return
	}

	running := false
	mode := "none"
	bind := result.Get("bind").String()
	apiPrefix := result.Get("proxy_api_prefix").String()
	webUrl := ""
	proxyUrl := ""
	if bind != "" {
		webUrl = fmt.Sprintf("http://%s", bind)
		// 发送 GET 请求
		resp, err := http.Get(webUrl + "/auth/login")
		if err == nil {
			// 检查状态码是否为 200
			if resp.StatusCode == http.StatusOK {
				running = true
				mode = "web"
			}
		}
		defer resp.Body.Close()

		if len(apiPrefix) > 7 {
			// 发送 GET 请求
			proxyUrl = fmt.Sprintf("http://%s/%s", bind, apiPrefix)
			resp2, err := http.Get(proxyUrl + "/v1/models")
			if err == nil {
				// 检查状态码是否为 200
				if resp2.StatusCode == http.StatusOK {
					running = true
					mode = "web & proxy"
				}

			}
			defer resp2.Body.Close()
		}
	}

	color.Cyan("%-15s %-10s \n", "bind: ", bind)
	color.Cyan("%-15s %-10s \n", "mode: ", mode)

	if running {
		color.Cyan("%-15s %-10s \n", "state: ", "running")
	} else {
		color.Red("%-15s %-10s \n", "state: ", "stoped")
	}
	color.Cyan("%-15s %-10s \n", "tls: ", result.Get("tls.enabled").String())
	if result.Get("license_id").String() != "" {
		color.Cyan("%-15s %-10s \n", "license: ", result.Get("license_id"))
	} else {
		color.Red("%-15s %-10s \n", "license: ", "no license")
	}
	color.Cyan("%-15s %-10s \n", "web url: ", webUrl)
	color.Cyan("%-15s %-10s \n", "proxy url: ", proxyUrl)
	if result.Get("public_share").String() != "" {
		color.Cyan("%-15s %-10s \n", "public share: ", result.Get("public_share"))
	}
	color.Cyan("%-15s %-10s \n", "site pass: ", result.Get("site_password"))
	color.Cyan("%-15s %-10s \n", "setup pass: ", result.Get("setup_password"))

}
