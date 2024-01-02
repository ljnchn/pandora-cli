/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/base64"
	"encoding/json"
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
	Short: "刷新 accounts.json 中的 access_token 和 fk",
	Long:  `刷新 accounts.json 中的 access_token 和 fk`,
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
	err := api.CheckService()
	if err != nil {
		api.SetBaseUrl("")
		color.Cyan("service not running, use fakeopen")
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
	// JSON 字符串
	jsonString := string(bytes)

	// 创建一个 map 用来存储解析后的数据
	var jsonMap map[string]interface{}

	// 解析 JSON 字符串到 map
	err = json.Unmarshal([]byte(jsonString), &jsonMap)
	if err != nil {
		fmt.Println(err)
	}

	result := gjson.ParseBytes(bytes)
	result.ForEach(func(email, item gjson.Result) bool {
		// 获取 access token
		accessToken, err := getAccessTokenByRefresh(email.String())
		if err != nil {
			accessToken, err = getAccessTokenBySession(email.String())
			if err != nil {
				color.Red("get access_token fail: %s", err.Error())
				return true
			}
		}
		// 解析 access token，判断过期时间
		var exp = int64(0)
		parts := strings.Split(accessToken, ".")
		if len(parts) == 3 {
			// 使用 base64 包的 RawStdEncoding.DecodeString 方法来解码
			accessData, err := base64.RawStdEncoding.DecodeString(parts[1])
			if err != nil {
				color.Red("access_token decode error")
			}
			exp = gjson.ParseBytes(accessData).Get("exp").Int()
		}
		tm := time.Unix(exp, 0)
		expired := tm.Format("2006-01-02 15:04:05")
		color.Cyan("%s, exp: %s", email.String(), expired)

		// 获取需要刷新的 fk
		share := item.Get("share")
		if share.Type != gjson.Null {
			share.ForEach(func(fkName, fkItems gjson.Result) bool {
				fmt.Print(fkName.String() + ": ")
				fk := ""
				fk, err = api.RefreshShare(accessToken, fkName.String(), fkItems)
				if err != nil {
					color.Red("refresh fail: " + err.Error())
				} else {
					jsonMap[email.String()].(map[string]interface{})["share"].(map[string]interface{})[fkName.String()].(map[string]interface{})["token_key"] = fk
					color.Green("refresh success: " + fk)
				}
				return true
			})
		}
		return true // keep iterating
	})
	// 保存到文件
	prettyJSON, err := json.MarshalIndent(jsonMap, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile("accounts.json", prettyJSON, 0644)
}

func getAccessTokenBySession(email string) (string, error) {
	accessToken := ""
	// 从文件获取access token
	path := "./sessions/" + email + ".json"
	content, err := os.ReadFile(path)
	if err != nil {
		color.Red("failed to read file: %q\n", path)
		return accessToken, fmt.Errorf("failed to read file: %q\n", path)
	}
	accessToken = gjson.ParseBytes(content).Get("access_token").String()
	sessionToken := gjson.ParseBytes(content).Get("session_token").String()
	if accessToken == "" {
		return accessToken, fmt.Errorf("access_token not found")
	}
	// 解析 access token，判断过期时间
	var exp = int64(0)
	parts := strings.Split(accessToken, ".")
	if len(parts) == 3 {
		// 使用 base64 包的 RawStdEncoding.DecodeString 方法来解码
		accessData, err := base64.RawStdEncoding.DecodeString(parts[1])
		if err != nil {
			return accessToken, fmt.Errorf("access_token decode error")
		}
		exp = gjson.ParseBytes(accessData).Get("exp").Int()
	}
	now := time.Now().Unix()
	// 如果过期时间小于当前时间，则刷新
	if exp < now {
		// 生成新的access token
		return api.GetAccessToken(email, sessionToken)
	}
	return accessToken, nil
}

func getAccessTokenByRefresh(email string) (string, error) {
	accessToken := ""
	// 从文件获取access token
	path := "./access/" + email + ".json"
	content, err := os.ReadFile(path)
	if err == nil {
		accessToken = gjson.ParseBytes(content).Get("access_token").String()
	}
	accessToken = gjson.ParseBytes(content).Get("access_token").String()

	if accessToken != "" {
		// 解析 access token，判断过期时间
		var exp = int64(0)
		parts := strings.Split(accessToken, ".")
		if len(parts) == 3 {
			// 使用 base64 包的 RawStdEncoding.DecodeString 方法来解码
			accessData, err := base64.RawStdEncoding.DecodeString(parts[1])
			if err != nil {
				return accessToken, fmt.Errorf("access_token decode error")
			}
			exp = gjson.ParseBytes(accessData).Get("exp").Int()
		}
		now := time.Now().Unix()
		if exp < now {
			accessToken = ""
		}
	}
	if accessToken == "" {
		// 生成新的access token
		// 从文件获取 refresh token
		path = "./refreshs/" + email + ".json"
		content, err = os.ReadFile(path)
		if err != nil {
			return accessToken, fmt.Errorf("failed to read file: %q\n", path)
		}
		refreshToken := gjson.ParseBytes(content).Get("refresh_token").String()
		if refreshToken == "" {
			return accessToken, fmt.Errorf("refresh_token not found")
		}
		return api.RefreshAccessToken(email, refreshToken)

	}
	return accessToken, nil

}
