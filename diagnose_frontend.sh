#!/bin/bash

echo "=== 前端棋子移动问题诊断脚本 ==="
echo "时间: $(date)"
echo

# 检查服务器状态
echo "1. 检查服务器状态..."
if curl -s http://localhost:8090/api/health > /dev/null; then
    echo "✅ 服务器运行正常"
else
    echo "❌ 服务器未运行，请先启动服务器"
    exit 1
fi

# 检查静态文件
echo
echo "2. 检查静态文件..."
if curl -s -I http://localhost:8090/static/app.js | grep -q "200 OK"; then
    echo "✅ app.js 可访问"
else
    echo "❌ app.js 无法访问"
fi

if curl -s -I http://localhost:8090/static/style.css | grep -q "200 OK"; then
    echo "✅ style.css 可访问"
else
    echo "❌ style.css 无法访问"
fi

if curl -s -I http://localhost:8090/ | grep -q "200 OK"; then
    echo "✅ 主页面可访问"
else
    echo "❌ 主页面无法访问"
fi

# 测试API功能
echo
echo "3. 测试API功能..."

# 测试创建游戏
echo "测试创建游戏..."
GAME_RESPONSE=$(curl -s -X POST http://localhost:8090/api/game/new \
    -H "Content-Type: application/json" \
    -d '{"playerColor": 1}')

if echo "$GAME_RESPONSE" | grep -q '"success":true'; then
    echo "✅ 游戏创建成功"
    GAME_ID=$(echo "$GAME_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo "   游戏ID: $GAME_ID"
    
    # 测试玩家移动
    echo "测试玩家移动..."
    MOVE_RESPONSE=$(curl -s -X POST http://localhost:8090/api/game/move \
        -H "Content-Type: application/json" \
        -d "{\"gameId\": \"$GAME_ID\", \"fromRow\": 9, \"fromCol\": 4, \"toRow\": 8, \"toCol\": 4}")
    
    if echo "$MOVE_RESPONSE" | grep -q '"success":true'; then
        echo "✅ 玩家移动成功"
    else
        echo "❌ 玩家移动失败"
        echo "   响应: $MOVE_RESPONSE"
    fi
else
    echo "❌ 游戏创建失败"
    echo "   响应: $GAME_RESPONSE"
fi

# 检查浏览器兼容性
echo
echo "4. 浏览器兼容性检查..."
echo "请在浏览器中打开以下页面进行测试："
echo
echo "🔧 调试页面: http://localhost:8090/frontend_debug.html"
echo "🎮 游戏页面: http://localhost:8090/"
echo
echo "调试步骤："
echo "1. 打开调试页面，检查系统状态是否全部为绿色"
echo "2. 点击'开始新游戏'按钮"
echo "3. 观察调试日志中的信息"
echo "4. 尝试点击棋子，查看是否有选中效果"
echo "5. 尝试移动棋子，观察日志中的移动请求"
echo
echo "如果问题仍然存在，请："
echo "- 打开浏览器开发者工具 (F12)"
echo "- 查看Console标签页的错误信息"
echo "- 查看Network标签页的网络请求"
echo "- 将错误信息反馈给开发者"

echo
echo "=== 诊断完成 ==="