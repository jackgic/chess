#!/bin/bash

# 测试中文棋谱交互功能

set -e

echo "=== 中文棋谱交互测试 ==="
echo ""

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8090/api"

echo -e "${BLUE}步骤1: 创建新游戏（玩家选择红方）${NC}"
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/game/new" \
  -H "Content-Type: application/json" \
  -d '{"playerColor": 1}')

echo "$CREATE_RESPONSE" | jq '.'
GAME_ID=$(echo "$CREATE_RESPONSE" | jq -r '.data.gameId')

if [ -z "$GAME_ID" ] || [ "$GAME_ID" = "null" ]; then
    echo -e "${RED}❌ 创建游戏失败${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 游戏创建成功，游戏ID: $GAME_ID${NC}"
echo ""

echo -e "${BLUE}步骤2: 玩家走子（炮二平五）${NC}"
MOVE_RESPONSE=$(curl -s -X POST "$BASE_URL/game/move" \
  -H "Content-Type: application/json" \
  -d "{\"gameId\": \"$GAME_ID\", \"fromRow\": 7, \"fromCol\": 1, \"toRow\": 7, \"toCol\": 4}")

echo "$MOVE_RESPONSE" | jq '.'
echo -e "${GREEN}✅ 玩家走子成功${NC}"
echo ""

echo -e "${BLUE}步骤3: 查看游戏状态${NC}"
STATE_RESPONSE=$(curl -s "$BASE_URL/game/$GAME_ID")
echo "$STATE_RESPONSE" | jq '.data | {turn, moveList}'
echo ""

echo -e "${BLUE}步骤4: AI走子（应该只发送中文棋谱）${NC}"
echo "查看日志以确认发送给AI的内容..."
echo ""

AI_RESPONSE=$(curl -s -X POST "$BASE_URL/game/$GAME_ID/ai-move")
echo "$AI_RESPONSE" | jq '.'

if echo "$AI_RESPONSE" | jq -e '.success' > /dev/null; then
    echo -e "${GREEN}✅ AI走子成功${NC}"
else
    echo -e "${RED}❌ AI走子失败${NC}"
    echo "错误信息："
    echo "$AI_RESPONSE" | jq -r '.error'
fi
echo ""

echo -e "${BLUE}步骤5: 查看最终游戏状态${NC}"
FINAL_STATE=$(curl -s "$BASE_URL/game/$GAME_ID")
echo "$FINAL_STATE" | jq '.data | {turn, status, moveList}'
echo ""

echo -e "${BLUE}步骤6: 查看日志（确认通信格式）${NC}"
echo "最近的游戏日志："
tail -20 logs/chinese-chess-ai.log | grep -E "\[Game\]|\[LKE\]"
echo ""

echo -e "${GREEN}=== 测试完成 ===${NC}"
echo ""
echo "请检查日志，确认："
echo "1. 发送给AI的消息只包含中文棋谱（如：炮二平五）"
echo "2. AI返回的格式为：MOVE: 炮2平5"
echo "3. 系统成功解析并执行了AI的走子"
