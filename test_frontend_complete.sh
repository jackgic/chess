#!/bin/bash

echo "=== 前端功能完整测试 ==="
echo "测试时间: $(date)"
echo

# 检查服务器状态
echo "1. 检查服务器状态..."
if curl -s http://localhost:8090/api/health > /dev/null; then
    echo "✅ 服务器运行正常"
else
    echo "❌ 服务器未运行"
    exit 1
fi

# 测试API功能
echo
echo "2. 测试API功能..."

# 创建游戏
echo "   创建新游戏..."
GAME_RESPONSE=$(curl -s -X POST http://localhost:8090/api/game/new \
    -H "Content-Type: application/json" \
    -d '{"playerColor": 1}')

if echo "$GAME_RESPONSE" | grep -q '"success":true'; then
    echo "   ✅ 游戏创建成功"
    GAME_ID=$(echo "$GAME_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo "   游戏ID: $GAME_ID"
else
    echo "   ❌ 游戏创建失败"
    echo "   响应: $GAME_RESPONSE"
    exit 1
fi

# 测试走子
echo "   测试玩家走子..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8090/api/game/move \
    -H "Content-Type: application/json" \
    -d "{\"gameId\": \"$GAME_ID\", \"fromRow\": 7, \"fromCol\": 1, \"toRow\": 7, \"toCol\": 4}")

if echo "$MOVE_RESPONSE" | grep -q '"success":true'; then
    echo "   ✅ 玩家走子成功"
else
    echo "   ❌ 玩家走子失败"
    echo "   响应: $MOVE_RESPONSE"
fi

# 测试前端页面访问
echo
echo "3. 测试前端页面访问..."

# 主页面
if curl -s http://localhost:8090/ | grep -q "中国象棋 AI 对弈"; then
    echo "   ✅ 主页面可访问"
else
    echo "   ❌ 主页面访问失败"
fi

# 测试页面
if curl -s http://localhost:8090/test_frontend.html | grep -q "前端功能测试"; then
    echo "   ✅ API测试页面可访问"
else
    echo "   ❌ API测试页面访问失败"
fi

if curl -s http://localhost:8090/test_canvas.html | grep -q "Canvas棋盘测试"; then
    echo "   ✅ Canvas测试页面可访问"
else
    echo "   ❌ Canvas测试页面访问失败"
fi

# 静态资源
if curl -s http://localhost:8090/static/app.js | grep -q "gameState"; then
    echo "   ✅ JavaScript文件可访问"
else
    echo "   ❌ JavaScript文件访问失败"
fi

if curl -s http://localhost:8090/static/style.css | grep -q "container"; then
    echo "   ✅ CSS文件可访问"
else
    echo "   ❌ CSS文件访问失败"
fi

echo
echo "=== 测试完成 ==="
echo
echo "📋 使用说明："
echo "1. 打开主页面: http://localhost:8090/"
echo "2. 选择颜色并点击'开始新游戏'"
echo "3. 点击棋子选择，再点击目标位置移动"
echo "4. 如有问题，查看浏览器控制台日志"
echo
echo "🔧 调试页面："
echo "- API测试: http://localhost:8090/test_frontend.html"
echo "- Canvas测试: http://localhost:8090/test_canvas.html"
echo
echo "📊 服务器状态: ./status.sh"
echo "📋 查看日志: ./logs.sh"