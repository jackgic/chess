# 中文棋谱交互 - 快速开始

## 重构说明

已完成与AI智能体通信方式的重构，现在仅使用中文棋谱格式进行交互。

## 核心变化

### 通信方式

**之前：**
```
发送给AI: "对手走了：炮二平五，现在轮到你了，请走子。"
AI回复: "MOVE: 71-74"
```

**现在：**
```
发送给AI: "炮二平五"
AI回复: "MOVE: 炮2平5"
```

### AI智能体提示词

确保在LKE平台配置的系统提示词为：

```
你是一个专业的中国象棋选手，棋力非凡，使用棋谱表示方法进行对弈

# 回答案例
MOVE: 炮2平5
MOVE: 马二进三
```

## 快速测试

### 1. 启动服务器

```bash
./start.sh
```

### 2. 运行测试脚本

```bash
./test_chinese_notation.sh
```

### 3. 查看日志

```bash
tail -f logs/chinese-chess-ai.log | grep -E "\[Game\]|\[LKE\]"
```

## 验证要点

在日志中确认：

1. ✅ **发送给AI的消息**：只包含中文棋谱
   ```
   [Game] 发送棋谱给AI: 炮二平五
   ```

2. ✅ **AI的回复**：包含MOVE格式的中文棋谱
   ```
   [Game] AI回复: MOVE: 炮2平5
   ```

3. ✅ **解析成功**：系统成功提取中文棋谱
   ```
   [LKE] 提取到的中文棋谱: 炮2平5
   ```

4. ✅ **执行成功**：系统成功执行走子
   ```
   [Game] 走子执行成功
   ```

## 手动测试

### 创建游戏

```bash
curl -X POST http://localhost:8090/api/game/new \
  -H "Content-Type: application/json" \
  -d '{"playerColor": 1}'
```

### 玩家走子（炮二平五）

```bash
curl -X POST http://localhost:8090/api/game/move \
  -H "Content-Type: application/json" \
  -d '{"gameId": "game_1", "fromRow": 7, "fromCol": 1, "toRow": 7, "toCol": 4}'
```

### AI走子

```bash
curl -X POST http://localhost:8090/api/game/game_1/ai-move
```

### 查看游戏状态

```bash
curl http://localhost:8090/api/game/game_1 | jq '.data.moveList'
```

## 支持的棋谱格式

### 棋子名称
- 红方：帅、仕、相、马、车、炮、兵
- 黑方：将、士、象、马、车、炮、卒

### 列号表示
- 红方：1-9（从右到左）或 一-九
- 黑方：1-9（从左到右）或 一-九

### 动作
- `平`：横向移动
- `进`：向对方方向移动
- `退`：向己方方向移动

### 示例
- `炮2平5`：炮从第2列平移到第5列
- `马二进三`：马从第二列进到第三列
- `车9进1`：车从第9列向前进1步

## 常见问题

### Q1: AI返回的格式不对怎么办？

**A:** 检查LKE平台的系统提示词配置，确保包含正确的回答案例。

### Q2: 解析中文棋谱失败？

**A:** 查看日志中的AI完整回复，确认是否包含 `MOVE:` 前缀和正确的中文棋谱格式。

### Q3: 执行走子失败？

**A:** 可能的原因：
- 棋子位置不对（同一列有多个相同棋子）
- 走子不合法（违反象棋规则）
- 列号表示错误（红方和黑方的列号方向不同）

## 技术细节

详细的技术说明请参考：[CHINESE_NOTATION_REFACTOR.md](CHINESE_NOTATION_REFACTOR.md)

## 下一步

1. 测试更多对局场景
2. 优化同列多子的处理
3. 增加错误提示和恢复机制
4. 性能优化
