# Pandora Cli
## 简单介绍
Pandora Cli 是一款帮你更好地使用 PandoraNext 开发的命令行工具。

## 安装
打开安装了 PandoraNext 目录，下载对应的 release 文件，修改 accounts.json.demo 为 accounts.json

示例:
``` bash
wget https://github.com/ljnchn/pandora-cli/releases/download/v0.01/pandora-cli-linux386-0.01.tar.gz

cp accounts.json.demo accounts.json

./pandora-cli

```

## 使用命令


- 查看 PandoraNext 服务状态
pandora-cli status

- 查看 token.json 列表
pandora-cli tokens

- 重载配置
pandora-cli relaod

- 刷新 share token
pandora-cli refresh

- 登陆账号
pandora-cli login



### 查看 PandoraNext 服务状态
```bash
./pandora-cli status
```
显示 config.json 中设置的参数，以及额度信息
![服务状态](./pic/image.png)

### 查看 token.json 列表
```bash
./pandora-cli tokens
```
显示 tokens.json 中的账户信息
![服务状态](./pic/image2.png)

### 重载配置
```bash
./pandora-cli reload
```
重载当前服务的config.json、tokens.json等配置

### 刷新 share token
```bash
./pandora-cli refresh
```
根据session token(有效期三个月)更新 accounts.json 中每个账号下 share token
![服务状态](./pic/image3.png)
### 登陆账号
```bash
./pandora-cli login
```
登陆 accounts.json 下面的账号