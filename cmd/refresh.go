/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	// "pandora-cli/pkg/api"

	"pandora-cli/pkg/api"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "重载当前服务的config.json、tokens.json等配置",
	Long:  `重载当前服务的config.json、tokens.json等配置`,
	Run: func(cmd *cobra.Command, args []string) {
		refresh()
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// refreshCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// refreshCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func refresh() {
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
	// 检查api服务
	// _, err = api.GetModels()
	// if err != nil {
	// 	color.Red("api server error")
	// 	return
	// }

	// 读取 accounts.json 文件内容
	// 打开文件
	file, err := os.Open("accounts.json")
	if err != nil {
		color.Red("accounts.json not found")
	}
	defer file.Close()

	// 读取文件内容
	bytes, err := io.ReadAll(file)
	if err != nil {
		color.Red("read accounts.json error")
	}
	result = gjson.ParseBytes(bytes)
	result.ForEach(func(key, value gjson.Result) bool {
		// fmt.Println("Key:", key.String(), "Value:", value.String())
		fmt.Println(value.Get("password"))
		return true // keep iterating
	})

}
