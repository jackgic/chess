# 中国象棋AI项目 - 脚本使用说明

## 概述

本项目提供了一套完整的脚本来管理中国象棋AI服务的启动、停止、状态检查和日志查看。

## 脚本列表

### 1. 启动脚本 (`start.sh`)

**功能**: 编译并在后台启动中国象棋AI服务

**使用方法**:
```bash
./start.sh
```

**功能特性**:
- 自动设置环境变量
- 编译Go程序
- 后台运行服务
- 输出重定向到日志文件
- PID管理
- 启动状态检查

**输出**:
- 主日志: `logs/chinese-chess-ai.log`
- 错误日志: `logs/chinese-chess-ai.error.log`
- PID文件: `chinese-chess-ai.pid`

### 2. 停止脚本 (`stop.sh`)

**功能**: 优雅地停止后台运行的服务

**使用方法**:
```bash
./stop.sh
```

**功能特性**:
- 优雅停止进程 (SIGTERM)
- 强制停止备用方案 (SIGKILL)
- 自动清理PID文件
- 进程查找和清理
- 停止状态记录

### 3. 状态检查脚本 (`status.sh`)

**功能**: 检查服务运行状态和系统信息

**使用方法**:
```bash
./status.sh
```

**检查内容**:
- ✅ 进程运行状态
- ✅ 端口占用情况
- ✅ API健康检查
- ✅ 日志文件状态
- ✅ 系统资源使用
- ✅ 可用命令提示

### 4. 日志查看脚本 (`logs.sh`)

**功能**: 方便地查看各种日志

**使用方法**:
```bash
# 查看最后50行主日志
./logs.sh

# 实时跟踪日志
./logs.sh -f

# 查看错误日志
./logs.sh -e

# 查看最后100行日志
./logs.sh -n 100

# 清空所有日志
./logs.sh -c

# 显示帮助
./logs.sh -h
```

**参数说明**:
- `-f, --follow`: 实时跟踪日志 (类似 tail -f)
- `-e, --error`: 查看错误日志
- `-n, --lines NUM`: 显示最后 NUM 行 (默认50)
- `-c, --clear`: 清空日志文件
- `-h, --help`: 显示帮助信息

## 使用流程

### 启动服务
```bash
# 1. 启动服务
./start.sh

# 2. 检查状态
./status.sh

# 3. 查看日志
./logs.sh -f
```

### 停止服务
```bash
# 停止服务
./stop.sh

# 确认停止
./status.sh
```

### 日常维护
```bash
# 检查服务状态
./status.sh

# 查看最新日志
./logs.sh -n 20

# 查看错误日志
./logs.sh -e

# 实时监控日志
./logs.sh -f
```

## 目录结构

```
chinese-chess-ai/
├── start.sh              # 启动脚本
├── stop.sh               # 停止脚本
├── status.sh             # 状态检查脚本
├── logs.sh               # 日志查看脚本
├── chinese-chess-ai.pid  # PID文件 (运行时生成)
├── logs/                 # 日志目录
│   ├── chinese-chess-ai.log       # 主日志文件
│   └── chinese-chess-ai.error.log # 错误日志文件
└── ...
```

## 环境变量

脚本会自动设置以下环境变量:
- `LKE_APP_ID`: 腾讯云LKE应用ID
- `LKE_APP_KEY`: 腾讯云LKE应用密钥
- `LKE_BOT_BIZ_ID`: 机器人业务ID
- `PORT`: 服务端口 (默认8090)

## API端点

服务启动后，可以通过以下端点访问:
- 健康检查: `http://localhost:8090/api/health`
- 创建游戏: `http://localhost:8090/api/game/new`
- 玩家走子: `http://localhost:8090/api/game/move`
- 获取游戏状态: `http://localhost:8090/api/game/:id`
- AI走子: `http://localhost:8090/api/game/:id/ai-move`

## 故障排除

### 启动失败
1. 检查端口是否被占用: `lsof -i :8090`
2. 查看错误日志: `./logs.sh -e`
3. 检查环境变量配置

### 进程无响应
1. 检查进程状态: `./status.sh`
2. 查看系统资源: `top` 或 `htop`
3. 重启服务: `./stop.sh && ./start.sh`

### 日志文件过大
1. 清空日志: `./logs.sh -c`
2. 或手动清理: `> logs/chinese-chess-ai.log`

## 注意事项

1. **权限**: 确保脚本有执行权限 (`chmod +x *.sh`)
2. **端口**: 默认使用8090端口，确保端口未被占用
3. **日志**: 日志文件会持续增长，建议定期清理
4. **环境**: 脚本在macOS和Linux上测试通过
5. **依赖**: 需要Go环境和相关依赖包

## 快速命令参考

```bash
# 服务管理
./start.sh          # 启动服务
./stop.sh           # 停止服务
./status.sh         # 检查状态

# 日志管理
./logs.sh           # 查看日志
./logs.sh -f        # 实时日志
./logs.sh -e        # 错误日志
./logs.sh -c        # 清空日志

# 测试API
curl http://localhost:8090/api/health
```