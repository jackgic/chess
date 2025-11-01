# 部署指南

## 本地开发部署

### 1. 环境准备
```bash
# 安装Go 1.21+
# macOS
brew install go

# Linux
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### 2. 配置项目
```bash
cd /Users/gongfeng/AI

# 复制配置文件
cp .env.example .env

# 编辑配置
vim .env
```

### 3. 运行项目
```bash
# 方式1: 使用启动脚本
chmod +x start.sh
./start.sh

# 方式2: 直接运行
go run main.go

# 方式3: 编译后运行
go build -o chinese-chess-ai
./chinese-chess-ai
```

## Docker部署

### 1. 创建Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o chinese-chess-ai main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/chinese-chess-ai .
COPY --from=builder /app/web ./web

EXPOSE 8080
CMD ["./chinese-chess-ai"]
```

### 2. 构建和运行
```bash
# 构建镜像
docker build -t chinese-chess-ai .

# 运行容器
docker run -d \
  -p 8080:8080 \
  -e TENCENT_SECRET_ID=your_id \
  -e TENCENT_SECRET_KEY=your_key \
  -e LKE_BOT_BIZ_ID=your_bot_id \
  -e LKE_APP_ID=your_app_id \
  --name chess-game \
  chinese-chess-ai
```

## 云服务器部署

### 腾讯云轻量应用服务器

1. **购买服务器**
   - 选择Ubuntu 20.04或更高版本
   - 至少1核2G配置

2. **安装环境**
```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装Git
sudo apt install git -y
```

3. **部署项目**
```bash
# 克隆项目
git clone <your-repo-url>
cd chinese-chess-ai

# 配置环境变量
cp .env.example .env
vim .env

# 运行项目
nohup go run main.go > app.log 2>&1 &
```

4. **配置Nginx反向代理**
```bash
# 安装Nginx
sudo apt install nginx -y

# 配置Nginx
sudo vim /etc/nginx/sites-available/chess
```

Nginx配置内容：
```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

```bash
# 启用配置
sudo ln -s /etc/nginx/sites-available/chess /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

5. **配置系统服务**
```bash
# 创建systemd服务
sudo vim /etc/systemd/system/chess-game.service
```

服务配置内容：
```ini
[Unit]
Description=Chinese Chess AI Game
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/chinese-chess-ai
ExecStart=/usr/local/go/bin/go run main.go
Restart=on-failure
Environment="TENCENT_SECRET_ID=your_id"
Environment="TENCENT_SECRET_KEY=your_key"
Environment="LKE_BOT_BIZ_ID=your_bot_id"
Environment="LKE_APP_ID=your_app_id"

[Install]
WantedBy=multi-user.target
```

```bash
# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable chess-game
sudo systemctl start chess-game
sudo systemctl status chess-game
```

## 性能优化

### 1. 编译优化
```bash
# 减小二进制文件大小
go build -ldflags="-s -w" -o chinese-chess-ai main.go

# 使用upx压缩
upx --best chinese-chess-ai
```

### 2. 并发优化
在 `internal/game/manager.go` 中调整并发参数：
```go
// 限制最大游戏数量
const MaxGames = 1000

// 添加游戏清理机制
func (m *Manager) CleanupOldGames() {
    // 清理超过1小时未活动的游戏
}
```

### 3. 缓存优化
添加Redis缓存游戏状态：
```go
import "github.com/go-redis/redis/v8"

// 缓存游戏状态
func (m *Manager) CacheGameState(gameID string, state *Game) error {
    // 实现Redis缓存
}
```

## 监控和日志

### 1. 日志配置
```go
import "github.com/sirupsen/logrus"

func init() {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.SetLevel(logrus.InfoLevel)
}
```

### 2. 性能监控
```go
import "github.com/prometheus/client_golang/prometheus"

// 添加Prometheus指标
var (
    gamesCreated = prometheus.NewCounter(...)
    aiMovesDuration = prometheus.NewHistogram(...)
)
```

## 安全建议

1. **使用HTTPS**
   - 配置SSL证书（Let's Encrypt）
   - 强制HTTPS重定向

2. **API限流**
   - 使用中间件限制请求频率
   - 防止DDoS攻击

3. **输入验证**
   - 验证所有用户输入
   - 防止SQL注入和XSS

4. **密钥管理**
   - 不要将密钥提交到Git
   - 使用环境变量或密钥管理服务

## 故障排查

### 常见问题

1. **端口被占用**
```bash
# 查找占用端口的进程
lsof -i :8080
# 杀死进程
kill -9 <PID>
```

2. **LKE连接失败**
   - 检查密钥是否正确
   - 检查网络连接
   - 查看LKE服务状态

3. **前端无法加载**
   - 检查静态文件路径
   - 查看浏览器控制台错误
   - 检查CORS配置

## 备份和恢复

### 备份游戏数据
```bash
# 如果使用数据库
mysqldump -u root -p chess_game > backup.sql

# 如果使用文件存储
tar -czf backup.tar.gz /path/to/game/data
```

### 恢复数据
```bash
# 恢复数据库
mysql -u root -p chess_game < backup.sql

# 恢复文件
tar -xzf backup.tar.gz -C /path/to/restore
```
