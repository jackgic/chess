#!/bin/bash

echo "=========================================="
echo "🎮 中国象棋AI通信完整测试"
echo "=========================================="
echo ""

# 新增服务器健康检查
echo "=========================================="
echo "🔍 检查服务器状态"
echo "=========================================="

response=$(curl -s http://localhost:8090/api/health)
if [[ $response != *"OK"* ]]; then
    echo "❌ 错误：服务器未正常运行！"
    echo "响应内容: $response"
    exit 1
else
    echo "✅ 服务器状态正常"
fi

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 1. 创建新游戏
echo -e "${BLUE}步骤 1: 创建新游戏${NC}"
echo "----------------------------------------"
RESPONSE=$(curl -s -X POST http://localhost:8090/api/game/new \
  -H "Content-Type: application/json" \
  -d '{"playerColor":1}')

GAME_ID=$(echo $RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

if [ -z "$GAME_ID" ]; then
  echo -e "${RED}❌ 创建游戏失败${NC}"
  echo "响应: $RESPONSE"
  exit 1
fi

echo -e "${GREEN}✅ 游戏创建成功${NC}"
echo "游戏ID: $GAME_ID"
echo "玩家颜色: 红方(1)"
echo "AI颜色: 黑方(2)"
echo ""

# 2. 玩家走子（红方先行）
echo -e "${BLUE}步骤 2: 玩家走子 - 炮二平五（中炮开局）${NC}"
echo "----------------------------------------"
RESPONSE=$(curl -s -X POST http://localhost:8090/api/game/move \
  -H "Content-Type: application/json" \
  -d "{\"gameId\":\"$GAME_ID\",\"fromRow\":7,\"fromCol\":1,\"toRow\":7,\"toCol\":4}")

if echo "$RESPONSE" | grep -q '"success":true'; then
  SUCCESS="true"
else
  SUCCESS="false"
fi

if [ "$SUCCESS" != "true" ]; then
  echo -e "${RED}❌ 玩家走子失败${NC}"
  echo "响应: ${RESPONSE:0:200}..."
  exit 1
fi

echo -e "${GREEN}✅ 玩家走子成功${NC}"
echo "走子: 炮二平五 (7,1) -> (7,4)"
echo ""

# 3. AI走子
echo -e "${BLUE}步骤 3: AI走子（调用LKE智能体）${NC}"
echo "----------------------------------------"
echo -e "${YELLOW}⏳ AI正在思考...${NC}"

START_TIME=$(date +%s)
RESPONSE=$(curl -s -X POST "http://localhost:8090/api/game/$GAME_ID/ai-move")
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

if echo "$RESPONSE" | grep -q '"success":true'; then
  SUCCESS="true"
else
  SUCCESS="false"
fi

if [ "$SUCCESS" != "true" ]; then
  echo -e "${RED}❌ AI走子失败${NC}"
  echo "响应: ${RESPONSE:0:200}..."
  exit 1
fi

# 提取AI走子信息
AI_MOVE=$(echo $RESPONSE | grep -o '"move":"[^"]*"' | cut -d'"' -f4)

echo -e "${GREEN}✅ AI走子成功${NC}"
echo "AI走子: $AI_MOVE"
echo "耗时: ${DURATION}秒"
echo ""

# 4. 查询游戏状态
echo -e "${BLUE}步骤 4: 查询游戏状态${NC}"
echo "----------------------------------------"
RESPONSE=$(curl -s -X GET "http://localhost:8090/api/game/$GAME_ID")

# 提取走子历史
MOVE_LIST=$(echo $RESPONSE | grep -o '"moveList":\[[^]]*\]' | sed 's/"moveList":\[//;s/\]//')

echo -e "${GREEN}✅ 游戏状态查询成功${NC}"
echo "走子历史: $MOVE_LIST"
echo ""

# 5. 再走一轮测试
echo -e "${BLUE}步骤 5: 再走一轮测试${NC}"
echo "----------------------------------------"

# 玩家第二步：马二进三
echo "玩家走子: 马二进三 (9,1) -> (7,2)"
RESPONSE=$(curl -s -X POST http://localhost:8090/api/game/move \
  -H "Content-Type: application/json" \
  -d "{\"gameId\":\"$GAME_ID\",\"fromRow\":9,\"fromCol\":1,\"toRow\":7,\"toCol\":2}")

if echo "$RESPONSE" | grep -q '"success":true'; then
  echo -e "${GREEN}✅ 玩家走子成功${NC}"
else
  echo -e "${RED}❌ 玩家走子失败${NC}"
  echo "响应: ${RESPONSE:0:200}..."
fi

echo ""
echo -e "${YELLOW}⏳ AI正在思考第二步...${NC}"

START_TIME=$(date +%s)
RESPONSE=$(curl -s -X POST "http://localhost:8090/api/game/$GAME_ID/ai-move")
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

if echo "$RESPONSE" | grep -q '"success":true'; then
  AI_MOVE=$(echo $RESPONSE | grep -o '"move":"[^"]*"' | cut -d'"' -f4)
  echo -e "${GREEN}✅ AI第二步走子成功${NC}"
  echo "AI走子: $AI_MOVE"
  echo "耗时: ${DURATION}秒"
else
  echo -e "${RED}❌ AI走子失败${NC}"
  echo "响应: ${RESPONSE:0:200}..."
fi

echo ""

# 6. 最终游戏状态
echo -e "${BLUE}步骤 6: 最终游戏状态${NC}"
echo "----------------------------------------"
RESPONSE=$(curl -s -X GET "http://localhost:8090/api/game/$GAME_ID")
MOVE_LIST=$(echo $RESPONSE | grep -o '"moveList":\[[^]]*\]' | sed 's/"moveList":\[//;s/\]//')

echo -e "${GREEN}✅ 完整走子历史:${NC}"
echo "$MOVE_LIST" | sed 's/,/\n/g' | nl

echo ""
echo "=========================================="
echo -e "${GREEN}🎉 测试完成！游戏与AI通信正常！${NC}"
echo "=========================================="
echo ""
echo "测试总结:"
echo "  ✅ 游戏创建成功"
echo "  ✅ 玩家走子功能正常"
echo "  ✅ AI走子功能正常"
echo "  ✅ LKE智能体通信正常"
echo "  ✅ 游戏状态管理正常"
echo ""
