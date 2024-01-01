/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"pandora-cli/pkg/api"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// login2Cmd represents the login command
var login2Cmd = &cobra.Command{
	Use:   "login2",
	Short: "获取 refresh token，参数: 'email,password'",
	Long:  `获取 refresh token，参数: 'email,password'`,
	Run: func(cmd *cobra.Command, args []string) {
		// 检查是否强制执行
		if !forceFlag {
			// 询问用户是否确认操作
			fmt.Println("本操作需要消耗1000额度，是否继续? [Y/n]")
			if !ask() {
				fmt.Println("已取消...")
				return
			}
		}
		if len(args) == 1 {
			parts := strings.SplitN(args[0], ",", 2)
			if len(parts) != 2 {
				color.Red("filename format error")
				return
			}
			email, password := parts[0], parts[1]
			fmt.Printf("email: %s, password: %s\n", email, password)
			err := api.CheckService()
			if err != nil {
				color.Red(err.Error())
				return
			}
			body, err := api.Login2(email, password)
			if err != nil {
				color.Red(err.Error())
			}
			fmt.Println(body)
			return
		}
		color.Red("参数错误")
	},
}

func init() {
	rootCmd.AddCommand(login2Cmd)
	login2Cmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "是否跳过确认")
}

// 添加 Ask 标志到命令
var forceFlag bool

// 实现 Ask 功能
func ask() bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		s, _ := reader.ReadString('\n')
		s = strings.TrimSuffix(s, "\n")
		s = strings.ToLower(s)
		if len(s) > 1 {
			fmt.Fprintln(os.Stderr, "Please enter Y or N")
			continue
		}
		if strings.Compare(s, "n") == 0 {
			return false
		} else if strings.Compare(s, "y") == 0 {
			break
		} else {
			continue
		}
	}
	return true
}
