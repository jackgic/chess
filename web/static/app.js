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
    document.querySelectorAll('.color-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            document.querySelectorAll('.color-btn').forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
        });
    });

    // 新游戏按钮
    const newGameBtn = document.getElementById('newGameBtn');
    if (newGameBtn) {
        newGameBtn.addEventListener('click', startNewGame);
    }

    // 棋盘点击事件
    canvas.addEventListener('click', handleBoardClick);
    
    console.log('前端初始化完成');
});

// 开始新游戏
async function startNewGame() {
    const playerColor = parseInt(document.querySelector('.color-btn.active').dataset.color);
    console.log('开始新游戏，玩家颜色:', playerColor);

    try {
        const response = await fetch('/api/game/new', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ playerColor })
        });

        const result = await response.json();
        console.log('游戏创建响应:', result);
        
        if (result.success) {
            gameState = result.data;
            selectedPiece = null;
            console.log('游戏状态设置完成:', gameState);

            // 显示游戏界面
            document.getElementById('gameInfo').style.display = 'flex';
            document.getElementById('gameArea').style.display = 'flex';
            document.getElementById('moveHistory').style.display = 'block';

            // 绘制棋盘
            drawBoard();
            updateGameInfo();

            // 如果AI先手，让AI走第一步
            if (gameState.turn === gameState.aiColor) {
                console.log('AI先手，准备AI走子');
                setTimeout(() => requestAIMove(), 500);
            } else {
                console.log('玩家先手，等待玩家走子');
            }
        } else {
            alert('创建游戏失败：' + result.error);
        }
    } catch (error) {
        console.error('创建游戏失败:', error);
        alert('创建游戏失败，请检查服务器连接');
    }
}

// 绘制棋盘
function drawBoard() {
    // 清空画布
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // 绘制背景
    ctx.fillStyle = '#F4A460';
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    // 绘制棋盘线
    ctx.strokeStyle = '#000';
    ctx.lineWidth = 2;

    // 横线
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
            ctx.lineTo(BOARD_PADDING + i * CELL_SIZE, BOARD_PADDING + 9 * CELL_SIZE);
        } else {
            ctx.lineTo(BOARD_PADDING + i * CELL_SIZE, BOARD_PADDING + 4 * CELL_SIZE);
            ctx.moveTo(BOARD_PADDING + i * CELL_SIZE, BOARD_PADDING + 5 * CELL_SIZE);
            ctx.lineTo(BOARD_PADDING + i * CELL_SIZE, BOARD_PADDING + 9 * CELL_SIZE);
        }
        ctx.stroke();
    }

    // 绘制九宫格斜线
    drawDiagonal(3, 0, 5, 2);
    drawDiagonal(5, 0, 3, 2);
    drawDiagonal(3, 7, 5, 9);
    drawDiagonal(5, 7, 3, 9);

    // 绘制楚河汉界
    ctx.font = 'bold 24px Arial';
    ctx.fillStyle = '#000';
    ctx.textAlign = 'center';
    ctx.fillText('楚河', BOARD_PADDING + 2 * CELL_SIZE, BOARD_PADDING + 4.7 * CELL_SIZE);
    ctx.fillText('汉界', BOARD_PADDING + 6 * CELL_SIZE, BOARD_PADDING + 4.7 * CELL_SIZE);

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

    // 绘制选中标记
    if (selectedPiece) {
        drawSelection(selectedPiece.row, selectedPiece.col);
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
    const x = BOARD_PADDING + col * CELL_SIZE;
    const y = BOARD_PADDING + row * CELL_SIZE;

    // 绘制棋子圆形
    ctx.beginPath();
    ctx.arc(x, y, PIECE_RADIUS, 0, 2 * Math.PI);
    ctx.fillStyle = piece.color === 1 ? '#FFE4B5' : '#2F4F4F';
    ctx.fill();
    ctx.strokeStyle = '#000';
    ctx.lineWidth = 2;
    ctx.stroke();

    // 绘制棋子文字
    ctx.font = 'bold 24px "Microsoft YaHei", "SimHei", "黑体", Arial, sans-serif';
    ctx.fillStyle = piece.color === 1 ? '#DC143C' : '#000';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    
    // 确保棋子名称存在
    const pieceName = piece.name || '?';
    ctx.fillText(pieceName, x, y);
    
    if (!piece.name) {
        console.warn(`棋子在位置 (${row},${col}) 没有名称:`, piece);
    }
}

// 绘制选中标记
function drawSelection(row, col) {
    const x = BOARD_PADDING + col * CELL_SIZE;
    const y = BOARD_PADDING + row * CELL_SIZE;

    ctx.strokeStyle = '#00FF00';
    ctx.lineWidth = 3;
    ctx.beginPath();
    ctx.arc(x, y, PIECE_RADIUS + 5, 0, 2 * Math.PI);
    ctx.stroke();
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

    // 计算点击的格子
    const col = Math.round((x - BOARD_PADDING) / CELL_SIZE);
    const row = Math.round((y - BOARD_PADDING) / CELL_SIZE);

    console.log('点击坐标:', { x, y, row, col });

    if (row < 0 || row > 9 || col < 0 || col > 8) {
        console.log('点击位置超出棋盘范围');
        return;
    }

    const piece = gameState.board[row][col];
    console.log('点击的棋子:', piece);

    if (!selectedPiece) {
        // 选择棋子
        if (piece.type !== 0 && piece.color === gameState.playerColor) {
            console.log('选择棋子:', piece);
            selectedPiece = { row, col };
            drawBoard();
        } else {
            console.log('无法选择棋子 - type:', piece.type, 'color:', piece.color, 'playerColor:', gameState.playerColor);
        }
    } else {
        // 移动棋子
        if (row === selectedPiece.row && col === selectedPiece.col) {
            // 取消选择
            console.log('取消选择棋子');
            selectedPiece = null;
            drawBoard();
        } else {
            // 尝试移动
            console.log('尝试移动棋子从', selectedPiece, '到', { row, col });
            playerMove(selectedPiece.row, selectedPiece.col, row, col);
        }
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
                document.getElementById('aiAnswer').textContent = result.data.answer;
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
            startNewGame();
        }
    }, 500);
}
