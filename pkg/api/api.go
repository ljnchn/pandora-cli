package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/viper"
)

const (
	// POST /api/auth/login 登录获取access token，使用urlencode form传递username 和 password 参数。
	authLoginPath = "/api/auth/login"
	// POST /api/auth/session 通过session token获取access token，使用urlencode form传递session_token参数。
	authSessionPath = "/api/auth/session"
	// GET /api/token/info/fk-xxx 获取share token信息，使用生成人的access token做为Authorization头，可查看各模型用量。
	tokenInfoPath = "/api/token/info"
	// POST /api/token/register 生成share token
	tokenRegisterPath = "/api/token/register"
	// POST /api/pool/update 生成更新pool token
	poolUpdatePath = "/api/pool/update"
	// POST /api/setup/reload 重载当前服务的config.json、tokens.json等配置。
	setupReloadPath = "/api/setup/reload"
	modelsPath      = "/v1/models"
)

type RequestOptions struct {
	Headers map[string]string
	Timeout time.Duration
	body    []byte
}

var baseUrl string

func SetBaseUrl(url string) {
	baseUrl = url
}

func Login(username string, password string) (string, error) {
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	options := RequestOptions{
		Headers: make(map[string]string),
		Timeout: 10 * time.Second,
		body:    []byte(data.Encode()),
	}
	options.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	body, err := Post(baseUrl+authLoginPath, &options)

	if err != nil {
		return body, fmt.Errorf("login fail")
	}
	return body, err
}

func GetAccessToken(session_token string) error {
	data := url.Values{}
	data.Set("session_token", session_token)
	options := RequestOptions{
		Headers: make(map[string]string),
		Timeout: 5 * time.Second,
		body:    []byte(data.Encode()),
	}

	body, err := Post(baseUrl+authSessionPath, &options)
	fmt.Println(body)
	if err != nil {
		return err
	}
	fmt.Println(body)
	err = viper.ReadConfig(bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}
	if viper.GetInt("code") != 0 {
		return err
	}
	return nil

}

// Reload 重载当前服务的config.json、tokens.json等配置。
func Reload() error {
	body, err := Get(baseUrl+setupReloadPath, NewRequestOptions())

	if err != nil {
		// 处理读取错误
		return err
	}
	err = viper.ReadConfig(bytes.NewBuffer([]byte(body))) // Find and read the config file
	if err != nil {                                       // Handle errors reading the config file
		return err
	}
	if viper.GetInt("code") != 0 {
		return err
	}
	return nil
}

func GetModels() (string, error) {
	body, err := Get(baseUrl+modelsPath, NewRequestOptions())

	if err != nil {
		// 处理读取错误
		return "", err
	}
	return body, nil
}

// NewRequestOptions 创建一个新的 RequestOptions 实例，设置默认值
func NewRequestOptions() *RequestOptions {
	return &RequestOptions{
		Headers: make(map[string]string),
		Timeout: 5 * time.Second,
		body:    []byte(""),
	}
}

// NewRequestOptions 创建一个新的 RequestOptions 实例，设置Bearer，设置默认值
func NewRequestOptionsWithBearer(token string) *RequestOptions {
	return &RequestOptions{
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", token),
		},
		Timeout: 5 * time.Second,
		body:    []byte(""),
	}
}

// HTTPGet 封装了一个带有超时和请求头的 HTTP GET 请求
func Get(url string, options *RequestOptions) (string, error) {
	client := &http.Client{
		Timeout: options.Timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	for key, value := range options.Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查状态码是否为 200
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("StatusFail")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// HTTPPostWithBody 封装了一个带有请求体、超时和请求头的 HTTP POST 请求
func Post(url string, options *RequestOptions) (string, error) {
	client := &http.Client{
		Timeout: options.Timeout,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(options.body))
	if err != nil {
		return "", err
	}

	for key, value := range options.Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// 检查状态码是否为 200
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("StatusFail")
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
