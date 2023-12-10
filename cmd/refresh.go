/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"pandora-cli/pkg/api"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	proxy_api_prefix := viper.GetString("proxy_api_prefix")
	if proxy_api_prefix == "" {
		color.Red("proxy_api_prefix not found")
		return
	}
	// 检查api服务
	_, err = api.GetModels()
	if err != nil {
		color.Red("api server error")
		return
	}

	// 遍历 accounts 文件夹下面的所有文件
	err = filepath.Walk("./accounts", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			color.Red("prevent panic by handling failure accessing a path %q: %v\n", path)
			return err
		}
		if !info.IsDir() {
			color.Cyan("visited file: %q\n", path)
			// 获取文件名
			filename := filepath.Base(path)
			fmt.Printf("filename: %s\n", filename)
			// 使用 strings.SplitN 函数分割字符串
			parts := strings.SplitN(filename, "，", 2)
			if len(parts) != 2 {
				fmt.Printf("username: %s\n", parts)

			}
			// 将结果赋值给两个变量
			username, password := parts[0], parts[1]
			fmt.Printf("username: %s\n", username)
			fmt.Printf("password: %s\n", password)
			// 读取文件内容
			content, err := os.ReadFile(path)
			if err != nil {
				color.Red("failed to read file: %q\n", path)
				return err
			}
			color.Green("file content: %s\n", content)
			return err
		}
		return nil
	})
	if err != nil {
		return
	}

	api.SetBaseUrl(fmt.Sprintf("http://%s/%s", bind, proxy_api_prefix))
	err = api.GetAccessToken("")

	if err != nil {
		color.Red("refresh fail")
		return
	}
	color.Green("refresh success")
}
