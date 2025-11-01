# 中文棋谱交互重构说明

## 重构目标

简化与AI智能体的通信方式，改为仅使用中文棋谱格式进行交互。

## 重构内容

### 1. 通信方式变更

**重构前：**
- 发送给AI：完整的提示词，包含棋盘状态、规则说明等
- AI返回：`MOVE: 数字坐标格式`（如 `MOVE: 71-74`）

**重构后：**
- 发送给AI：仅发送中文棋谱（如 `炮二平五`）
- AI返回：`MOVE: 中文棋谱格式`（如 `MOVE: 炮2平5` 或 `MOVE: 炮二平五`）

### 2. 代码修改

#### 2.1 游戏管理器 (`internal/game/manager.go`)

修改 `AIMove()` 方法：

```go
// 重构前
var prompt string
if len(g.Board.MoveList) == 0 {
    prompt = "游戏开始，你是先手，请走第一步。"
} else {
    lastMove := g.Board.MoveList[len(g.Board.MoveList)-1]
    prompt = fmt.Sprintf("对手走了：%s，现在轮到你了，请走子。", lastMove)
}

// 重构后
var prompt string
if len(g.Board.MoveList) == 0 {
    prompt = "开始"
} else {
    // 只发送对手刚才的走子（中文棋谱格式）
    lastMove := g.Board.MoveList[len(g.Board.MoveList)-1]
    prompt = lastMove
}
```

#### 2.2 LKE客户端 (`internal/lke/client.go`)

新增 `ExtractChineseMove()` 方法：

```go
// ExtractChineseMove 从AI回答中提取中文棋谱
func (c *Client) ExtractChineseMove(answer string) (string, error) {
    // 匹配中文象棋记谱格式 MOVE: 炮2平5 或 MOVE: 炮二平五
    chineseMovePattern := regexp.MustCompile(`MOVE:\s*([炮车马相象士将帅兵卒])([\d一二三四五六七八九])([进退平])([\d一二三四五六七八九])`)
    allChineseMatches := chineseMovePattern.FindAllStringSubmatch(answer, -1)
    if len(allChineseMatches) > 0 {
        matches := allChineseMatches[len(allChineseMatches)-1]
        if len(matches) >= 5 {
            chineseMove := matches[1] + matches[2] + matches[3] + matches[4]
            return chineseMove, nil
        }
    }
    return "", fmt.Errorf("无法从回答中提取中文棋谱")
}
```

#### 2.3 棋盘 (`internal/chess/board.go`)

新增 `MoveByChineseNotation()` 方法：

```go
// MoveByChineseNotation 根据中文棋谱执行走子
// 例如：炮2平5、马二进三、车9进1
func (b *Board) MoveByChineseNotation(notation string, color Color) error {
    // 解析中文棋谱
    // 1. 解析棋子名称
    // 2. 解析起始列
    // 3. 解析动作（进/退/平）
    // 4. 解析目标列或步数
    // 5. 查找棋子位置
    // 6. 计算目标位置
    // 7. 执行走子
}
```

支持的辅助方法：
- `parsePieceType()`: 解析棋子类型
- `parseColumn()`: 解析列号（支持红方和黑方的不同表示）
- `parseNumber()`: 解析数字（支持阿拉伯数字和中文数字）
- `findPiece()`: 查找指定棋子的位置

### 3. AI智能体提示词配置

在LKE平台配置的系统提示词应该包含：

```
你是一个专业的中国象棋选手，棋力非凡，使用棋谱表示方法进行对弈

# 回答案例
MOVE: 炮2平5
MOVE: 马二进三

# 规则说明
- 你会收到对手的走子（中文棋谱格式）
- 你需要回复你的走子，格式为：MOVE: 棋子名+起始列+动作+目标列/步数
- 支持阿拉伯数字（1-9）或中文数字（一-九）
- 动作包括：进、退、平
- 例如：炮2平5、马二进三、车9进1
```

## 使用示例

### 游戏流程

1. **创建游戏**
   - 系统创建新的会话ID
   - 初始化棋盘

2. **玩家走子**
   - 玩家走：炮二平五
   - 系统记录：`炮二平五`

3. **AI走子**
   - 发送给AI：`炮二平五`
   - AI回复：`MOVE: 炮2平5`
   - 系统解析：`炮2平5`
   - 系统执行：将黑方炮从第2列平移到第5列
   - 系统记录：`炮2平5`

4. **继续对弈**
   - 玩家走：马二进三
   - 发送给AI：`马二进三`
   - AI回复：`MOVE: 马8进7`
   - 系统执行并记录

## 测试方法

### 方法1: 使用测试脚本

```bash
# 启动服务器
./start.sh

# 运行测试
./test_chinese_notation.sh
```

### 方法2: 手动测试

```bash
# 1. 创建游戏
curl -X POST http://localhost:8090/api/game/new \
  -H "Content-Type: application/json" \
  -d '{"playerColor": 1}'

# 2. 玩家走子
curl -X POST http://localhost:8090/api/game/move \
  -H "Content-Type: application/json" \
  -d '{"gameId": "game_1", "fromRow": 7, "fromCol": 1, "toRow": 7, "toCol": 4}'

# 3. AI走子
curl -X POST http://localhost:8090/api/game/game_1/ai-move

# 4. 查看日志
tail -f logs/chinese-chess-ai.log
```

### 验证要点

查看日志，确认：

1. ✅ 发送给AI的消息只包含中文棋谱（如：`炮二平五`）
2. ✅ AI返回格式为：`MOVE: 炮2平5`
3. ✅ 系统成功解析并执行了AI的走子
4. ✅ 棋谱记录正确

## 优势

1. **简化通信**：每次只发送一个中文棋谱字符串，极大减少通信内容
2. **符合习惯**：使用标准的中国象棋记谱法，更加直观
3. **降低成本**：减少token消耗
4. **易于调试**：通信内容简洁明了，便于排查问题
5. **AI友好**：AI智能体已经配置好提示词，可以直接理解中文棋谱

## 注意事项

1. **列号表示**：
   - 红方：从右到左为1-9（或一-九）
   - 黑方：从左到右为1-9

2. **动作说明**：
   - `平`：横向移动，行不变
   - `进`：向对方方向移动（红方向上，黑方向下）
   - `退`：向己方方向移动（红方向下，黑方向上）

3. **数字格式**：
   - 支持阿拉伯数字：1-9
   - 支持中文数字：一、二、三、四、五、六、七、八、九

4. **棋子查找**：
   - 系统会在指定列上查找对应颜色和类型的棋子
   - 如果同一列有多个相同棋子，会选择第一个找到的

## 后续优化建议

1. **处理同列多子**：当同一列有多个相同棋子时，需要更精确的定位方式
2. **错误处理**：增加更详细的错误提示，帮助AI纠正错误
3. **棋谱验证**：在执行前验证中文棋谱的合法性
4. **性能优化**：缓存棋子位置，避免每次都遍历整个棋盘
