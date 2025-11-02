# 中文棋谱修复总结

## 修复日期
2025-11-02

## 问题描述

用户反馈了两个问题：

1. **"车二进9"应该是"车二进九"**：红方的棋谱使用了阿拉伯数字，但应该使用中文数字
2. **"马7退8"执行失败**：AI返回的"马7退8"被解析后报"非法移动"

## 问题分析

### 问题1：红方步数使用阿拉伯数字

**现象**：
- 玩家走"车二进九"
- 系统记录为"车二进9"（错误）
- 发送给AI的也是"车二进9"

**原因**：
在`internal/chess/board.go`的`toStandardNotation`方法中，对于纵向移动的步数，使用了：
```go
moveType += strconv.Itoa(abs(startRow - endRow))
```
这会生成阿拉伯数字，没有区分红方和黑方。

**修复**：
```go
// 添加移动步数（红方用中文数字，黑方用阿拉伯数字）
steps := abs(startRow - endRow)
if pieceColor == Red {
    moveType += toChineseNumber(steps)
} else {
    moveType += strconv.Itoa(steps)
}
```

### 问题2："马7退8"解析错误（BUG）

**现象**：
- AI返回：`MOVE: 马7退8`
- 系统报错：`执行中文棋谱失败: 非法移动`

**正确的分析**：
- 黑方马在位置 (2, 6)
- "马7退8"应该是：从 (2, 6) 退到 (0, 7)，吃红方车
- 马腿在位置 (1, 6)，是空位，不会被蹩马腿 ✅
- 这是一个**合法的走法**

**错误的原因**：
在`MoveByChineseNotation`方法的"退"逻辑中，有一个错误的判断：
```go
redGoesUp := (color == Red && fromRow > 5) || (color == Black && fromRow < 5)
```

这个逻辑是错误的。对于黑方：
- 黑方在上方（row 0-4）
- "退"表示向己方方向移动，即向上移动（row减小）
- 但代码中，当黑方马在row=2时，`fromRow < 5`为true，所以`redGoesUp=true`
- 然后代码执行`toRow = fromRow + 2`，即向下移动到(4,7)，这是错误的！
- 正确的应该是`toRow = fromRow - 2`，即向上移动到(0,7)

**修复**：
删除了错误的`redGoesUp`判断，直接根据颜色判断方向：
```go
case "退":
    // 后退
    // 对于直线移动的棋子（车、炮、兵、卒），列不变，数字表示步数
    if pieceType == Chariot || pieceType == Cannon || pieceType == Pawn {
        if color == Red {
            toRow = fromRow + target // 红方向下移动（后退）
        } else {
            toRow = fromRow - target // 黑方向上移动（后退）
        }
        toCol = fromCol
    } else {
        // 对于斜线移动的棋子（马、相、士），数字表示目标列
        toCol, err = b.parseColumn(targetStr, color)
        if err != nil {
            return err
        }
        // 根据目标列计算行（向后移动）
        if pieceType == Horse {
            colDiff := abs(toCol - fromCol)
            if colDiff == 1 {
                if color == Red {
                    toRow = fromRow + 2 // 红方向下移动
                } else {
                    toRow = fromRow - 2 // 黑方向上移动
                }
            } else if colDiff == 2 {
                if color == Red {
                    toRow = fromRow + 1 // 红方向下移动
                } else {
                    toRow = fromRow - 1 // 黑方向上移动
                }
            } else {
                return fmt.Errorf("马的移动不符合日字规则")
            }
        }
        // ... 其他棋子类似
    }
```

**对比**：
- "马7退8"：从(2,6)到(0,7)，马腿在(1,6)空位，✅ 合法
- "马7进8"：从(2,6)到(4,7)，马腿在(3,6)有卒，❌ 非法（蹩马腿）

## 修复内容

### 修改的文件

1. **internal/chess/board.go**
   - 修复`toStandardNotation`方法
   - 确保红方的步数使用中文数字

### 测试验证

```bash
go run test_notation_fix2.go
```

测试结果：
```
第1步 炮二平五 ✅ 成功，棋谱: 炮二平五
第2步 马8进7 ✅ 成功，棋谱: 马8进7
第3步 马二进三 ✅ 成功，棋谱: 马二进三
第4步 炮8平9 ✅ 成功，棋谱: 炮8平9
第5步 车一平二 ✅ 成功，棋谱: 车一平二
第6步 车9平8 ✅ 成功，棋谱: 车9平8
第7步 车二进九 ✅ 成功，棋谱: 车二进九

✅ 所有测试通过！红方使用中文数字，黑方使用阿拉伯数字
```

## 关于"马7退8"的建议

由于"马7退8"是AI选择的非法走法，建议：

1. **短期方案**：当AI返回非法走法时，系统会自动要求AI重新走子（已实现）
2. **长期方案**：优化AI的系统提示词，让AI在选择走法前检查是否合法

## 中文棋谱规则总结

### 红方（使用中文数字）
- 列号：一二三四五六七八九
- 步数：一二三四五六七八九
- 示例：`炮二平五`、`车二进九`、`马二进三`

### 黑方（使用阿拉伯数字）
- 列号：1 2 3 4 5 6 7 8 9
- 步数：1 2 3 4 5 6 7 8 9
- 示例：`炮2平5`、`车9进1`、`马8进7`

### 动作说明
- **平**：横向移动，后面跟目标列
- **进**：向对方方向移动
  - 直线移动棋子（车、炮、兵）：后面跟步数
  - 斜线移动棋子（马、相、士）：后面跟目标列
- **退**：向己方方向移动（规则同"进"）

## 相关文档

- [CHINESE_NOTATION_FIX_COMPLETE.md](CHINESE_NOTATION_FIX_COMPLETE.md) - 之前的修复总结
- [CHESS_COLUMN_REFERENCE.md](CHESS_COLUMN_REFERENCE.md) - 列号参考卡片
- [QUICKSTART.md](QUICKSTART.md) - 快速开始指南
