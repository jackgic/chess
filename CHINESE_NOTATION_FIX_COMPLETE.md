# 中文棋谱解析修复总结

## 修复日期
2025-11-01

## 问题描述

### 问题1：马的"进"和"退"操作解析错误
- **现象**：AI返回 `MOVE: 马8进7`，但执行时报错"非法移动"
- **原因**：对于马来说，"进7"中的"7"应该表示**目标列**，而不是**移动步数**
- **影响**：所有马的走子都无法正确执行

### 问题2：新对局开始时缺少角色说明
- **现象**：新对局开始后，AI不知道自己是哪方，谁先走
- **原因**：首次对话没有告诉AI角色信息
- **影响**：AI可能无法正确理解对局状态

## 修复方案

### 修复1：马的移动规则

#### 修改文件
`internal/chess/board.go` - `MoveByChineseNotation` 方法

#### 修改内容

**原逻辑（错误）**：
```go
case "进":
    // 对于所有棋子，数字都表示步数
    if color == Red {
        toRow = fromRow - target
    } else {
        toRow = fromRow + target
    }
    // 对于马，再解析目标列
    toCol, err = b.parseColumn(targetStr, color)
```

**新逻辑（正确）**：
```go
case "进":
    // 对于直线移动的棋子（车、炮、兵、卒），数字表示步数
    if pieceType == Chariot || pieceType == Cannon || pieceType == Pawn {
        if color == Red {
            toRow = fromRow - target
        } else {
            toRow = fromRow + target
        }
        toCol = fromCol
    } else {
        // 对于斜线移动的棋子（马、相、士），数字表示目标列
        toCol, err = b.parseColumn(targetStr, color)
        if err != nil {
            return err
        }
        // 根据目标列计算行（向前移动）
        if color == Red {
            if pieceType == Horse {
                // 马走日字：如果列差为1，则行差为2；如果列差为2，则行差为1
                colDiff := abs(toCol - fromCol)
                if colDiff == 1 {
                    toRow = fromRow - 2
                } else if colDiff == 2 {
                    toRow = fromRow - 1
                } else {
                    return fmt.Errorf("马的移动不符合日字规则")
                }
            } else {
                // 相、士向前移动1步
                toRow = fromRow - 1
            }
        } else {
            // 黑方向下移动
            if pieceType == Horse {
                colDiff := abs(toCol - fromCol)
                if colDiff == 1 {
                    toRow = fromRow + 2
                } else if colDiff == 2 {
                    toRow = fromRow + 1
                } else {
                    return fmt.Errorf("马的移动不符合日字规则")
                }
            } else {
                toRow = fromRow + 1
            }
        }
    }
```

#### 关键改进
1. **区分棋子类型**：直线移动（车炮兵）vs 斜线移动（马相士）
2. **马的特殊处理**：
   - "马8进7" = 从第8列移动到第7列
   - 根据列差计算行差：列差1→行差2，列差2→行差1（日字规则）
3. **同时修复"退"操作**：使用相同的逻辑

### 修复2：新对局角色说明

#### 修改文件
`internal/game/manager.go` - `AIMove` 方法

#### 修改内容

**原逻辑**：
```go
if len(g.Board.MoveList) == 0 {
    prompt = "开始"
} else {
    lastMove := g.Board.MoveList[len(g.Board.MoveList)-1]
    prompt = lastMove
}
```

**新逻辑**：
```go
if len(g.Board.MoveList) == 0 {
    // 第一步，AI先手（红方）
    if g.AIColor == chess.Red {
        prompt = "你是红方，请先走"
    } else {
        // 这种情况不应该发生，因为如果AI是黑方，玩家应该先走
        return "", fmt.Errorf("AI是黑方但棋盘为空，逻辑错误")
    }
} else {
    // 发送对手刚才的走子（中文棋谱格式）
    lastMove := g.Board.MoveList[len(g.Board.MoveList)-1]
    prompt = lastMove
}
```

#### 关键改进
1. **明确角色**：首次对话告诉AI"你是红方，请先走"
2. **逻辑检查**：如果AI是黑方但棋盘为空，说明逻辑错误
3. **后续对话**：只发送对手的棋谱，保持简洁

## 测试结果

### 测试用例

```go
// 测试1: 红方炮二平五
✅ 成功: 炮二平五

// 测试2: 黑方炮2平5
✅ 成功: 炮2平5

// 测试3: 黑方马8进7
✅ 成功: 马8进7

// 测试4: 红方马二进三
✅ 成功: 马二进三

// 测试5: 完整对局
第1步 红方炮二平五 ✅ 成功
第2步 黑方马8进7 ✅ 成功
第3步 红方马二进三 ✅ 成功
第4步 黑方车9平8 ✅ 成功
```

### 棋谱记录
```
1. 炮二平五
2. 马8进7
3. 马二进三
4. 车9平8
```

## 通信示例

### 场景1：AI是红方（先手）

**第1回合**：
- 发送给AI：`你是红方，请先走`
- AI回复：`MOVE: 炮二平五`
- 系统执行：✅ 成功

**第2回合**：
- 玩家走棋：`马8进7`
- 发送给AI：`马8进7`
- AI回复：`MOVE: 马二进三`
- 系统执行：✅ 成功

### 场景2：AI是黑方（后手）

**第1回合**：
- 玩家走棋：`炮二平五`
- 发送给AI：`炮二平五`
- AI回复：`MOVE: 马8进7`
- 系统执行：✅ 成功

**第2回合**：
- 玩家走棋：`马二进三`
- 发送给AI：`马二进三`
- AI回复：`MOVE: 车9平8`
- 系统执行：✅ 成功

## 中文棋谱规则总结

### 直线移动棋子（车、炮、兵、卒）
- **格式**：`棋子名 + 起始列 + 动作 + 步数`
- **示例**：
  - `炮二平五`：炮从第二列平移到第五列
  - `车9进1`：车从第9列向前移动1步
  - `兵五进一`：兵从第五列向前移动1步

### 斜线移动棋子（马、相、士）
- **格式**：`棋子名 + 起始列 + 动作 + 目标列`
- **示例**：
  - `马8进7`：马从第8列移动到第7列（向前）
  - `马二进三`：马从第二列移动到第三列（向前）
  - `相三进五`：相从第三列移动到第五列（向前）

### 马的日字规则
- **列差1 → 行差2**：`马8进7`（col 7→6，row 0→2）
- **列差2 → 行差1**：`马8进6`（col 7→5，row 0→1）

## 相关文件

### 修改的文件
1. `internal/chess/board.go` - 修复马的移动规则
2. `internal/game/manager.go` - 添加角色说明

### 测试文件
1. `test_ma8_debug.go` - 马8进7调试测试
2. `test_chinese_notation_complete.go` - 完整功能测试

### 文档文件
1. `CHINESE_NOTATION_FIX_SUMMARY.md` - 本文档
2. `CHESS_COLUMN_REFERENCE.md` - 列号参考卡片

## 验证方法

### 1. 运行单元测试
```bash
go run test_chinese_notation_complete.go
```

### 2. 启动服务器测试
```bash
./start.sh
```

### 3. 查看日志
```bash
tail -f logs/chinese-chess-ai.log | grep -E "\[Game\]|\[LKE\]"
```

### 4. 前端测试
1. 打开浏览器访问 `http://localhost:8080`
2. 选择"黑方（后手）"，让AI先走
3. 观察AI的第一步是否正确
4. 走一步"马8进7"，观察是否成功

## 注意事项

1. **马的移动**：必须符合日字规则，否则会报错
2. **回合检查**：必须设置正确的 `board.Turn`
3. **角色一致**：调用 `MoveByChineseNotation` 时，`color` 参数必须与棋子颜色一致
4. **AI提示词**：AI的系统提示词应该包含"使用中文棋谱格式回复，格式为 MOVE: 炮二平五"

## 后续优化建议

1. **支持前后马**：当同一列有两个马时，使用"前马"、"后马"区分
2. **支持前后兵**：当同一列有多个兵时，使用"前兵"、"后兵"区分
3. **错误提示优化**：提供更详细的错误信息，帮助调试
4. **AI回复解析**：支持更多格式，如"我走炮二平五"、"炮二平五"等

## 总结

✅ **问题已完全修复**
- 马的"进"和"退"操作现在可以正确解析
- 新对局开始时会告诉AI角色信息
- 所有测试用例均通过
- 代码已编译成功，可以直接使用

🎉 **现在可以正常与AI进行中国象棋对弈了！**
