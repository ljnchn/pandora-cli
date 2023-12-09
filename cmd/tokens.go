/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// tokensCmd represents the tokens command
var tokensCmd = &cobra.Command{
	Use:   "tokens",
	Short: "获取 Pandora Next 的 tokens 信息",
	Long:  `获取 Pandora Next 的 tokens 信息`,
	Run: func(cmd *cobra.Command, args []string) {
		getTokens()
	},
}

func init() {
	rootCmd.AddCommand(tokensCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tokensCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tokensCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getTokens() {
	viper.SetConfigName("tokens") // name of config file (without extension)
	viper.SetConfigType("json")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	topLevelKeys := viper.AllSettings()

	// Loop through the keys and print them
	color.Cyan("%-10s %-10s %-10s %-10s %-10s %-10s \n", "account", "type", "pass", "plus", "shared", "expired")

	for key := range topLevelKeys {
		var expired = ""
		var token = viper.GetString(key + ".token")
		if token == "" {
			continue
		}
		var types string

		if strings.HasPrefix(token, "fk-") {
			types = "share"
		} else if strings.Contains(token, ",") {
			types = "account"
		} else {
			types = "session"
			parts := strings.Split(token, ".")
			if len(parts) == 3 {
				// 使用 base64 包的 RawStdEncoding.DecodeString 方法来解码
				tokenData, err := base64.RawStdEncoding.DecodeString(parts[1])
				if err == nil {
					type JsonStruct struct {
						Exp int64 `json:"exp"`
					}
					var jsondata JsonStruct
					err2 := json.Unmarshal(tokenData, &jsondata)
					if err2 == nil {
						types = "access"
						tm := time.Unix(jsondata.Exp, 0)
						expired = tm.Format("2006-01-02 15:04:05")
					}
				}
			}
		}

		var pass = viper.Get(key+".password") == nil
		var plus = viper.Get(key+".plus") == nil
		var shared = viper.Get(key+".shared") == nil
		fmt.Printf("%-10s %-10s %-10t %-10t %-10t %-10s \n", key, types, pass, plus, shared, expired)
	}
}
