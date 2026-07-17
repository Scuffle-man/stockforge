下面从**一台没有 Go 的电脑**开始，直到在 VS Code 里调试成功你给出的 MySQL 连接脚本。  
我会把 Go 安装、项目初始化、依赖拉取、数据库准备、VS Code 配置和调试选择 Go 版本全部串起来。

---

## 1. 安装 Go（选一个版本，建议 1.21+）
### Windows
- 去 https://go.dev/dl/ 下载 `.msi` 安装包（如 `go1.22.5.windows-amd64.msi`）
- 双击安装，全部默认（会自动添加 `PATH`）
- 重新打开终端，输入 `go version` 验证

### macOS
```bash
brew install go
```
或者下载官方 `.pkg` 安装。同样验证：`go version`

### Linux（以 Ubuntu 为例）
```bash
sudo apt update
sudo apt install golang-go
```
如果版本较老，推荐从官网下载 `tar.gz`，解压到 `/usr/local` 并配置环境变量（网上很多教程）。  
验证：`go version`

---

## 2. 准备 MySQL 连接环境
脚本里 DSN 是：
```
root:123456@tcp(127.0.0.1:3306)/mysql
```
**确保本机安装了 MySQL 且已启动**：
- Windows 可用 XAMPP、MySQL Installer；macOS 可 `brew install mysql` 并启动；Linux `apt install mysql-server`。
- 用命令行或工具（如 MySQL Workbench）登录测试：
```bash
mysql -u root -p123456 -h 127.0.0.1 -P 3306
```
如果能进入 MySQL 命令行，则连接可用。  
如果密码或用户不对，请修改 DSN 或创建用户：
```sql
ALTER USER 'root'@'localhost' IDENTIFIED BY '123456';
FLUSH PRIVILEGES;
```
另外防火墙不要拦截 `127.0.0.1:3306`。

> 如果你暂时没有 MySQL，也可以先跳过真实连接，先把 Go 代码跑通（初始化模块、下载依赖、编译通过），等有了数据库再真正运行。

---

## 3. 创建项目并拉取依赖
在你想要放代码的目录，打开终端：

```bash
mkdir go-mysql-demo
cd go-mysql-demo
go mod init go-mysql-demo
```
这样会生成 `go.mod` 文件。

把脚本代码保存为 `main.go`（就在 `go-mysql-demo` 目录里）。

接着拉取 MySQL 驱动依赖：
```bash
go get github.com/go-sql-driver/mysql
```
这会更新 `go.mod` 并下载包到本地缓存。

现在可以直接在终端试运行：
```bash
go run main.go
```
如果 MySQL 连上了，你会看到类似：
```
✅ 连接成功，MySQL 版本： 8.0.36
```

---

## 4. 在 VS Code 中打开并配置
### 安装 Go 扩展
1. 打开 VS Code，打开项目文件夹 `go-mysql-demo`
2. 左侧扩展市场搜索 `Go`（作者 Go Team at Google），安装
3. 安装后，按 `Ctrl+Shift+P` 输入 `Go: Install/Update Tools`，全选并安装（这样会安装 `gopls`、`dlv` 等调试工具）

### 选择 Go 版本（“解释器”）
Go 没有解释器，但你可以告诉扩展用哪个 `go` 二进制文件。

- 如果系统只装了一个 Go，什么都不用做，扩展会自动找到。
- 如果你装了多个 Go 版本，可以在 VS Code 设置里指定：

```json
"go.alternateTools": {
    "go": "/usr/local/go/bin/go1.22.5"
}
```
或者在 `launch.json` 里针对**某个调试配置**指定 `go` 路径（接下来马上讲）。

---

## 5. 配置调试（launch.json）
点击 VS Code 左侧“运行和调试”图标（或 `Ctrl+Shift+D`），点击“创建 launch.json 文件”，选择 **Go**。

默认会生成一个配置，把它改成这样：

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch MySQL Demo",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}",
            "go": "go"                     // 可改成绝对路径指定Go版本
        }
    ]
}
```

- `"program": "${workspaceFolder}"` 代表运行整个 main 包
- `"go": "go"` 使用系统 PATH 中的 `go`；若要指定版本，写完整路径，如 `"/usr/local/go/bin/go1.21.0"`
- 如果不想在配置里写死，可以省略 `go` 字段，扩展会自动用 `GOROOT` 或系统默认

---

## 6. 启动调试
1. 打开 `main.go`，在 `fmt.Println("✅ 连接成功...")` 这行打一个断点（点击行号左边）
2. 按 `F5` 选择 `Launch MySQL Demo` 配置开始调试
3. 程序会编译、运行，停在你打的断点处
4. 你可以看变量，单步执行等
5. 调试控制台会显示实际执行的 `go` 命令路径，可以验证是否用了你指定的版本

---

## 7. 常见问题与排查
| 问题 | 解决方法 |
|------|----------|
| 运行报 `Access denied for user` | 检查 MySQL 用户名密码和权限，`root` 是否允许从 `127.0.0.1` 登录 |
| `Error 1049: Unknown database 'mysql'` | DSN 里的数据库名是 `mysql`（系统库），通常存在。若没有，可换成 `information_schema` 测试 |
| `dial tcp 127.0.0.1:3306: connect: connection refused` | MySQL 服务没启动，或端口不对，启动服务再试 |
| `go: command not found` | Go 未正确安装或未加入 PATH，重新检查安装 |
| 调试时打不上断点 | 确保安装了 `dlv` 调试器（通过 `Go: Install/Update Tools` 安装） |
| 模块下载慢或超时 | 设置国内代理：`go env -w GOPROXY=https://goproxy.cn,direct` |

---

现在你有了一个从零到调试成功的完整流程。  
如果只是想调试那段 MySQL 连接代码，保证前三步（Go 安装、MySQL 可用、项目依赖拉取）完成后，直接在 VS Code 里 `F5` 就可以了。需要我帮你检查哪一步出错吗？