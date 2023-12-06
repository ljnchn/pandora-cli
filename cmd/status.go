/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"net/http"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		getConfig()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("json")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		color.Red("config.json not found")
		return
	}

	running := false
	bind := viper.GetString("bind")
	if bind != "" {
		// 发送 GET 请求
		resp, err := http.Get("http://" + bind + "/auth/login")
		if err == nil {
			// 检查状态码是否为 200
			if resp.StatusCode == http.StatusOK {
				running = true
			}

		}
		defer resp.Body.Close()
	}

	color.Cyan("%-15s %-10s \n", "bind: ", bind)
	color.Cyan("%-15s %-10s \n", "mode: ", viper.GetString("server_mode"))

	if running {
		color.Cyan("%-15s %-10s \n", "state: ", "running")
	} else {
		color.Red("%-15s %-10s \n", "state: ", "stoped")
	}
	if viper.GetString("license_id") != "" {
		color.Cyan("%-15s %-10s \n", "license: ", viper.GetString("license_id"))
	} else {
		color.Red("%-15s %-10s \n", "license: ", "no license")
	}
	if viper.GetString("public_share") != "" {
		color.Cyan("%-15s %-10s \n", "public share: ", viper.GetString("public_share"))
	}
	if viper.GetString("site_password") != "" {
		color.Cyan("%-15s %-10s \n", "site pass: ", viper.GetString("site_password"))
	}
	if viper.GetString("setup_password") != "" {
		color.Cyan("%-15s %-10s \n", "setup pass: ", viper.GetString("setup_password"))
	}

}
