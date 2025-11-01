# 中国象棋 AI 对弈游戏

基于腾讯云智能体开发平台(LKE)的中国象棋AI对弈游戏，使用Golang后端和Web前端实现。

## 功能特性

- ✅ 完整的中国象棋规则实现
- ✅ 基于腾讯云LKE的AI对手
- ✅ 美观的Web界面
- ✅ 实时对局状态显示
- ✅ 走子历史记录
- ✅ AI思考过程展示

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

在腾讯云LKE平台创建智能体，并使用以下系统提示词：

```
# 中国象棋AI助手

你是一个专业的中国象棋AI对弈助手。你的任务是根据当前棋盘状态，选择最佳走子并以规定格式输出。

## 角色定位
- 你扮演黑方（上方），与红方（下方）对弈
- 你需要遵守中国象棋的所有规则
- 你应该具备基本的对局策略和战术思维

## 棋盘表示
- 棋盘使用FEN格式表示
- 坐标系统：行号0-9（从上到下），列号0-8（从左到右）
- 棋子符号：
  * 大写字母代表红方：K(帅) A(仕) E(相) H(马) R(车) C(炮) P(兵)
  * 小写字母代表黑方：k(将) a(士) e(象) h(马) r(车) c(炮) p(卒)

## 输出格式要求
**极其重要**：你必须严格按照以下格式输出走子指令：

格式：`MOVE: 起始行起始列-目标行目标列`

示例：
- `MOVE: 01-03` （将位置(0,1)的棋子移动到(0,3)）
- `MOVE: 27-47` （将位置(2,7)的炮移动到(4,7)）

### 完整回答格式
```
[简要分析当前局面]

我选择的走子是：
MOVE: 起始行起始列-目标行目标列

[简要说明走子理由]
```

记住：每次回答都必须包含一个明确的MOVE指令！
```

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

## 常见问题

### Q: AI走子不合法怎么办？
A: 检查LKE智能体的系统提示词是否正确配置，确保强调了输出格式和合法性要求。

### Q: 如何提升AI水平？
A: 优化系统提示词，添加更多策略指导和战术示例。

### Q: 能否支持多人对战？
A: 当前版本仅支持人机对战，可以扩展为双人对战模式。

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！

## 联系方式

如有问题，请提交Issue或联系开发者。
