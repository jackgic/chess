#!/bin/bash

# 测试会话式对话 - 验证AI只接收简洁的走子信息

echo "=========================================="
echo "测试会话式AI对话"
echo "=========================================="
echo ""

# 服务地址
BASE_URL="http://localhost:8080"

# 1. 检查服务状态
echo "1. 检查服务状态..."
if ! curl -s "$BASE_URL/health" > /dev/null; then
    echo "❌ 服务未启动，请先运行: ./start.sh"
    exit 1
fi
echo "✅ 服务正常运行"
echo ""

# 2. 创建新游戏
echo "2. 创建新游戏..."
GAME_RESPONSE=$(curl -s -X POST "$BASE_URL/game/new" \
    -H "Content-Type: application/json" \
    -d '{"player_color": "red"}')

GAME_ID=$(echo $GAME_RESPONSE | grep -o '"game_id":"[^"]*"' | cut -d'"' -f4)
SESSION_ID=$(echo $GAME_RESPONSE | grep -o '"session_id":"[^"]*"' | cut -d'"' -f4)

if [ -z "$GAME_ID" ]; then
    echo "❌ 创建游戏失败"
    echo "响应: $GAME_RESPONSE"
    exit 1
fi

echo "✅ 游戏创建成功"
echo "   游戏ID: $GAME_ID"
echo "   会话ID: $SESSION_ID"
echo ""

# 3. 玩家走第一步（炮二平五）
echo "3. 玩家走第一步: 炮二平五 (从7,1到7,4)..."
MOVE_RESPONSE=$(curl -s -X POST "$BASE_URL/game/$GAME_ID/move" \
    -H "Content-Type: application/json" \
    -d '{
        "from_row": 7,
        "from_col": 1,
        "to_row": 7,
        "to_col": 4
    }')

echo "响应: $MOVE_RESPONSE"
echo ""

# 4. 查看日志中AI收到的提示词
echo "4. 查看日志中AI收到的提示词..."
echo "----------------------------------------"
tail -20 logs/chinese-chess-ai.log | grep -A 5 "发送提示词给AI"
echo "----------------------------------------"
echo ""

# 5. AI走子
echo "5. 请求AI走子..."
AI_RESPONSE=$(curl -s -X POST "$BASE_URL/game/$GAME_ID/ai-move")

echo "AI响应: $AI_RESPONSE"
echo ""

# 6. 再次查看日志
echo "6. 查看AI走子后的日志..."
echo "----------------------------------------"
tail -30 logs/chinese-chess-ai.log | grep -A 5 "发送提示词给AI"
echo "----------------------------------------"
echo ""

# 7. 玩家再走一步
echo "7. 玩家再走一步: 马二进三 (从7,1到5,2)..."
MOVE_RESPONSE2=$(curl -s -X POST "$BASE_URL/game/$GAME_ID/move" \
    -H "Content-Type: application/json" \
    -d '{
        "from_row": 7,
        "from_col": 1,
        "to_row": 5,
        "to_col": 2
    }')

echo "响应: $MOVE_RESPONSE2"
echo ""

# 8. AI再次走子
echo "8. AI再次走子..."
AI_RESPONSE2=$(curl -s -X POST "$BASE_URL/game/$GAME_ID/ai-move")

echo "AI响应: $AI_RESPONSE2"
echo ""

# 9. 最终查看日志，验证提示词简洁性
echo "9. 验证提示词简洁性..."
echo "=========================================="
echo "查看最近3次发送给AI的提示词："
echo "=========================================="
tail -50 logs/chinese-chess-ai.log | grep -B 1 "发送提示词给AI"
echo ""

echo "=========================================="
echo "测试完成！"
echo "=========================================="
echo ""
echo "✅ 预期结果："
echo "   - 第一次AI走子：提示词应该是 '游戏开始，你是先手，请走第一步。'"
echo "   - 第二次AI走子：提示词应该是 '对手走了：炮二平五，现在轮到你了，请走子。'"
echo "   - 第三次AI走子：提示词应该是 '对手走了：马二进三，现在轮到你了，请走子。'"
echo ""
echo "📝 如果提示词仍然很长，说明优化未生效"
echo "📝 如果提示词简洁，说明优化成功！"
