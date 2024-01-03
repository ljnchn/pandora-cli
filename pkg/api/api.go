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
	// POST 获取 refresh token
	authRefreshPath = "/api/auth/login2"
	// POST /api/auth/platform/login 获取sess-开头的sess key，使用urlencode form传递username 和 password 参数。
	authSessPath = "/api/auth/platform/login"
	// 使用 refres token
	tokenRefreshPath = "/api/auth/refresh"
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

type FkParams struct {
	access_token       string
	unique_name        string
	site_limit         string
	expires_in         string
	show_conversations string
	show_userinfo      string
}

var baseUrl string

func SetBaseUrl(url string) {
	baseUrl = url
}

// 登陆操作，结果保存到 session 下面
func Sess(email string, password string) (string, error) {
	body := ""
	// 检查目录是否存在
	path := "./sess"
	file := path + "/" + email + ".json"
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// 如果目录不存在，则创建目录
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			return body, fmt.Errorf("failed to create directory: %q\n", path)
		}
	}

	data := url.Values{}
	data.Set("username", email)
	data.Set("password", password)
	data.Set("prompt", "login")
	options := RequestOptions{
		Headers: make(map[string]string),
		Timeout: 10 * time.Second,
		body:    []byte(data.Encode()),
	}
	options.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	body, err = Post(baseUrl+authSessPath, &options)
	if err != nil {
		return body, err
	}
	err = os.WriteFile(file, []byte(body), 0644)
	if err != nil {
		return body, fmt.Errorf("failed to writeFile %s", file)
	}
	return body, err
}

// 登陆操作，结果保存到 session 下面
func Login(email string, password string) (string, error) {
	body := ""
	// 检查目录是否存在
	path := "./sessions"
	file := path + "/" + email + ".json"
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// 如果目录不存在，则创建目录
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			return body, fmt.Errorf("failed to create directory: %q\n", path)
		}
	}

	data := url.Values{}
	data.Set("username", email)
	data.Set("password", password)
	options := RequestOptions{
		Headers: make(map[string]string),
		Timeout: 10 * time.Second,
		body:    []byte(data.Encode()),
	}
	options.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	body, err = Post(baseUrl+authLoginPath, &options)
	if err != nil {
		return body, err
	}
	err = os.WriteFile(file, []byte(body), 0644)
	if err != nil {
		return body, fmt.Errorf("failed to writeFile %s", file)
	}
	return body, err
}

// 登陆操作，结果保存到 refreshs 下面
func Login2(email string, password string) (string, error) {
	body := ""
	// 检查目录是否存在
	path := "./refreshs"
	file := path + "/" + email + ".json"
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// 如果目录不存在，则创建目录
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			return body, fmt.Errorf("failed to create directory: %q\n", path)
		}
	}
	data := url.Values{}
	data.Set("username", email)
	data.Set("password", password)
	options := RequestOptions{
		Headers: make(map[string]string),
		Timeout: 10 * time.Second,
		body:    []byte(data.Encode()),
	}
	options.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	body, err = Post(baseUrl+authRefreshPath, &options)
	if err != nil {
		return body, err
	}

	err = os.WriteFile(file, []byte(body), 0644)
	if err != nil {
		return body, fmt.Errorf("failed to writeFile %s", file)
	}
	RefreshAccessToken(email, gjson.Parse(body).Get("refresh_token").String())
	return body, err
}

// 通过session token获取access token
func RefreshAccessToken(email, refresh_token string) (string, error) {
	accessToken := ""
	// 检查目录是否存在
	path := "./access"
	savePath := path + "/" + email + ".json"
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// 如果目录不存在，则创建目录
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			return accessToken, fmt.Errorf("failed to create directory: %q\n", path)
		}
	}
	data := url.Values{}
	data.Set("refresh_token", refresh_token)
	options := RequestOptions{
		Headers: make(map[string]string),
		Timeout: 5 * time.Second,
		body:    []byte(data.Encode()),
	}
	options.Headers["Content-Type"] = "application/x-www-form-urlencoded"

	url := baseUrl + tokenRefreshPath
	if baseUrl == "" {
		return accessToken, fmt.Errorf("baseUrl is empty")
	}

	body, err := Post(url, &options)
	if err != nil {
		return accessToken, fmt.Errorf("请求失败")
	}
	res := gjson.Parse(body)
	accessToken = res.Get("access_token").String()
	if accessToken == "" {
		return accessToken, fmt.Errorf("获取失败")
	}
	// 保存到文件
	err = os.WriteFile(savePath, []byte(body), 0644)
	if err != nil {
		return accessToken, fmt.Errorf("save file to " + savePath + " fail")
	}
	return accessToken, nil
}

// 通过session token获取access token
func GetAccessToken(email, session_token string) (string, error) {
	savePath := "sessions/" + email + ".json"
	accessToken := ""
	data := url.Values{}
	data.Set("session_token", session_token)
	options := RequestOptions{
		Headers: make(map[string]string),
		Timeout: 5 * time.Second,
		body:    []byte(data.Encode()),
	}
	options.Headers["Content-Type"] = "application/x-www-form-urlencoded"

	url := baseUrl + authSessionPath
	if baseUrl == "" {
		url = "https://ai.fakeopen.com/auth/session"
	}

	body, err := Post(url, &options)
	if err != nil {
		return accessToken, fmt.Errorf("请求失败")
	}
	res := gjson.Parse(body)
	accessToken = res.Get("access_token").String()
	if accessToken == "" {
		return accessToken, fmt.Errorf("获取失败")
	}
	// 保存到文件
	err = os.WriteFile(savePath, []byte(body), 0644)
	if err != nil {
		return accessToken, fmt.Errorf("save file to " + savePath + " fail")
	}
	return accessToken, nil
}

// 根据 access token 获取 share token 信息
func RefreshShare(access_token string, unique_name string, json gjson.Result) (string, error) {
	if access_token == "" {
		return "", fmt.Errorf("access_token is empty")
	}
	if unique_name == "" {
		return "", fmt.Errorf("unique_name is empty")
	}
	fk := ""
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	data := url.Values{}
	// 赋值 FkParams
	data.Set("access_token", access_token)
	data.Set("unique_name", unique_name)
	data.Set("site_limit", json.Get("site_limit").String())
	data.Set("expires_in", "0")
	data.Set("show_conversations", "true")
	data.Set("show_userinfo", "true")
	if json.Get("expires_in").String() != "" {
		data.Set("expires_in", json.Get("expires_in").String())
	}
	if json.Get("show_conversations").String() != "" {
		data.Set("show_conversations", json.Get("show_conversations").String())
	}
	if json.Get("show_userinfo").String() != "" {
		data.Set("show_userinfo", json.Get("show_userinfo").String())
	}

	options := RequestOptions{
		Headers: headers,
		Timeout: 5 * time.Second,
		body:    []byte(data.Encode()),
	}
	url := baseUrl + tokenRegisterPath
	if baseUrl == "" {
		url = "https://ai.fakeopen.com/token/register"
	}

	body, err := Post(url, &options)
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
	body, err := Post(baseUrl+setupReloadPath, NewRequestOptions())

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
	body, err := Get(baseUrl+modelsPath, NewRequestOptionsWithBearer("pk-this-is-a-real-free-pool-token-for-everyone"))

	if err != nil {
		// 处理读取错误
		return "", err
	}
	return body, nil
}

func GetUsage(license string) (string, error) {
	url := fmt.Sprintf("https://dash.pandoranext.com/api/%s/usage", license)
	body, err := Get(url, NewRequestOptions())
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

func NewFkParams() *FkParams {
	return &FkParams{
		access_token:       "",
		unique_name:        "",
		site_limit:         "",
		expires_in:         "0",
		show_conversations: "true",
		show_userinfo:      "true",
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("status code: %d, read body error", resp.StatusCode)
	}
	// 检查状态码是否为 200
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code: %d\n%s", resp.StatusCode, body)
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

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("status code: %d, read body error", resp.StatusCode)
	}

	// 检查状态码是否为 200
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code: %d\n%s", resp.StatusCode, respBody)
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

func CheckService() error {
	result, err := GetJsonFromFile("./config.json")
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
	SetBaseUrl(fmt.Sprintf("http://%s/%s", bind, proxy_api_prefix))
	// 检查api服务
	_, err = GetModels()
	if err != nil {
		return fmt.Errorf("api server error")
	}
	return nil
}
