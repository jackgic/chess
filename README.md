# 中国象棋 AI 对弈游戏

基于腾讯云智能体开发平台(LKE)的中国象棋AI对弈游戏，使用Golang后端和Web前端实现。

## 功能特性

- ✅ 完整的中国象棋规则实现
- ✅ 基于腾讯云LKE的AI对手
- ✅ 美观的Web界面
- ✅ 实时对局状态显示
- ✅ 走子历史记录
- ✅ AI思考过程展示
- ✅ **会话式对话优化**（节省90%+ token消耗）

## 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin (HTTP服务器)
- **SDK**: 腾讯云SDK (tencentcloud-sdk-go)

### 前端
- **技术**: HTML5 Canvas + JavaScript
- **样式**: CSS3

## 项目结构

```
chinese-chess-ai/
├── main.go                 # 程序入口
├── go.mod                  # Go模块依赖
├── internal/
│   ├── api/               # HTTP API处理器
│   │   └── handler.go
│   ├── chess/             # 象棋逻辑
│   │   └── board.go
│   ├── config/            # 配置管理
│   │   └── config.go
│   ├── game/              # 游戏管理
│   │   └── manager.go
│   └── lke/               # LKE客户端
│       └── client.go
├── web/
│   ├── index.html         # 主页面
│   └── static/
│       ├── style.css      # 样式文件
│       └── app.js         # 前端逻辑
├── .env.example           # 环境变量示例
└── README.md              # 项目文档
```

## 快速开始

### 1. 环境准备

确保已安装：
- Go 1.21 或更高版本
- Git

### 2. 克隆项目

```bash
cd /Users/gongfeng/AI
```

### 3. 安装依赖

```bash
go mod download
```

### 4. 配置环境变量

复制 `.env.example` 为 `.env` 并填写配置：

```bash
cp .env.example .env
```

编辑 `.env` 文件，填入你的腾讯云配置：

```env
TENCENT_SECRET_ID=your_secret_id
TENCENT_SECRET_KEY=your_secret_key
TENCENT_REGION=ap-guangzhou
LKE_BOT_BIZ_ID=your_bot_biz_id
LKE_APP_ID=your_app_id
PORT=8080
```

### 5. 配置LKE智能体

在腾讯云LKE平台创建智能体，并配置系统提示词。

**重要提示**：本项目使用会话式对话机制，系统提示词包含完整的游戏规则和格式要求，每次对话只发送简洁的走子信息，可节省90%以上的token消耗。

详细的系统提示词配置请参考：[LKE_SYSTEM_PROMPT_OPTIMIZED.md](LKE_SYSTEM_PROMPT_OPTIMIZED.md)

**配置步骤：**
1. 登录腾讯云LKE平台
2. 创建或编辑智能体
3. 在"系统提示词"中粘贴优化后的提示词（见上述文档）
4. 保存配置

**关键特性：**
- ✅ 系统提示词包含完整规则（一次性配置）
- ✅ 用户消息只发送本次走子（如："对手走了：炮二平五，现在轮到你了，请走子。"）
- ✅ 充分利用会话上下文，大幅减少token消耗
- ✅ 对话更加自然流畅

**优化效果对比：**
- 优化前：每次 ~800 tokens
- 优化后：每次 ~25 tokens
- 节省：**96.9%**

详细的优化说明请参考：
- [SESSION_CHAT_OPTIMIZATION.md](SESSION_CHAT_OPTIMIZATION.md) - 优化方案详解
- [OPTIMIZATION_COMPARISON.md](OPTIMIZATION_COMPARISON.md) - 前后对比

### 6. 运行程序

```bash
# 设置环境变量（或使用.env文件）
export TENCENT_SECRET_ID=your_secret_id
export TENCENT_SECRET_KEY=your_secret_key
export LKE_BOT_BIZ_ID=your_bot_biz_id
export LKE_APP_ID=your_app_id

# 运行程序
go run main.go
```

### 7. 访问游戏

打开浏览器访问：`http://localhost:8080`

## API接口

### 创建新游戏
```
POST /api/game/new
Content-Type: application/json

{
  "playerColor": 1  // 1=红方, 2=黑方
}
```

### 玩家走子
```
POST /api/game/move
Content-Type: application/json

{
  "gameId": "game_1",
  "fromRow": 6,
  "fromCol": 4,
  "toRow": 5,
  "toCol": 4
}
```

### 获取游戏状态
```
GET /api/game/:id
```

### AI走子
```
POST /api/game/:id/ai-move
```

## LKE智能体提示词详解

### 核心要求

1. **角色定位**：AI扮演黑方，需要遵守所有象棋规则
2. **输出格式**：必须输出 `MOVE: 数字数字-数字数字` 格式
3. **策略思维**：具备开局、中局、残局的基本策略

### 提示词优化建议

1. **增强合法性检查**：
   - 明确每种棋子的移动规则
   - 强调不能走出非法步骤
   - 提供常见错误示例

2. **提升对局水平**：
   - 添加常见开局定式
   - 提供战术组合示例（如双车错、马后炮）
   - 强调子力价值和兑子原则

3. **改进输出稳定性**：
   - 多次强调输出格式要求
   - 提供正确和错误示例对比
   - 使用特殊标记（如MOVE:）便于解析

## 象棋规则说明

### 棋子移动规则

- **将/帅**：九宫格内，每次一格
- **士**：九宫格内斜走
- **象/相**：田字格，不过河，象眼不能堵
- **马**：日字格，马腿不能堵
- **车**：直线，无限距离
- **炮**：直线移动，吃子需隔一子
- **兵/卒**：未过河向前，过河可横移

### 胜负判定

- 吃掉对方将/帅
- 对方无子可走（困毙）
- 将帅对面（白脸将）

## 开发说明

### 添加新功能

1. **悔棋功能**：在 `game/manager.go` 中添加 `Undo()` 方法
2. **保存对局**：实现对局序列化和反序列化
3. **难度等级**：调整LKE提示词或添加思考时间限制

### 调试技巧

1. **查看LKE响应**：检查日志中的 "LKE响应" 输出
2. **测试走子合法性**：使用 `board.IsValidMove()` 方法
3. **前端调试**：打开浏览器开发者工具查看网络请求

## 测试验证

### 验证会话式对话优化

运行测试脚本验证AI只接收简洁的走子信息：

```bash
# 重新编译和启动服务
./start.sh

# 运行会话式对话测试
./test_session_chat.sh
```

测试脚本会自动：
1. 创建新游戏
2. 玩家走子
3. AI走子
4. 显示发送给AI的提示词
5. 验证提示词是否简洁

**预期结果：**
```
[Game] 发送提示词给AI: 游戏开始，你是先手，请走第一步。
[Game] 发送提示词给AI: 对手走了：炮二平五，现在轮到你了，请走子。
```

如果看到大量的棋盘状态和规则说明，说明优化未生效，请检查代码是否重新编译。

## 常见问题

### Q: AI走子不合法怎么办？
A: 检查LKE智能体的系统提示词是否正确配置，确保强调了输出格式和合法性要求。参考 [LKE_SYSTEM_PROMPT_OPTIMIZED.md](LKE_SYSTEM_PROMPT_OPTIMIZED.md)

### Q: 如何提升AI水平？
A: 优化系统提示词，添加更多策略指导和战术示例。

### Q: 能否支持多人对战？
A: 当前版本仅支持人机对战，可以扩展为双人对战模式。

### Q: 为什么要使用会话式对话？
A: 会话式对话可以：
- 大幅减少token消耗（节省90%+）
- 提升响应速度（减少50%+处理时间）
- 让对话更加自然流畅
- 充分利用LKE的会话上下文能力

详见：[SESSION_CHAT_OPTIMIZATION.md](SESSION_CHAT_OPTIMIZATION.md)

### Q: 如何查看AI收到的提示词？
A: 查看日志文件：
```bash
tail -f logs/chinese-chess-ai.log | grep "发送提示词给AI"
```

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！

## 联系方式

如有问题，请提交Issue或联系开发者。
