/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"pandora-cli/pkg/api"
	"strings"
	"time"

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
	_, err = api.GetModels()
	if err != nil {
		color.Red("api server error")
		return
	}

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
	result.ForEach(func(email, item gjson.Result) bool {
		color.Cyan(email.String() + ": ")
		// 从文件获取access token
		path := "./sessions/" + email.String()
		content, err := os.ReadFile(path)
		if err != nil {
			color.Red("failed to read file: %q\n", path)
			return true
		}
		accessToken := gjson.ParseBytes(content).Get("access_token").String()
		sessionToken := gjson.ParseBytes(content).Get("session_token").String()
		if accessToken == "" {
			color.Red("access_token not found")
			return true
		}
		// 解析 access token，判断过期时间
		parts := strings.Split(accessToken, ".")
		if len(parts) == 3 {
			// 使用 base64 包的 RawStdEncoding.DecodeString 方法来解码
			accessData, err := base64.RawStdEncoding.DecodeString(parts[1])
			if err != nil {
				color.Red("access_token decode error")
			}
			exp := gjson.ParseBytes(accessData).Get("exp").Int()
			now := time.Now().Unix()
			// 如果过期时间小于当前时间，则刷新
			if exp < now {
				color.Green("refresh success")
				// 生成新的access token
				accessToken, err = api.GetAccessToken(email.String(), sessionToken)
				if err != nil {
					color.Red("get access_token fail")
					return true
				}
			}
		}
		// 获取需要刷新的 fk
		share := item.Get("share")
		if share.Type != gjson.Null {
			share.ForEach(func(fkName, fkItems gjson.Result) bool {
				fmt.Print(fkName.String() + ": ")
				fk := ""
				fk, err = api.RefreshShare(accessToken, fkName.String(), fkItems)
				if err != nil {
					color.Red("refresh fail")
				} else {
					color.Green("refresh success" + fk)
				}
				return true
			})
		}
		fmt.Println()
		return true // keep iterating
	})

}
