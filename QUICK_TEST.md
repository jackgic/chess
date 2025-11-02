# 🎮 快速测试指南

## 立即测试游戏

### 方法1: 使用浏览器（推荐）

1. **打开浏览器**
   ```
   http://localhost:8090
   ```

2. **检查页面显示**
   - ✅ 看到精美的中国风界面
   - ✅ 看到"中国象棋 AI 对弈"标题
   - ✅ 看到红方/黑方选择按钮
   - ✅ 看到"开始新游戏"按钮

3. **开始游戏**
   - 点击"红方（先手）"或"黑方（后手）"
   - 点击"开始新游戏"按钮
   - 等待棋盘加载

4. **验证功能**
   - ✅ 棋盘正常显示
   - ✅ 棋子正常显示
   - ✅ 可以点击棋子
   - ✅ 可以移动棋子
   - ✅ AI能正常走子

### 方法2: 使用开发者工具

1. **打开浏览器开发者工具**
   - Windows: `F12` 或 `Ctrl+Shift+I`
   - Mac: `Cmd+Option+I`

2. **检查控制台 (Console)**
   - ✅ 应该看到: "Canvas初始化成功"
   - ✅ 应该看到: "前端初始化完成"
   - ❌ 不应该有红色错误信息

3. **检查网络 (Network)**
   - ✅ 所有资源状态码应该是 200
   - ✅ `index.html` - 200 OK
   - ✅ `app.js` - 200 OK
   - ✅ `style.css` - 200 OK

### 方法3: 命令行测试

```bash
# 1. 测试服务器
curl http://localhost:8090/api/health

# 2. 测试HTML
curl -I http://localhost:8090/

# 3. 测试JavaScript
curl -I http://localhost:8090/static/app.js

# 4. 测试CSS
curl -I http://localhost:8090/static/style.css
```

## 常见问题

### Q1: 页面空白
**解决方法**:
1. 清除浏览器缓存 (`Ctrl+Shift+Delete`)
2. 强制刷新页面 (`Ctrl+F5` 或 `Cmd+Shift+R`)
3. 检查浏览器控制台是否有错误

### Q2: 棋盘不显示
**解决方法**:
1. 确认已点击"开始新游戏"
2. 检查浏览器控制台是否有错误
3. 尝试使用其他浏览器

### Q3: 无法走子
**解决方法**:
1. 确认是你的回合（不是AI回合）
2. 确认选择的是己方棋子
3. 确认移动符合象棋规则

### Q4: AI不走子
**解决方法**:
1. 等待AI思考（通常5-10秒）
2. 检查浏览器控制台是否有错误
3. 查看服务器日志: `tail -f logs/chinese-chess-ai.log`

## 推荐浏览器

✅ **最佳体验**:
- Chrome 90+
- Edge 90+
- Firefox 88+
- Safari 14+

❌ **不支持**:
- Internet Explorer

## 获取帮助

如果遇到问题，请：

1. **查看文档**:
   - [FRONTEND_FIX.md](FRONTEND_FIX.md) - 前端问题修复
   - [COORDINATE_SYSTEM_FIX.md](COORDINATE_SYSTEM_FIX.md) - 坐标系统修复
   - [QUICKSTART.md](QUICKSTART.md) - 快速开始

2. **查看日志**:
   ```bash
   tail -f logs/chinese-chess-ai.log
   ```

3. **检查服务器**:
   ```bash
   ./status.sh
   ```

## 🎉 开始游戏吧！

访问: **http://localhost:8090**

祝你玩得开心！🏮♟️
