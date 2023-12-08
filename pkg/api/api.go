package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

const (
	// POST /api/auth/login 登录获取access token，使用urlencode form传递username 和 password 参数。
	authLoginPath = "/api/auth/login"
	// POST /api/auth/session 通过session token获取access token，使用urlencode form传递session_token参数。
	authSessionPath = "/api/auth/session"
	// POST /api/auth/refresh 通过refresh token获取access token，使用urlencode form传递refresh_token参数。
	authRefreshPath = "/api/auth/refresh"
	// GET /api/token/info/fk-xxx 获取share token信息，使用生成人的access token做为Authorization头，可查看各模型用量。
	tokenInfoPath = "/api/token/info"
	// POST /api/token/register 生成share token
	tokenRegisterPath = "/api/token/register"
	// POST /api/pool/update 生成更新pool token
	poolUpdatePath = "/api/pool/update"
	// POST /api/setup/reload 重载当前服务的config.json、tokens.json等配置。
	setupReloadPath = "/api/setup/reload"
)

var baseUrl string

func SetBaseUrl(url string) {
	baseUrl = url
}

func Reload() error {
	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", baseUrl+setupReloadPath, nil)
	if err != nil {
		return err
	}

	// 发送 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	// 发送 GET 请求
	if err != nil {
		// 处理请求错误
		return err
	}
	defer resp.Body.Close()

	// 检查状态码是否为 200
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("reload fail")
	}
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// 处理读取错误
		return err
	}
	err = viper.ReadConfig(bytes.NewBuffer(body)) // Find and read the config file
	if err != nil {                               // Handle errors reading the config file
		return err
	}
	if viper.GetInt("code") != 0 {
		return err
	}
	return nil
}
