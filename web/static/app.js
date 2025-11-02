// 游戏状态
let gameState = null;
let selectedPiece = null;
let canvas = null;
let ctx = null;

// 常量
const BOARD_PADDING = 30;
const CELL_SIZE = 60;
const PIECE_RADIUS = 25;

// 棋子名称映射（备用，后端已提供name字段）
const PIECE_NAMES = {
    1: { 1: '帅', 2: '将' },
    2: { 1: '仕', 2: '士' },
    3: { 1: '相', 2: '象' },
    4: { 1: '马', 2: '马' },
    5: { 1: '车', 2: '车' },
    6: { 1: '炮', 2: '炮' },
    7: { 1: '兵', 2: '卒' }
};

// 初始化
document.addEventListener('DOMContentLoaded', () => {
    canvas = document.getElementById('chessBoard');
    if (!canvas) {
        console.error('找不到棋盘Canvas元素');
        return;
    }
    
    ctx = canvas.getContext('2d');
    if (!ctx) {
        console.error('无法获取Canvas 2D上下文');
        return;
    }
    
    console.log('Canvas初始化成功，尺寸:', canvas.width, 'x', canvas.height);

    // 颜色选择按钮
    document.querySelectorAll('.color-option').forEach(btn => {
        btn.addEventListener('click', () => {
            document.querySelectorAll('.color-option').forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
        });
    });

    // 开始游戏按钮
    const startGameBtn = document.getElementById('startGameBtn');
    console.log('startGameBtn元素:', startGameBtn);
    if (startGameBtn) {
        console.log('正在为startGameBtn添加点击事件监听器');
        startGameBtn.addEventListener('click', () => {
            console.log('=== 开始游戏按钮被点击 ===');
            startNewGame();
        });
        console.log('startGameBtn点击事件监听器已添加');
    } else {
        console.error('找不到startGameBtn元素！');
    }

    // 移除了重复的新游戏按钮，使用棋盘控制区域的重新开始按钮

    // 操作按钮
    const undoBtn = document.getElementById('undoBtn');
    if (undoBtn) {
        undoBtn.addEventListener('click', undoMove);
    }

    const restartBtn = document.getElementById('restartBtn');
    if (restartBtn) {
        restartBtn.addEventListener('click', () => {
            if (confirm('确定要重新开始对局吗？')) {
                startNewGame();
            }
        });
    }

    const hintBtn = document.getElementById('hintBtn');
    if (hintBtn) {
        hintBtn.addEventListener('click', getHint);
    }

    // 棋盘点击事件
    canvas.addEventListener('click', handleBoardClick);
    
    console.log('前端初始化完成');
});

// 开始新游戏
async function startNewGame() {
    console.log('=== startNewGame 函数被调用 ===');
    
    const activeBtn = document.querySelector('.color-option.active');
    console.log('选中的按钮:', activeBtn);
    
    if (!activeBtn) {
        alert('请先选择阵营（红方或黑方）');
        return;
    }
    
    const playerColor = parseInt(activeBtn.dataset.color);
    console.log('开始新游戏，玩家颜色:', playerColor);

    try {
        console.log('正在发送请求到 /api/game/new...');
        const response = await fetch('/api/game/new', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ playerColor })
        });

        console.log('收到响应，状态码:', response.status);
        const result = await response.json();
        console.log('游戏创建响应:', result);
        
        if (result.success) {
            gameState = result.data;
            selectedPiece = null;
            console.log('游戏状态设置完成:', gameState);

            // 隐藏阵营选择，显示游戏界面
            document.getElementById('gameControls').style.display = 'none';
            document.getElementById('gameInfo').style.display = 'block';
            document.getElementById('gameArea').style.display = 'flex';
            
            // 更新玩家颜色显示
            const playerColorText = document.getElementById('playerColorText');
            if (playerColorText) {
                playerColorText.textContent = gameState.playerColor === 1 ? '红方（先手）' : '黑方（后手）';
            }

            // 绘制棋盘
            drawBoard();
            updateGameInfo();

            // 红方先走，检查是否需要AI先手
            // gameState.turn 始终是1（红方先手）
            if (gameState.turn === gameState.aiColor) {
                console.log('AI是红方，AI先手');
                setTimeout(() => requestAIMove(), 500);
            } else {
                console.log('玩家是红方，玩家先手');
            }
        } else {
            alert('创建游戏失败：' + result.error);
        }
    } catch (error) {
        console.error('创建游戏失败:', error);
        alert('创建游戏失败，请检查服务器连接');
    }
}

// 将棋盘坐标转换为显示坐标（根据玩家方翻转）
function boardToDisplay(row, col) {
    if (gameState && gameState.playerColor === 2) {
        // 玩家是黑方，翻转棋盘
        return { row: 9 - row, col: 8 - col };
    }
    return { row, col };
}

// 将显示坐标转换为棋盘坐标
function displayToBoard(row, col) {
    if (gameState && gameState.playerColor === 2) {
        // 玩家是黑方，翻转棋盘
        return { row: 9 - row, col: 8 - col };
    }
    return { row, col };
}

// 绘制棋盘
function drawBoard() {
    // 清空画布
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // 绘制木质纹理背景
    const gradient = ctx.createLinearGradient(0, 0, canvas.width, canvas.height);
    gradient.addColorStop(0, '#F5DEB3');
    gradient.addColorStop(0.5, '#E8C5A0');
    gradient.addColorStop(1, '#F5DEB3');
    ctx.fillStyle = gradient;
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    
    // 添加木质纹理效果
    ctx.globalAlpha = 0.08;
    for (let i = 0; i < 50; i++) {
        ctx.strokeStyle = Math.random() > 0.5 ? '#C9A882' : '#B8956A';
        ctx.lineWidth = Math.random() * 2;
        ctx.beginPath();
        ctx.moveTo(Math.random() * canvas.width, 0);
        ctx.lineTo(Math.random() * canvas.width, canvas.height);
        ctx.stroke();
    }
    ctx.globalAlpha = 1.0;

    // 绘制棋盘线（带阴影效果）
    ctx.shadowColor = 'rgba(0, 0, 0, 0.2)';
    ctx.shadowBlur = 1.5;
    ctx.shadowOffsetX = 0.5;
    ctx.shadowOffsetY = 0.5;
    ctx.strokeStyle = '#5A5A5A';
    ctx.lineWidth = 2;

    // 横线（楚河汉界处不断开，保持完整）
    for (let i = 0; i < 10; i++) {
        ctx.beginPath();
        ctx.moveTo(BOARD_PADDING, BOARD_PADDING + i * CELL_SIZE);
        ctx.lineTo(BOARD_PADDING + 8 * CELL_SIZE, BOARD_PADDING + i * CELL_SIZE);
        ctx.stroke();
    }

    // 竖线
    for (let i = 0; i < 9; i++) {
        ctx.beginPath();
        ctx.moveTo(BOARD_PADDING + i * CELL_SIZE, BOARD_PADDING);
        if (i === 0 || i === 8) {
            // 边线：完整的竖线
            ctx.lineTo(BOARD_PADDING + i * CELL_SIZE, BOARD_PADDING + 9 * CELL_SIZE);
        } else {
            // 中间的竖线：在楚河汉界处断开，间隔更紧凑
            // 上半部分：从顶部到第4行稍微延伸一点
            ctx.lineTo(BOARD_PADDING + i * CELL_SIZE, BOARD_PADDING + 4.2 * CELL_SIZE);
            ctx.stroke();
            
            // 下半部分：从第5行稍微提前开始到底部
            ctx.beginPath();
            ctx.moveTo(BOARD_PADDING + i * CELL_SIZE, BOARD_PADDING + 4.8 * CELL_SIZE);
            ctx.lineTo(BOARD_PADDING + i * CELL_SIZE, BOARD_PADDING + 9 * CELL_SIZE);
        }
        ctx.stroke();
    }

    // 绘制九宫格斜线（根据玩家颜色调整位置）
    if (gameState && gameState.playerColor === 2) {
        // 玩家是黑方，翻转显示
        drawDiagonal(3, 9, 5, 7);  // 下方九宫格（黑方）
        drawDiagonal(5, 9, 3, 7);
        drawDiagonal(3, 2, 5, 0);  // 上方九宫格（红方）
        drawDiagonal(5, 2, 3, 0);
    } else {
        // 玩家是红方，正常显示
        drawDiagonal(3, 0, 5, 2);  // 上方九宫格（黑方）
        drawDiagonal(5, 0, 3, 2);
        drawDiagonal(3, 7, 5, 9);  // 下方九宫格（红方）
        drawDiagonal(5, 7, 3, 9);
    }

    // 重置阴影
    ctx.shadowColor = 'transparent';
    ctx.shadowBlur = 0;
    ctx.shadowOffsetX = 0;
    ctx.shadowOffsetY = 0;
    
    // 绘制楚河汉界（居中显示在棋盘正中间）
    ctx.font = 'bold 24px "STKaiti", "KaiTi", "SimKai", serif';
    ctx.textAlign = 'center';
    
    // 添加文字阴影效果（减少阴影避免覆盖线条）
    ctx.shadowColor = 'rgba(0, 0, 0, 0.3)';
    ctx.shadowBlur = 2;
    ctx.shadowOffsetX = 1;
    ctx.shadowOffsetY = 1;
    
    // 计算楚河汉界的正确位置（第4行和第5行之间，精确居中）
    const centerY = BOARD_PADDING + 4.5 * CELL_SIZE; // 垂直居中位置，不偏移
    const leftTextX = BOARD_PADDING + 1.8 * CELL_SIZE;   // 左侧文字位置，向左移动
    const rightTextX = BOARD_PADDING + 6.2 * CELL_SIZE;  // 右侧文字位置，向右移动
    
    // 设置文字样式（减少描边宽度避免覆盖线条）
    ctx.strokeStyle = '#654321';  // 深棕色描边
    ctx.lineWidth = 1.5;
    ctx.fillStyle = '#DAA520';    // 金色填充
    
    // 按照传统象棋标准显示楚河汉界
    // 楚河汉界是水平排列的，位置是固定的：
    // 从红方（下方）视角：左侧是"楚河"，右侧是"汉界"
    // 这个排列不会因为玩家颜色而改变，是象棋的标准布局
    
    if (gameState && gameState.playerColor === 2) {
        // 玩家是黑方，棋盘翻转显示
        // 棋盘翻转180度后，左右位置互换
        
        // 左侧显示"汉界"（翻转后原右侧的汉界到了左侧）
        ctx.strokeText('汉界', leftTextX, centerY);
        ctx.fillText('汉界', leftTextX, centerY);
        
        // 右侧显示"楚河"（翻转后原左侧的楚河到了右侧）
        ctx.strokeText('楚河', rightTextX, centerY);
        ctx.fillText('楚河', rightTextX, centerY);
    } else {
        // 玩家是红方，正常显示
        // 按照传统象棋标准布局
        
        // 左侧显示"楚河"
        ctx.strokeText('楚河', leftTextX, centerY);
        ctx.fillText('楚河', leftTextX, centerY);
        
        // 右侧显示"汉界"
        ctx.strokeText('汉界', rightTextX, centerY);
        ctx.fillText('汉界', rightTextX, centerY);
    }
    
    // 在楚河汉界中间绘制装饰性分割线（避免与棋盘线重叠）
    ctx.strokeStyle = '#8B4513';
    ctx.lineWidth = 1;
    ctx.setLineDash([6, 3]);
    ctx.beginPath();
    // 中间装饰线，位于楚河汉界文字之间，更加精致
    ctx.moveTo(BOARD_PADDING + 3.6 * CELL_SIZE, centerY - 6);
    ctx.lineTo(BOARD_PADDING + 4.4 * CELL_SIZE, centerY - 6);
    ctx.moveTo(BOARD_PADDING + 3.6 * CELL_SIZE, centerY + 6);
    ctx.lineTo(BOARD_PADDING + 4.4 * CELL_SIZE, centerY + 6);
    ctx.stroke();
    ctx.setLineDash([]); // 重置虚线
    
    // 重置阴影
    ctx.shadowColor = 'transparent';
    ctx.shadowBlur = 0;
    ctx.shadowOffsetX = 0;
    ctx.shadowOffsetY = 0;

    // 绘制棋子
    if (gameState && gameState.board) {
        let pieceCount = 0;
        for (let row = 0; row < 10; row++) {
            for (let col = 0; col < 9; col++) {
                const piece = gameState.board[row][col];
                if (piece.type !== 0) {
                    drawPiece(row, col, piece);
                    pieceCount++;
                }
            }
        }
        console.log('绘制了', pieceCount, '个棋子');
    } else {
        console.log('没有游戏状态或棋盘数据');
    }
}

// 绘制斜线
function drawDiagonal(col1, row1, col2, row2) {
    ctx.beginPath();
    ctx.moveTo(BOARD_PADDING + col1 * CELL_SIZE, BOARD_PADDING + row1 * CELL_SIZE);
    ctx.lineTo(BOARD_PADDING + col2 * CELL_SIZE, BOARD_PADDING + row2 * CELL_SIZE);
    ctx.stroke();
}

// 绘制棋子
function drawPiece(row, col, piece) {
    const isSelected = selectedPiece && selectedPiece.row === row && selectedPiece.col === col;
    const radius = isSelected ? PIECE_RADIUS * 1.2 : PIECE_RADIUS;

    // 转换为显示坐标
    const displayPos = boardToDisplay(row, col);
    const x = BOARD_PADDING + displayPos.col * CELL_SIZE;
    const y = BOARD_PADDING + displayPos.row * CELL_SIZE;

    // 检查是否是将军状态下的将帅
    const isKingInCheck = piece.type === 1 && 
        ((piece.color === 1 && gameState.inCheck && gameState.inCheck.player) ||
         (piece.color === 2 && gameState.inCheck && gameState.inCheck.ai));

    ctx.save();
    
    // ========== 超强3D立体效果 ==========
    
    // 1. 底部投影（模拟棋子悬浮）
    const shadowGradient = ctx.createRadialGradient(x, y + 5, 0, x, y + 5, radius + 10);
    shadowGradient.addColorStop(0, 'rgba(0, 0, 0, 0.6)');
    shadowGradient.addColorStop(0.4, 'rgba(0, 0, 0, 0.3)');
    shadowGradient.addColorStop(1, 'rgba(0, 0, 0, 0)');
    ctx.fillStyle = shadowGradient;
    ctx.beginPath();
    ctx.arc(x, y + 5, radius + 8, 0, 2 * Math.PI);
    ctx.fill();
    
    // 2. 棋子底座（厚度感）
    for (let i = 4; i >= 0; i--) {
        const offset = i * 0.8;
        const baseRadius = radius + 3 - i * 0.5;
        const alpha = 0.15 - i * 0.02;
        
        ctx.fillStyle = piece.color === 1 ? 
            `rgba(139, 115, 85, ${alpha})` : 
            `rgba(42, 42, 42, ${alpha})`;
        ctx.beginPath();
        ctx.arc(x, y + offset, baseRadius, 0, 2 * Math.PI);
        ctx.fill();
    }
    
    // 3. 棋子主体（球面渐变 - 模拟真实光照）
    const mainGradient = ctx.createRadialGradient(
        x - radius * 0.5,  // 光源位置：左上
        y - radius * 0.5, 
        0, 
        x, 
        y, 
        radius * 1.4
    );
    
    if (piece.color === 1) {
        // 红方 - 温暖的木质感
        mainGradient.addColorStop(0, '#FFFEF5');      // 最亮点
        mainGradient.addColorStop(0.15, '#FFF8E1');   // 高光区
        mainGradient.addColorStop(0.35, '#F5E6D0');   // 亮部
        mainGradient.addColorStop(0.55, '#E8D4B8');   // 中间调
        mainGradient.addColorStop(0.75, '#D4C0A0');   // 暗部
        mainGradient.addColorStop(0.9, '#B8A080');    // 深暗部
        mainGradient.addColorStop(1, '#9A8060');      // 边缘最暗
    } else {
        // 黑方 - 石质感
        mainGradient.addColorStop(0, '#C0C0C0');      // 最亮点
        mainGradient.addColorStop(0.15, '#A0A0A0');   // 高光区
        mainGradient.addColorStop(0.35, '#808080');   // 亮部
        mainGradient.addColorStop(0.55, '#606060');   // 中间调
        mainGradient.addColorStop(0.75, '#404040');   // 暗部
        mainGradient.addColorStop(0.9, '#2A2A2A');    // 深暗部
        mainGradient.addColorStop(1, '#1A1A1A');      // 边缘最暗
    }
    
    ctx.fillStyle = mainGradient;
    ctx.beginPath();
    ctx.arc(x, y, radius, 0, 2 * Math.PI);
    ctx.fill();
    
    // 4. 强烈的高光反射（模拟光泽）
    const highlightGradient = ctx.createRadialGradient(
        x - radius * 0.4, 
        y - radius * 0.4, 
        0, 
        x - radius * 0.4, 
        y - radius * 0.4, 
        radius * 0.8
    );
    highlightGradient.addColorStop(0, 'rgba(255, 255, 255, 0.95)');
    highlightGradient.addColorStop(0.3, 'rgba(255, 255, 255, 0.6)');
    highlightGradient.addColorStop(0.6, 'rgba(255, 255, 255, 0.2)');
    highlightGradient.addColorStop(1, 'rgba(255, 255, 255, 0)');
    ctx.fillStyle = highlightGradient;
    ctx.beginPath();
    ctx.arc(x - radius * 0.3, y - radius * 0.3, radius * 0.5, 0, 2 * Math.PI);
    ctx.fill();
    
    // 5. 次级高光（增强球面感）
    const secondaryHighlight = ctx.createRadialGradient(
        x - radius * 0.2, 
        y - radius * 0.3, 
        0, 
        x - radius * 0.2, 
        y - radius * 0.3, 
        radius * 0.4
    );
    secondaryHighlight.addColorStop(0, 'rgba(255, 255, 255, 0.4)');
    secondaryHighlight.addColorStop(0.5, 'rgba(255, 255, 255, 0.15)');
    secondaryHighlight.addColorStop(1, 'rgba(255, 255, 255, 0)');
    ctx.fillStyle = secondaryHighlight;
    ctx.beginPath();
    ctx.arc(x - radius * 0.15, y - radius * 0.2, radius * 0.35, 0, 2 * Math.PI);
    ctx.fill();
    
    // 6. 底部阴影（增强立体感）
    const bottomShadow = ctx.createRadialGradient(
        x + radius * 0.4, 
        y + radius * 0.4, 
        0, 
        x + radius * 0.4, 
        y + radius * 0.4, 
        radius * 0.9
    );
    bottomShadow.addColorStop(0, 'rgba(0, 0, 0, 0.5)');
    bottomShadow.addColorStop(0.5, 'rgba(0, 0, 0, 0.25)');
    bottomShadow.addColorStop(1, 'rgba(0, 0, 0, 0)');
    ctx.fillStyle = bottomShadow;
    ctx.beginPath();
    ctx.arc(x + radius * 0.25, y + radius * 0.25, radius * 0.8, 0, 2 * Math.PI);
    ctx.fill();
    
    // 7. 环境光遮蔽（AO效果）
    ctx.strokeStyle = 'rgba(0, 0, 0, 0.3)';
    ctx.lineWidth = 2;
    ctx.beginPath();
    ctx.arc(x, y, radius - 1, 0, 2 * Math.PI);
    ctx.stroke();
    
    // 8. 边缘高光（Rim Light）
    const rimGradient = ctx.createLinearGradient(
        x - radius, 
        y - radius, 
        x + radius, 
        y + radius
    );
    rimGradient.addColorStop(0, 'rgba(255, 255, 255, 0.4)');
    rimGradient.addColorStop(0.5, 'rgba(255, 255, 255, 0)');
    rimGradient.addColorStop(1, 'rgba(0, 0, 0, 0.3)');
    ctx.strokeStyle = rimGradient;
    ctx.lineWidth = 3;
    ctx.beginPath();
    ctx.arc(x, y, radius, 0, 2 * Math.PI);
    ctx.stroke();
    
    // 9. 外边框（根据状态）
    if (isKingInCheck) {
        ctx.strokeStyle = '#FF0000';
        ctx.lineWidth = 4;
        ctx.beginPath();
        ctx.arc(x, y, radius + 1, 0, 2 * Math.PI);
        ctx.stroke();
        
        ctx.strokeStyle = 'rgba(255, 0, 0, 0.5)';
        ctx.lineWidth = 6;
        ctx.beginPath();
        ctx.arc(x, y, radius + 4, 0, 2 * Math.PI);
        ctx.stroke();
    } else if (isSelected) {
        ctx.strokeStyle = '#FFD700';
        ctx.lineWidth = 4;
        ctx.beginPath();
        ctx.arc(x, y, radius + 1, 0, 2 * Math.PI);
        ctx.stroke();
        
        ctx.strokeStyle = 'rgba(255, 215, 0, 0.5)';
        ctx.lineWidth = 6;
        ctx.beginPath();
        ctx.arc(x, y, radius + 4, 0, 2 * Math.PI);
        ctx.stroke();
    } else {
        const borderGradient = ctx.createLinearGradient(
            x - radius, y - radius, 
            x + radius, y + radius
        );
        if (piece.color === 1) {
            borderGradient.addColorStop(0, '#A08060');
            borderGradient.addColorStop(0.5, '#806040');
            borderGradient.addColorStop(1, '#604020');
        } else {
            borderGradient.addColorStop(0, '#505050');
            borderGradient.addColorStop(0.5, '#303030');
            borderGradient.addColorStop(1, '#101010');
        }
        ctx.strokeStyle = borderGradient;
        ctx.lineWidth = 3.5;
        ctx.beginPath();
        ctx.arc(x, y, radius, 0, 2 * Math.PI);
        ctx.stroke();
    }
    
    // 10. 棋子文字（超强立体效果）
    ctx.font = 'bold 28px "STKaiti", "KaiTi", "SimKai", serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    
    const pieceName = piece.name || '?';
    
    // 文字深层阴影
    ctx.fillStyle = 'rgba(0, 0, 0, 0.6)';
    ctx.fillText(pieceName, x + 2, y + 2);
    
    // 文字中层阴影
    ctx.fillStyle = 'rgba(0, 0, 0, 0.4)';
    ctx.fillText(pieceName, x + 1, y + 1);
    
    // 文字外描边（最深）
    ctx.strokeStyle = piece.color === 1 ? '#6B2A3F' : '#000000';
    ctx.lineWidth = 5;
    ctx.strokeText(pieceName, x, y);
    
    // 文字中描边
    ctx.strokeStyle = piece.color === 1 ? '#8B3A4F' : '#1A1A1A';
    ctx.lineWidth = 3.5;
    ctx.strokeText(pieceName, x, y);
    
    // 文字内描边
    ctx.strokeStyle = piece.color === 1 ? '#AB4A5F' : '#2A2A2A';
    ctx.lineWidth = 2;
    ctx.strokeText(pieceName, x, y);
    
    // 文字主体填充
    const textGradient = ctx.createLinearGradient(x, y - 15, x, y + 15);
    if (piece.color === 1) {
        textGradient.addColorStop(0, '#FF6B7A');
        textGradient.addColorStop(0.5, '#E85D75');
        textGradient.addColorStop(1, '#D04D65');
    } else {
        textGradient.addColorStop(0, '#FFFFFF');
        textGradient.addColorStop(0.5, '#F0F0F0');
        textGradient.addColorStop(1, '#E0E0E0');
    }
    ctx.fillStyle = textGradient;
    ctx.fillText(pieceName, x, y);
    
    // 文字顶部高光
    ctx.fillStyle = piece.color === 1 ? 
        'rgba(255, 200, 200, 0.8)' : 
        'rgba(255, 255, 255, 0.7)';
    ctx.fillText(pieceName, x - 1, y - 1);
    
    ctx.restore();
    
    if (!piece.name) {
        console.warn(`棋子在位置 (${row},${col}) 没有名称:`, piece);
    }
}

// 处理棋盘点击
function handleBoardClick(event) {
    console.log('棋盘被点击');
    
    if (!gameState || gameState.status !== 'playing') {
        console.log('游戏状态无效:', gameState?.status);
        return;
    }

    // 检查是否是玩家回合
    if (gameState.turn !== gameState.playerColor) {
        console.log('不是玩家回合:', gameState.turn, 'vs', gameState.playerColor);
        alert('现在是AI的回合');
        return;
    }

    const rect = canvas.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;

    // 计算点击的格子（显示坐标）
    const displayCol = Math.round((x - BOARD_PADDING) / CELL_SIZE);
    const displayRow = Math.round((y - BOARD_PADDING) / CELL_SIZE);

    if (displayRow < 0 || displayRow > 9 || displayCol < 0 || displayCol > 8) {
        console.log('点击位置超出棋盘范围');
        return;
    }

    // 转换为棋盘坐标
    const boardPos = displayToBoard(displayRow, displayCol);
    const row = boardPos.row;
    const col = boardPos.col;

    console.log('点击坐标:', { x, y, displayRow, displayCol, row, col });

    const piece = gameState.board[row][col];
    console.log('点击的棋子:', piece);

    // 如果点击的是本方棋子
    if (piece.type !== 0 && piece.color === gameState.playerColor) {
        // 如果已经选中了该棋子，则取消选中
        if (selectedPiece && selectedPiece.row === row && selectedPiece.col === col) {
            console.log('取消选择棋子');
            selectedPiece = null;
        } else {
            // 否则选中该棋子
            console.log('选择棋子:', piece);
            selectedPiece = { row, col };
        }
        drawBoard();
        return;
    }

    // 如果已有选中的棋子，尝试移动
    if (selectedPiece) {
        console.log('尝试移动棋子从', selectedPiece, '到', { row, col });
        playerMove(selectedPiece.row, selectedPiece.col, row, col);
    }
}

// 玩家走子
async function playerMove(fromRow, fromCol, toRow, toCol) {
    try {
        console.log('发送移动请求:', { gameId: gameState.id, fromRow, fromCol, toRow, toCol });
        
        const response = await fetch('/api/game/move', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                gameId: gameState.id,
                fromRow,
                fromCol,
                toRow,
                toCol
            })
        });

        const result = await response.json();
        console.log('移动响应:', result);
        
        if (result.success) {
            gameState = result.data;
            selectedPiece = null;
            drawBoard();
            updateGameInfo();

            // 检查游戏是否结束
            if (gameState.status !== 'playing') {
                showGameOver();
                return;
            }

            // AI走子
            setTimeout(() => requestAIMove(), 500);
        } else {
            // 处理错误响应
            const errorMsg = result.error || '移动失败';
            console.error('移动失败:', errorMsg);
            alert('移动失败：' + errorMsg);
            selectedPiece = null;
            drawBoard();
        }
    } catch (error) {
        console.error('移动请求异常:', error);
        alert('移动失败，请重试');
        selectedPiece = null;
        drawBoard();
    }
}

// 请求AI走子
async function requestAIMove() {
    // 显示思考中
    document.getElementById('aiThinking').style.display = 'flex';
    document.getElementById('aiResponse').style.display = 'none';

    try {
        const response = await fetch(`/api/game/${gameState.id}/ai-move`, {
            method: 'POST'
        });

        const result = await response.json();
        document.getElementById('aiThinking').style.display = 'none';

        if (result.success) {
            gameState = result.data.state;
            drawBoard();
            updateGameInfo();

            // 显示AI分析
            if (result.data.answer) {
                const aiAnswer = document.getElementById('aiAnswer');
                let displayText = result.data.answer;
                
                // 检查是否包含思考过程
                if (displayText.includes('思考过程：') && displayText.includes('最终决策：')) {
                    // 分割思考过程和最终决策
                    const parts = displayText.split('最终决策：');
                    const thoughtProcess = parts[0].replace('思考过程：', '').trim();
                    const finalDecision = parts[1].trim();
                    
                    // 创建格式化显示
                    aiAnswer.innerHTML = `
                        <div class="thought-process">
                            <h4>🤔 AI思考过程</h4>
                            <div class="thought-content">${thoughtProcess.replace(/\n/g, '<br>')}</div>
                        </div>
                        <div class="final-decision">
                            <h4>🎯 最终决策</h4>
                            <div class="decision-content">${finalDecision.replace(/\n/g, '<br>')}</div>
                        </div>
                    `;
                } else {
                    // 直接显示完整的思考内容
                    aiAnswer.innerHTML = displayText.replace(/\n/g, '<br>');
                }
                
                document.getElementById('aiResponse').style.display = 'block';
            }

            // 检查游戏是否结束
            if (gameState.status !== 'playing') {
                showGameOver();
            }
        } else {
            // AI走子失败，显示详细错误信息
            console.error('AI走子失败:', result.error);
            const errorMsg = result.error || 'AI走子失败';
            alert('AI走子失败：' + errorMsg + '\n\n游戏将继续，请等待AI重新思考...');
            
            // 重新尝试AI走子（最多重试一次）
            if (!result.retried) {
                setTimeout(() => {
                    console.log('重新尝试AI走子...');
                    requestAIMove();
                }, 1000);
            }
        }
    } catch (error) {
        console.error('AI走子失败:', error);
        document.getElementById('aiThinking').style.display = 'none';
        alert('AI走子失败，请重试');
    }
}

// 更新游戏信息
function updateGameInfo() {
    const turnText = gameState.turn === 1 ? '红方' : '黑方';
    document.getElementById('currentTurn').textContent = turnText;

    const statusMap = {
        'playing': '进行中',
        'red_win': '红方胜',
        'black_win': '黑方胜',
        'draw': '和棋'
    };
    document.getElementById('gameStatus').textContent = statusMap[gameState.status] || '未知';

    // 将军状态展示已移除

    // 更新走子历史
    const moveList = document.getElementById('moveList');
    moveList.innerHTML = '';
    gameState.moveList.forEach((move, index) => {
        const moveItem = document.createElement('div');
        moveItem.className = 'move-item ' + (index % 2 === 0 ? 'red' : 'black');
        moveItem.textContent = `${index + 1}. ${move}`;
        moveList.appendChild(moveItem);
    });
}

// 显示游戏结束
function showGameOver() {
    const statusMap = {
        'red_win': '红方获胜！',
        'black_win': '黑方获胜！',
        'draw': '和棋！'
    };

    const message = statusMap[gameState.status] || '游戏结束';
    
    setTimeout(() => {
        if (confirm(message + '\n是否开始新游戏？')) {
            // 显示阵营选择
            document.getElementById('gameControls').style.display = 'flex';
            startNewGame();
        }
    }, 500);
}

// 悔棋
async function undoMove() {
    if (!gameState || gameState.status !== 'playing') {
        alert('当前无法悔棋');
        return;
    }

    if (gameState.moveList.length < 2) {
        alert('还没有足够的步数可以悔棋');
        return;
    }

    if (!confirm('确定要悔棋吗？将撤销最近两步（你和AI各一步）')) {
        return;
    }

    try {
        const response = await fetch(`/api/game/${gameState.id}/undo`, {
            method: 'POST'
        });

        const result = await response.json();
        
        if (result.success) {
            gameState = result.data;
            selectedPiece = null;
            drawBoard();
            updateGameInfo();
        } else {
            alert('悔棋失败：' + (result.error || '未知错误'));
        }
    } catch (error) {
        console.error('悔棋失败:', error);
        alert('悔棋失败，请重试');
    }
}

// 获取提示
async function getHint() {
    if (!gameState || gameState.status !== 'playing') {
        alert('当前无法获取提示');
        return;
    }

    if (gameState.turn !== gameState.playerColor) {
        alert('现在是AI的回合，无需提示');
        return;
    }

    try {
        const response = await fetch(`/api/game/${gameState.id}/hint`, {
            method: 'POST'
        });

        const result = await response.json();
        
        if (result.success && result.data) {
            const hint = result.data;
            alert(`AI建议：\n从 (${hint.fromRow}, ${hint.fromCol}) 移动到 (${hint.toRow}, ${hint.toCol})\n\n${hint.reason || ''}`);
        } else {
            alert('获取提示失败：' + (result.error || '未知错误'));
        }
    } catch (error) {
        console.error('获取提示失败:', error);
        alert('获取提示失败，请重试');
    }
}
