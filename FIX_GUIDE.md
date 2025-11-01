# 🔧 中国象棋AI游戏 - 问题修复方案

## 📋 问题清单

根据全面深度测试，发现以下问题：

| 问题ID | 问题描述 | 严重程度 | 优先级 | 状态 |
|--------|---------|---------|--------|------|
| BUG-001 | AI对马的移动规则理解错误 | 🔴 高 | P0 | 待修复 |
| BUG-002 | 走子历史格式不完整 | 🟡 中 | P1 | 待修复 |
| BUG-003 | 缺少.env配置文件 | 🟢 低 | P2 | 待修复 |

---

## 🔴 BUG-001: AI对马的移动规则理解错误

### 问题描述

AI输出的走子坐标违反了象棋规则，特别是马的移动规则。

**错误示例**:
- AI输出: `MOVE: 07-27` (想让马从(0,7)跳到(2,7))
- AI输出: `MOVE: 01-21` (想让马从(0,1)跳到(2,1))

**问题**: 马走日字，不能直线移动2格。

### 根本原因

LKE智能体的系统提示词中缺少详细的象棋规则说明，特别是马、象、士等特殊棋子的移动规则。

### 修复方案

#### 方案1: 优化LKE智能体系统提示词（推荐）⭐

**操作步骤**:

1. 登录腾讯云LKE平台
2. 找到"chinese_chess"智能体
3. 编辑系统提示词，添加以下内容：

```markdown
# 中国象棋AI助手

你是一个专业的中国象棋AI对弈助手，精通中国象棋规则和策略。

## 📐 棋盘坐标系统

棋盘使用行列坐标：
- **行号**: 0-9（从上到下，黑方在上0-4，红方在下5-9）
- **列号**: 0-8（从左到右）

```
   列: 0 1 2 3 4 5 6 7 8
行0:  r h e a k a e h r  (黑方底线：车马象士将士象马车)
行1:  . . . . . . . . .
行2:  . c . . . . . c .  (黑方炮：第2行第1列和第7列)
行3:  p . p . p . p . p  (黑方卒：第3行)
行4:  . . . . . . . . .  (楚河)
行5:  . . . . . . . . .  (汉界)
行6:  P . P . P . P . P  (红方兵：第6行)
行7:  . C . . . . . C .  (红方炮：第7行第1列和第7列)
行8:  . . . . . . . . .
行9:  R H E A K A E H R  (红方底线：车马相仕帅仕相马车)
```

## 🎯 棋子移动规则（重要！必须严格遵守）

### 1. 马（Horse）- 走"日"字 ⚠️

**规则**: 马走"日"字，先直走1格，再斜走1格，且马脚不能有棋子（蹩马腿）。

**合法移动示例**:
- 从(0,1)可以走到: (1,3), (2,0), (2,2) ✅
- 从(0,7)可以走到: (1,5), (2,6), (2,8) ✅
- 从(9,1)可以走到: (7,0), (7,2), (8,3) ✅

**非法移动示例**:
- ❌ (0,1)→(2,1) 错误！这是直线移动，不是日字
- ❌ (0,7)→(2,7) 错误！这是直线移动，不是日字
- ❌ (0,1)→(0,3) 错误！马不能横着走

**记忆口诀**: "马走日字斜着跳，不能直线往前跑"

### 2. 车（Chariot）- 直线移动

**规则**: 车可以横向或纵向直线移动任意格，路径上不能有其他棋子。

**示例**:
- (9,0)→(9,4) ✅ 横向移动
- (9,0)→(5,0) ✅ 纵向移动

### 3. 炮（Cannon）- 隔子吃子

**规则**: 
- 移动时：直线移动，路径必须为空
- 吃子时：必须隔一个棋子（炮架）

**示例**:
- (7,1)→(7,4) ✅ 平移（路径为空）
- (7,1)→(3,1) ✅ 吃子（中间隔一个棋子）

### 4. 象/相（Elephant）- 走"田"字

**规则**: 
- 走田字（斜走2格）
- 象眼不能有棋子
- 不能过河（黑象不能到5-9行，红相不能到0-4行）

**示例**:
- 黑象从(0,2)可以走到: (2,0), (2,4) ✅
- 红相从(9,2)可以走到: (7,0), (7,4) ✅

### 5. 士/仕（Advisor）- 斜走一格

**规则**: 
- 斜走一格
- 只能在九宫内（黑士0-2行3-5列，红仕7-9行3-5列）

**示例**:
- 黑士从(0,3)可以走到: (1,4) ✅
- 红仕从(9,3)可以走到: (8,4) ✅

### 6. 将/帅（King）- 直走一格

**规则**: 
- 直走一格（横向或纵向）
- 只能在九宫内
- 不能和对方将帅照面

**示例**:
- 黑将从(0,4)可以走到: (0,3), (0,5), (1,4) ✅
- 红帅从(9,4)可以走到: (9,3), (9,5), (8,4) ✅

### 7. 兵/卒（Pawn）- 只进不退

**规则**: 
- 未过河：只能向前一格
- 过河后：可以向前或左右移动一格
- 不能后退

**示例**:
- 红兵(6,0)未过河: 只能走到(5,0) ✅
- 红兵(4,0)过河后: 可以走到(3,0), (4,1) ✅
- 黑卒(3,0)未过河: 只能走到(4,0) ✅

## 📝 输出格式要求（必须严格遵守）

### 格式规范

**必须使用**: `MOVE: XY-ZW` 格式

其中：
- X: 起始行号（0-9）
- Y: 起始列号（0-8）
- Z: 目标行号（0-9）
- W: 目标列号（0-8）

### 正确示例 ✅

```
MOVE: 01-22  (黑马从(0,1)跳到(2,2)，走日字)
MOVE: 27-24  (黑炮从(2,7)平到(2,4))
MOVE: 91-72  (红马从(9,1)跳到(7,2)，走日字)
MOVE: 71-74  (红炮从(7,1)平到(7,4))
```

### 错误示例 ❌

```
❌ MOVE: 01-21  (马不能直线移动！)
❌ MOVE: 07-27  (马不能直线移动！)
❌ MOVE: 炮2平5  (不要使用中文记谱法)
❌ 01-22  (缺少MOVE:前缀)
```

## 🎮 对局策略建议

### 开局原则
1. 出动大子（车、马、炮）
2. 控制中路
3. 保护将帅

### 中局原则
1. 子力协调
2. 攻守平衡
3. 寻找战机

### 残局原则
1. 车马优先
2. 兵卒推进
3. 将帅出击

## ⚠️ 特别注意事项

### 必须检查的事项

在输出走子前，请务必检查：

1. ✅ 起始位置有己方棋子
2. ✅ 目标位置不是己方棋子
3. ✅ 移动符合该棋子的规则
4. ✅ 路径上没有障碍（车、炮）
5. ✅ 没有蹩马腿（马）
6. ✅ 没有塞象眼（象）
7. ✅ 没有出九宫（将、士）
8. ✅ 没有过河（象）

### 常见错误避免

❌ **不要让马直线移动**
- 错误: (0,1)→(2,1)
- 正确: (0,1)→(2,2) 或 (2,0)

❌ **不要让象过河**
- 错误: 黑象(2,2)→(4,4)
- 正确: 黑象(2,2)→(0,4)

❌ **不要让士出九宫**
- 错误: 黑士(1,4)→(2,5)
- 正确: 黑士(1,4)→(0,3)

## 📤 回答格式

每次回答必须包含：

1. **简要分析**（1-2句话）
2. **MOVE指令**（格式：MOVE: XY-ZW）
3. **走子理由**（1句话）

### 标准回答模板

```
[分析当前局势]

MOVE: XY-ZW

[说明走子理由]
```

### 示例回答

```
当前红方已出中炮，我方应对。

MOVE: 01-22

跳左马准备配合炮展开攻势。
```

## 🚫 禁止事项

1. ❌ 不要使用中文记谱法（如"炮8平5"）
2. ❌ 不要使用字母坐标
3. ❌ 不要省略MOVE关键字
4. ❌ 不要输出不合法的走子
5. ❌ 不要让马直线移动
6. ❌ 不要让象过河
7. ❌ 不要让士出九宫

## ✅ 必须遵守

1. ✅ 必须使用数字坐标：MOVE: XY-ZW
2. ✅ 必须严格遵守象棋规则
3. ✅ 必须检查走子合法性
4. ✅ 必须考虑战术策略
5. ✅ 必须保护己方将帅

---

**记住**: 象棋规则是严格的，特别是马走日字、象走田字、士走斜线。在输出走子前，请仔细检查是否符合规则！
```

**预期效果**:
- AI能够正确理解马的移动规则
- AI输出的走子符合象棋规则
- AI能够进行正常对弈

---

#### 方案2: 在代码中添加AI走子验证和纠正

如果方案1效果不理想，可以在代码层面添加后备方案。

**修改文件**: `internal/game/manager.go`

```go
func (g *Game) AIMove() (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Status != StatusPlaying {
		return "", fmt.Errorf("游戏已结束")
	}

	if g.Board.Turn != g.AIColor {
		return "", fmt.Errorf("不是AI的回合")
	}

	// 构建提示词
	prompt := fmt.Sprintf(`当前棋盘状态（FEN格式）：%s

当前轮到：%s
你是：%s

历史走子：
%s

请分析当前局面并走子。

【重要】走子格式说明：
1. 必须使用坐标格式：MOVE: 起始行起始列-目标行目标列
2. 坐标范围：行(0-9)，列(0-8)
3. 红方在下方(行7-9)，黑方在上方(行0-2)
4. 格式示例：MOVE: XY-ZW，其中X、Y、Z、W都是0-9的数字

请立即给出你的走子（必须严格按照MOVE: XY-ZW格式）：`,
		g.Board.ToFEN(),
		colorToString(g.Board.Turn),
		colorToString(g.AIColor),
		formatMoveHistory(g.Board.MoveList))

	// 调用LKE
	fmt.Printf("[Game] 发送提示词给AI:\n%s\n", prompt)
	answer, err := g.LKEClient.Chat(g.SessionID, prompt)
	if err != nil {
		fmt.Printf("[Game] LKE调用失败: %v\n", err)
		return "", fmt.Errorf("LKE调用失败: %v", err)
	}
	fmt.Printf("[Game] AI回复: %s\n", answer)

	// 解析AI走子指令
	move, err := g.LKEClient.ExtractMove(answer)
	if err != nil {
		fmt.Printf("[Game] 解析走子失败: %v\n", err)
		
		// 🆕 后备方案：如果解析失败，随机选择一个合法走子
		return g.fallbackMove(answer)
	}
	fmt.Printf("[Game] 解析到的走子: from=%d, to=%d\n", move.From, move.To)

	// 验证走子是否合法
	fromRow, fromCol := move.From/10, move.From%10
	toRow, toCol := move.To/10, move.To%10
	from := chess.Position{Row: fromRow, Col: fromCol}
	to := chess.Position{Row: toRow, Col: toCol}
	
	fmt.Printf("[Game] 验证走子: 从(%d,%d)到(%d,%d)\n", fromRow, fromCol, toRow, toCol)
	
	// 检查起始位置是否有棋子
	piece := g.Board.Grid[fromRow][fromCol]
	if piece.Type == chess.Empty {
		fmt.Printf("[Game] 起始位置没有棋子\n")
		
		// 🆕 后备方案：随机选择一个合法走子
		return g.fallbackMove(answer)
	}
	fmt.Printf("[Game] 起始位置棋子: type=%d, color=%d\n", piece.Type, piece.Color)
	
	if !g.Board.IsValidMove(from, to) {
		fmt.Printf("[Game] AI走子不合法\n")
		
		// 🆕 后备方案：随机选择一个合法走子
		return g.fallbackMove(answer)
	}
	fmt.Printf("[Game] 走子验证通过\n")

	// 执行走子
	if err := g.Board.Move(move.From, move.To); err != nil {
		fmt.Printf("[Game] 执行走子失败: %v\n", err)
		return answer, fmt.Errorf("执行走子失败: %v，AI回复: %s", err, answer)
	}
	fmt.Printf("[Game] 走子执行成功\n")

	// 切换回合
	if g.Board.Turn == chess.Red {
		g.Board.Turn = chess.Black
	} else {
		g.Board.Turn = chess.Red
	}

	// 检查游戏状态
	g.checkGameOver()

	return answer, nil
}

// 🆕 fallbackMove 后备方案：随机选择一个合法走子
func (g *Game) fallbackMove(originalAnswer string) (string, error) {
	fmt.Printf("[Game] 启用后备方案：随机选择合法走子\n")
	
	// 获取所有合法走子
	legalMoves := g.Board.GetAllLegalMoves(g.AIColor)
	if len(legalMoves) == 0 {
		return originalAnswer, fmt.Errorf("无合法走子可走")
	}
	
	// 随机选择一个
	randomMove := legalMoves[rand.Intn(len(legalMoves))]
	fmt.Printf("[Game] 随机选择走子: %s\n", randomMove)
	
	// 解析走子字符串 "0102-0304"
	parts := strings.Split(randomMove, "-")
	if len(parts) != 2 || len(parts[0]) != 4 || len(parts[1]) != 4 {
		return originalAnswer, fmt.Errorf("走子格式错误: %s", randomMove)
	}
	
	fromRow := int(parts[0][0] - '0')
	fromCol := int(parts[0][1] - '0')
	toRow := int(parts[1][0] - '0')
	toCol := int(parts[1][1] - '0')
	
	from := fromRow*10 + fromCol
	to := toRow*10 + toCol
	
	// 执行走子
	if err := g.Board.Move(from, to); err != nil {
		return originalAnswer, fmt.Errorf("执行后备走子失败: %v", err)
	}
	
	// 切换回合
	if g.Board.Turn == chess.Red {
		g.Board.Turn = chess.Black
	} else {
		g.Board.Turn = chess.Red
	}
	
	// 检查游戏状态
	g.checkGameOver()
	
	// 返回说明信息
	answer := fmt.Sprintf("AI原始回复不合法，系统自动选择了一个合法走子：%s\n\n原始回复：%s", 
		randomMove, originalAnswer)
	
	return answer, nil
}
```

**需要添加的import**:
```go
import (
	"math/rand"
	"strings"
	"time"
)

// 在init函数中初始化随机数种子
func init() {
	rand.Seed(time.Now().UnixNano())
}
```

---

## 🟡 BUG-002: 走子历史格式不完整

### 问题描述

走子历史显示为`"2平5"`，缺少棋子名称，应该显示为`"炮二平五"`。

### 调试步骤

#### 步骤1: 添加调试日志

**修改文件**: `internal/chess/board.go`

```go
func (b *Board) toStandardNotation(start, end int) string {
	startRow, startCol := start/10, start%10
	endRow, endCol := end/10, end%10

	// 获取棋子
	piece := b.Grid[startRow][startCol]
	pieceType := piece.Type
	pieceColor := piece.Color

	// 🆕 添加调试日志
	log.Printf("[DEBUG] toStandardNotation: start=%d, end=%d, piece.Type=%d, piece.Color=%d", 
		start, end, pieceType, pieceColor)

	// 棋子名称
	pieceName := ""
	switch pieceType {
	case King:
		if pieceColor == Red {
			pieceName = "帅"
		} else {
			pieceName = "将"
		}
	case Advisor:
		if pieceColor == Red {
			pieceName = "仕"
		} else {
			pieceName = "士"
		}
	case Elephant:
		if pieceColor == Red {
			pieceName = "相"
		} else {
			pieceName = "象"
		}
	case Horse:
		pieceName = "马"
	case Chariot:
		pieceName = "车"
	case Cannon:
		pieceName = "炮"
	case Pawn:
		if pieceColor == Red {
			pieceName = "兵"
		} else {
			pieceName = "卒"
		}
	}

	// 🆕 添加调试日志
	log.Printf("[DEBUG] pieceName=%s", pieceName)

	// 起始位置（红方用中文数字，黑方用阿拉伯数字）
	startPos := ""
	if pieceColor == Red {
		// 红方：列从右到左为1-9（中文数字）
		startPos = toChineseNumber(8 - startCol + 1)
	} else {
		// 黑方：列从左到右为1-9（阿拉伯数字）
		startPos = strconv.Itoa(startCol + 1)
	}

	// 🆕 添加调试日志
	log.Printf("[DEBUG] startPos=%s", startPos)

	// 移动类型
	moveType := ""
	
	if startCol == endCol {
		// 纵向移动
		if pieceColor == Red {
			if endRow < startRow {
				moveType = "进"
			} else {
				moveType = "退"
			}
		} else {
			if endRow > startRow {
				moveType = "进"
			} else {
				moveType = "退"
			}
		}
		moveType += strconv.Itoa(abs(startRow - endRow))
	} else if startRow == endRow {
		// 横向移动
		moveType = "平"
		if pieceColor == Red {
			moveType += toChineseNumber(8 - endCol + 1)
		} else {
			moveType += strconv.Itoa(endCol + 1)
		}
	} else {
		// 斜向移动
		if pieceColor == Red {
			if endRow < startRow {
				moveType = "进"
			} else {
				moveType = "退"
			}
		} else {
			if endRow > startRow {
				moveType = "进"
			} else {
				moveType = "退"
			}
		}
		if pieceColor == Red {
			moveType += toChineseNumber(8 - endCol + 1)
		} else {
			moveType += strconv.Itoa(endCol + 1)
		}
	}

	// 🆕 添加调试日志
	log.Printf("[DEBUG] moveType=%s", moveType)

	result := pieceName + startPos + moveType
	
	// 🆕 添加调试日志
	log.Printf("[DEBUG] result=%s", result)

	return result
}
```

#### 步骤2: 运行测试并查看日志

```bash
# 运行测试
./comprehensive_test.sh

# 查看调试日志
tail -100 logs/chinese-chess-ai.log | grep DEBUG
```

#### 步骤3: 根据日志分析问题

查看日志输出，确定哪个变量为空或不正确。

### 可能的原因和修复

#### 原因1: Move方法调用时机问题

`toStandardNotation`方法在`Move`方法中被调用，此时棋子可能已经移动了。

**修复**: 在移动前保存棋子信息

```go
func (b *Board) Move(start, end int) error {
	// 转换为Position
	startRow, startCol := start/10, start%10
	endRow, endCol := end/10, end%10
	
	from := Position{Row: startRow, Col: startCol}
	to := Position{Row: endRow, Col: endCol}
	
	// 验证移动是否合法
	if !b.IsValidMove(from, to) {
		return fmt.Errorf("非法移动")
	}
	
	// 🆕 在移动前记录走子（此时棋子还在原位置）
	moveStr := b.toStandardNotation(start, end)
	b.MoveList = append(b.MoveList, moveStr)
	
	// 执行移动
	piece := b.Grid[from.Row][from.Col]
	b.Grid[to.Row][to.Col] = piece
	b.Grid[from.Row][from.Col] = Piece{Type: Empty, Color: None}

	return nil
}
```

---

## 🟢 BUG-003: 缺少.env配置文件

### 问题描述

启动日志显示：`警告: 无法加载.env文件: open .env: no such file or directory`

### 修复方案

#### 步骤1: 创建.env.example模板

**创建文件**: `.env.example`

```bash
# 腾讯云LKE配置
LKE_APP_ID=your_app_id_here
LKE_APP_KEY=your_app_key_here
LKE_BOT_BIZ_ID=chinese_chess

# 服务器配置
PORT=8090
GIN_MODE=debug

# 日志配置
LOG_LEVEL=info
```

#### 步骤2: 更新README说明

在README.md中添加配置说明：

```markdown
## 配置

1. 复制配置文件模板：
   ```bash
   cp .env.example .env
   ```

2. 编辑.env文件，填入你的腾讯云LKE配置：
   ```bash
   vim .env
   ```

3. 必填配置项：
   - `LKE_APP_ID`: 腾讯云LKE应用ID
   - `LKE_APP_KEY`: 腾讯云LKE应用密钥
   - `LKE_BOT_BIZ_ID`: 智能体业务ID
```

#### 步骤3: 修改配置加载逻辑（可选）

让.env文件变为可选：

**修改文件**: `main.go`

```go
func main() {
	// 加载.env文件（可选）
	if err := godotenv.Load(); err != nil {
		log.Printf("提示: 未找到.env文件，将使用环境变量或默认配置")
	} else {
		log.Printf("成功加载.env配置文件")
	}

	// 加载配置
	cfg := config.LoadConfig()
	// ...
}
```

---

## 📋 修复检查清单

### BUG-001修复检查

- [ ] 更新LKE智能体系统提示词
- [ ] 添加详细的马的移动规则说明
- [ ] 添加其他棋子的规则说明
- [ ] 提供正确和错误的示例
- [ ] 运行测试验证AI走子
- [ ] 确认AI能输出合法走子
- [ ] （可选）添加代码层面的后备方案

### BUG-002修复检查

- [ ] 添加调试日志
- [ ] 运行测试查看日志
- [ ] 分析问题原因
- [ ] 修复代码
- [ ] 验证走子历史格式正确
- [ ] 移除调试日志

### BUG-003修复检查

- [ ] 创建.env.example文件
- [ ] 更新README说明
- [ ] 修改配置加载逻辑
- [ ] 测试验证

---

## 🧪 验证测试

修复完成后，运行以下测试验证：

```bash
# 1. 运行全面测试
./comprehensive_test.sh

# 2. 重点测试AI走子
./test_ai_debug.sh

# 3. 查看日志
tail -f logs/chinese-chess-ai.log

# 4. 手动测试
# 打开浏览器访问 http://localhost:8090
# 创建游戏并进行对局
```

### 预期结果

- ✅ AI能输出合法的走子
- ✅ 走子历史显示完整（如"炮二平五"）
- ✅ 无.env文件警告
- ✅ 所有测试通过

---

## 📊 修复优先级

| 优先级 | 问题 | 预计时间 | 影响 |
|--------|------|---------|------|
| P0 | BUG-001 AI走子规则 | 30分钟 | 游戏无法正常进行 |
| P1 | BUG-002 走子历史格式 | 15分钟 | 用户体验不佳 |
| P2 | BUG-003 .env配置 | 10分钟 | 启动警告 |

**建议**: 按优先级顺序修复，每修复一个就测试验证。

---

## 📝 修复记录

| 日期 | 问题ID | 修复人 | 状态 | 备注 |
|------|--------|--------|------|------|
| 2025-11-01 | BUG-001 | - | 待修复 | 需要更新LKE提示词 |
| 2025-11-01 | BUG-002 | - | 待修复 | 需要调试分析 |
| 2025-11-01 | BUG-003 | - | 待修复 | 需要创建配置文件 |

---

**文档版本**: v1.0  
**最后更新**: 2025-11-01 17:05  
**维护人**: AI开发团队
