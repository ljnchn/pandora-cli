/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"pandora-cli/pkg/api"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type accessJsonStruct struct {
	Exp int64 `json:"exp"`
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		login()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func login() {
	
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
	if proxy_api_prefix == "" {
		color.Red("proxy_api_prefix not found")
		return
	}
	api.SetBaseUrl(fmt.Sprintf("http://%s/%s", bind, proxy_api_prefix))
	// 检查api服务
	_, err = api.GetModels()
	if err != nil {
		color.Red("api server error")
		return
	}

	// 遍历 accounts 文件夹下面的所有文件
	err = filepath.Walk("./sessions", func(path string, info os.FileInfo, err error) error {
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
			parts := strings.SplitN(filename, ",", 2)
			if len(parts) != 2 {
				color.Red("filename format error")
				return err
			}
			fmt.Printf("parts0: %s\n", parts[0])
			fmt.Printf("parts1: %s\n", parts[1])
			// 将结果赋值给两个变量
			username, password := parts[0], parts[1]
			fmt.Printf("username: %s\n", username)
			fmt.Printf("password: %s\n", password)
			// 读取文件内容
			content, err := os.ReadFile(path)
			if err != nil {
				color.Red("failed to read file: %q\n", path)
				return nil
			}
			var expired = ""
			if len(content) == 0 {
				// 生成token
				body, err := api.Login(username, password)
				if err != nil {
					color.Red("login fail")
				}
				color.Green("login success")
				data := []byte(body)
				err = os.WriteFile(path, data, 0644)
				if err != nil {
					color.Red("failed to write file: %q\n", path)
					return nil
				}
			} else {
				// 获取文件中的session与access
				err = viper.ReadConfig(bytes.NewBuffer(content))
				session_token := viper.GetString("session_token")
				access_token := viper.GetString("access_token")
				if session_token != "" && access_token != "" {
					// 使用 base64 包的 RawStdEncoding.DecodeString 方法来解码
					accessData, err := ParseAccess(access_token)
					if err == nil {
						tm := time.Unix(accessData.Exp, 0)
						expired = tm.Format("2006-01-02 15:04:05")
					}
				}
			}
			fmt.Printf("expired: %s\n", expired)
		}
		return nil
	})
}

func NewLogin(username string, password string) (string, error) {
	body, err := api.Login(username, password)
	if err != nil {
		color.Red("login fail")
	}
	color.Green("login success")
	return body, err
}

func ParseAccess(access string) (accessJsonStruct, error) {
	// 使用 base64 包的 RawStdEncoding.DecodeString 方法来解码
	tokenData, err := base64.RawStdEncoding.DecodeString(access)
	var jsondata accessJsonStruct
	if err == nil {
		err := json.Unmarshal(tokenData, &jsondata)
		if err == nil {
			return jsondata, nil
		}
	}
	return jsondata, err
}
