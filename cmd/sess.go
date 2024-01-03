/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"pandora-cli/pkg/api"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

// sessCmd represents the login command
var sessCmd = &cobra.Command{
	Use:   "sess",
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
			fmt.Printf("email: %s, password: %s\n", email, password)
			body, err := getSess(email, password)
			if err != nil {
				color.Red("登陆失败: %s", err.Error())
				color.Red(body)
				return
			}
			color.Green("登陆成功")
			// fmt.Println(body)
			fmt.Println(gjson.Parse(body).Get("login_info.user.session.sensitive_id").String())
			return
		} else {
			getSessAuto()
		}
	},
}

func init() {
	rootCmd.AddCommand(sessCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sessCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sessCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// 通过邮箱密码获取sess
func getSess(email string, password string) (string, error) {
	body := ""
	// 生成token
	body, err := api.Sess(email, password)
	if err != nil {
		return body, fmt.Errorf(body)
	}
	return body, nil
}

// 遍历 accounts.json 文件，自动获取sess
func getSessAuto() {
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
		path := "./sess/" + email.String() + ".json"
		// 判断文件是否存在
		_, err := os.Stat(path)
		if os.IsExist(err) {
			return true
		}
		password := accoount.Get("password").String()
		// 执行登陆
		body, err := getSess(email.String(), password)
		if err != nil {
			color.Red("登陆失败: %s", err.Error())
			color.Red(body)
			return false
		}
		color.Green("登陆成功")
		// fmt.Println(body)
		fmt.Println(gjson.Parse(body).Get("login_info.user.session.sensitive_id").String())
		return true
	})
}
