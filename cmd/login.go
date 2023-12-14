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

type accessJsonStruct struct {
	Exp int64 `json:"exp"`
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "自动刷新登陆 accoounts.json 文件中的账号",
	Long:  `自动刷新登陆 accoounts.json 文件中的账号`,
	Run: func(cmd *cobra.Command, args []string) {
		err := api.CheckService()
		if err != nil {
			color.Red(err.Error())
			return
		}
		if len(args) == 1 {
			parts := strings.SplitN(args[0], ",", 2)
			if len(parts) != 2 {
				color.Red("filename format error")
				return
			}
			email, password := parts[0], parts[1]
			_, err := loginOne(email, password)
			if err != nil {
				color.Red(err.Error())
				return
			}
			color.Green("登陆成功")
			return
		} else {
			loginAll()
		}
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

func loginAll() {
	// 读取 accounts.json 文件内容
	file, err := os.Open("accounts.json")
	if err != nil {
		color.Red("open accounts.json fail")
	}
	defer file.Close()

	// 读取文件内容
	bytes, err := io.ReadAll(file)
	if err != nil {
		color.Red("read accounts.json error")
	}
	result := gjson.ParseBytes(bytes)
	result.ForEach(func(email, accoount gjson.Result) bool {
		fmt.Printf("%s: ", email.String())
		// 从文件获取session，获取 access_token，成功则跳过登陆
		path := "./sessions/" + email.String() + ".json"
		content, err := os.ReadFile(path)
		if err == nil {
			accountInfo := gjson.ParseBytes(content)
			// 是否有 session
			sessionToken := accountInfo.Get("session_token").String()
			if sessionToken == "" {
				fmt.Println("session_token not found")
				return true
			}
			// session_tokne 换取 access_token
			body, err := api.GetAccessToken(email.String(), sessionToken)
			if err == nil {
				res, err := ParseAccess(body)
				if err == nil {
					tm := time.Unix(res.Exp, 0)
					expired := tm.Format("2006-01-02 15:04:05")
					fmt.Println("无需刷新，" + expired)
					return true
				}
			}
		}
		password := accoount.Get("password").String()
		// 执行登陆
		_, err = loginOne(email.String(), password)
		if err != nil {
			color.Red("登陆失败: %s", err.Error())
			return true
		}
		color.Green("登陆成功")
		return true
	})
}

// 通过邮箱密码登陆
func loginOne(email string, password string) (string, error) {
	body := ""
	// 生成token
	body, err := api.Login(email, password)
	if err != nil {
		return body, fmt.Errorf(body)
	}
	return body, nil
}

func NewLogin(email string, password string) (string, error) {
	body, err := api.Login(email, password)
	if err != nil {
		color.Red(body)
	}
	color.Green("login success")
	return body, err
}

func ParseAccess(access string) (accessJsonStruct, error) {
	var jsondata accessJsonStruct
	parts := strings.Split(access, ".")
	if len(parts) != 3 {
		return jsondata, fmt.Errorf("access_token format error")
	}
	// 使用 base64 包的 RawStdEncoding.DecodeString 方法来解码
	tokenData, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err == nil {
		err := json.Unmarshal(tokenData, &jsondata)
		if err == nil {
			return jsondata, nil
		}
	}
	return jsondata, err
}
