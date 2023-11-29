/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"io"
	"net/http"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reloadCmd represents the reload command
var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", "htts://"+bind+"/setup/reload", nil)
	if err != nil {
		color.Red("创建请求失败:", err)
		return
	}

	// 添加自定义的 Header
	req.Header.Add("Authorization", "Bearer "+setup_password)
	req.Header.Add("Content-Type", "application/json")

	// 发送 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		color.Red("HTTP 请求失败:", err)
		return
	}
	// 发送 GET 请求
	if err != nil {
		// 处理请求错误
		color.Red("reload fail")
		return
	}
	defer resp.Body.Close()

	// 检查状态码是否为 200
	if resp.StatusCode != http.StatusOK {
		color.Red("reload fail")
		return
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// 处理读取错误
		color.Red("reload fail")
		return
	}

	viper.SetConfigType("json")

	err = viper.ReadConfig(bytes.NewBuffer(body)) // Find and read the config file
	if err != nil {                               // Handle errors reading the config file
		color.Red("reload fail")
		return
	}
	if viper.GetInt("code") != 0 {
		color.Red("reload fail")
		return
	}
	color.Green("reload success")
}
