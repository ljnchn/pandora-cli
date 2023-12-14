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
	"path/filepath"
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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := checkService()
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
			username, password := parts[0], parts[1]
			_, err := loginOne(username, password)
			if err != nil {
				color.Red(err.Error())
				return
			}
			color.Green("登陆成功")
			return
		} else {
			// login()
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

func login() {
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
	result := gjson.ParseBytes(bytes)
	result.ForEach(func(email, item gjson.Result) bool {
		// 从文件获取access token，有效则跳过，无效则登陆
		path := "./sessions/" + email.String() + ".json"
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
	})
	// 遍历 accounts 文件夹下面的所有文件
	filepath.Walk("./sessions", func(path string, info os.FileInfo, err error) error {
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
				res := gjson.ParseBytes(content)
				// 获取文件中的session与access
				session_token := res.Get("session_token").String()
				access_token := res.Get("access_token").String()
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

func checkService() error {
	result, err := api.GetJsonFromFile("./config.json")
	if err != nil { // Handle errors reading the config file
		return fmt.Errorf("read config.json error")
	}

	bind := result.Get("bind").String()
	if bind == "" {
		return fmt.Errorf("bind not found")
	}
	proxy_api_prefix := result.Get("proxy_api_prefix").String()
	if proxy_api_prefix == "" {
		return fmt.Errorf("proxy_api_prefix not found")
	}
	api.SetBaseUrl(fmt.Sprintf("http://%s/%s", bind, proxy_api_prefix))
	// 检查api服务
	_, err = api.GetModels()
	if err != nil {
		return fmt.Errorf("api server error")
	}
	return nil
}

func loginOne(username string, password string) (string, error) {
	body := ""
	// 生成token
	body, err := api.Login(username, password)
	if err != nil {
		return body, fmt.Errorf("login fail")
	}
	color.Green("login success")
	data := []byte(body)

	// 检查目录是否存在
	path := "./sessions"
	file := path + "/" + username + ".json"
	_, err = os.Stat("./sessions")
	if os.IsNotExist(err) {
		// 如果目录不存在，则创建目录
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			return body, fmt.Errorf("failed to create directory: %q\n", path)
		}
	}
	err = os.WriteFile(file, data, 0644)
	if err != nil {
		return body, fmt.Errorf("failed to write file: %q\n", file)
	}
	return body, nil
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
