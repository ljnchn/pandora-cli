/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	// "pandora-cli/pkg/api"

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
	// _, err = api.GetModels()
	// if err != nil {
	// 	color.Red("api server error")
	// 	return
	// }

	// 获取 accounts.json 的数据
	viper.SetConfigName("accounts") // name of config file (without extension)
	viper.SetConfigType("json")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")        // optionally look for config in the working directory
	err = viper.ReadInConfig()      // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		color.Red("accounts.json not found")
		return
	}
	if err != nil {
		fmt.Println("Failed to read config file:", err)
		return
	}

	// 读取包含点号的键
	value := viper.GetStringMap("key.with.dot")
	fmt.Println("Value:", value)
	// 遍历配置
	for key, value := range viper.GetStringMap("admin@qq.com.fk") { // 获取 fk 下的配置
		fmt.Printf("Key: %s, Value: %#v\n", key, value)
	}
	fmt.Println(viper.GetString("admin@qq.password"))
	topLevelKeys := viper.AllSettings()
	if len(topLevelKeys) == 0 {
		color.Cyan("no accounts")
	}
	for email := range topLevelKeys {
		details := viper.GetString("admin@qq.share")
		color.Cyan("email: %s\n", email)
		color.Cyan("details: %s\n", details)
		// password := details["password"].(string)
		fmt.Printf("password: %s\n", details)
		// fks := viper.GetStringMap(email + ".fks")
		// for fkName, fkDetails := range fks {
		// 	fmt.Printf("fkName: %s\n", fkName)
		// 	token := fkDetails.(map[string]interface{})["token"].(string)
		// 	fmt.Printf("token: %s\n", token)
		// }
		return
	}

}
