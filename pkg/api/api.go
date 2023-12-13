package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/tidwall/gjson"
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

// 登陆操作，结果保存到 session 下面
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

// 通过session token获取access token
func GetAccessToken(email, session_token string) (string, error) {
	accessToken := ""
	data := url.Values{}
	data.Set("session_token", session_token)
	options := RequestOptions{
		Headers: make(map[string]string),
		Timeout: 5 * time.Second,
		body:    []byte(data.Encode()),
	}
	options.Headers["Content-Type"] = "application/x-www-form-urlencoded"

	body, err := Post(baseUrl+authSessionPath, &options)
	if err != nil {
		return accessToken, fmt.Errorf("请求失败")
	}
	res := gjson.Parse(body)
	accessToken = res.Get("access_token").String()
	if accessToken != "" {
		return accessToken, fmt.Errorf("获取失败")
	}
	// 保存到文件
	path := "session/" + email
	err = os.WriteFile(path, []byte(body), 0644)
	if err != nil {
		return accessToken, fmt.Errorf("save file to " + path + "fail")
	}
	return accessToken, nil
}

// 根据 access token 获取 share token 信息
func refreshShare(accessToken, fkName string) (string, error) {
	fk := ""
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	data := url.Values{}
	data.Set("access_token", accessToken)
	data.Set("unique_name", fkName)
	data.Set("access_token", accessToken)
	data.Set("site_limit", "")
	data.Set("expires_in", "0")
	data.Set("show_conversations", "true")
	data.Set("show_userinfo", "true")
	options := RequestOptions{
		Headers: headers,
		Timeout: 5 * time.Second,
		body:    []byte(data.Encode()),
	}
	body, err := Post(baseUrl+tokenRegisterPath, &options)
	if err != nil {
		return fk, fmt.Errorf("请求失败")
	}
	res := gjson.Parse(body)
	fk = res.Get("token_key").String()
	if fk == "" {
		return fk, fmt.Errorf("获取失败")

	}
	return fk, nil
}

// Reload 重载当前服务的config.json、tokens.json等配置。
func Reload() error {
	body, err := Get(baseUrl+setupReloadPath, NewRequestOptions())

	if err != nil {
		// 处理读取错误
		return err
	}
	res := gjson.Parse(body)
	code := res.Get("code")
	if code.Type != gjson.Null && code.Int() == 0 {
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

func GetJsonFromFile(path string) (gjson.Result, error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return gjson.Result{}, fmt.Errorf(path + " not found")
	}
	defer file.Close()

	// 读取文件内容
	bytes, err := io.ReadAll(file)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("read file error:")
	}
	return gjson.ParseBytes(bytes), nil
}
