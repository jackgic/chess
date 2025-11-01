#!/bin/bash

# 中国象棋AI项目状态检查脚本
# 作者: AI Assistant
# 日期: $(date +%Y-%m-%d)

# 项目配置
PROJECT_NAME="chinese-chess-ai"
PID_FILE="${PROJECT_NAME}.pid"
LOG_DIR="logs"
LOG_FILE="${LOG_DIR}/${PROJECT_NAME}.log"
ERROR_LOG="${LOG_DIR}/${PROJECT_NAME}.error.log"
PORT="8090"

echo "=== 中国象棋AI项目状态检查 ==="
echo "检查时间: $(date)"
echo ""

# 检查PID文件
if [ -f ${PID_FILE} ]; then
    PID=$(cat ${PID_FILE})
    echo "PID文件存在: ${PID_FILE}"
    echo "记录的PID: ${PID}"
    
    # 检查进程是否运行
    if ps -p ${PID} > /dev/null 2>&1; then
        echo "✅ 进程状态: 运行中"
        
        # 获取进程信息 (兼容 macOS)
        if [[ "$OSTYPE" == "darwin"* ]]; then
            PROCESS_INFO=$(ps -p ${PID} -o pid,ppid,etime,pcpu,pmem,command)
        else
            PROCESS_INFO=$(ps -p ${PID} -o pid,ppid,etime,pcpu,pmem,cmd --no-headers)
        fi
        echo "进程信息:"
        echo "${PROCESS_INFO}" | sed 's/^/  /'
        
        # 检查端口占用
        if lsof -i :${PORT} > /dev/null 2>&1; then
            echo "✅ 端口状态: ${PORT} 端口已占用"
            PORT_INFO=$(lsof -i :${PORT} | grep LISTEN)
            echo "端口信息: ${PORT_INFO}"
        else
            echo "❌ 端口状态: ${PORT} 端口未占用"
        fi
        
        # 测试API连接
        echo ""
        echo "测试API连接..."
        if curl -s -f http://localhost:${PORT}/api/health > /dev/null 2>&1; then
            echo "✅ API状态: 服务正常响应"
        else
            echo "❌ API状态: 服务无响应"
        fi
        
    else
        echo "❌ 进程状态: 进程不存在"
        echo "PID文件存在但进程已停止，建议清理PID文件"
    fi
else
    echo "❌ PID文件不存在: ${PID_FILE}"
    
    # 检查是否有相关进程在运行
    PIDS=$(pgrep -f ${PROJECT_NAME} 2>/dev/null || true)
    if [ -n "${PIDS}" ]; then
        echo "⚠️  发现运行中的相关进程:"
        for pid in ${PIDS}; do
            if [[ "$OSTYPE" == "darwin"* ]]; then
                PROCESS_INFO=$(ps -p ${pid} -o pid,ppid,etime,pcpu,pmem,command)
            else
                PROCESS_INFO=$(ps -p ${pid} -o pid,ppid,etime,pcpu,pmem,cmd --no-headers)
            fi
            echo "${PROCESS_INFO}" | sed 's/^/  /'
        done
    else
        echo "❌ 进程状态: 没有运行中的进程"
    fi
fi

echo ""
echo "=== 日志文件状态 ==="

# 检查日志目录
if [ -d ${LOG_DIR} ]; then
    echo "✅ 日志目录存在: ${LOG_DIR}"
    
    # 检查主日志文件
    if [ -f ${LOG_FILE} ]; then
        LOG_SIZE=$(du -h ${LOG_FILE} | cut -f1)
        LOG_LINES=$(wc -l < ${LOG_FILE})
        echo "✅ 主日志文件: ${LOG_FILE} (大小: ${LOG_SIZE}, 行数: ${LOG_LINES})"
        
        echo "最近5行日志:"
        tail -5 ${LOG_FILE} | sed 's/^/  /'
    else
        echo "❌ 主日志文件不存在: ${LOG_FILE}"
    fi
    
    # 检查错误日志文件
    if [ -f ${ERROR_LOG} ]; then
        ERROR_SIZE=$(du -h ${ERROR_LOG} | cut -f1)
        ERROR_LINES=$(wc -l < ${ERROR_LOG})
        echo "⚠️  错误日志文件: ${ERROR_LOG} (大小: ${ERROR_SIZE}, 行数: ${ERROR_LINES})"
        
        if [ ${ERROR_LINES} -gt 0 ]; then
            echo "最近5行错误日志:"
            tail -5 ${ERROR_LOG} | sed 's/^/  /'
        fi
    else
        echo "✅ 错误日志文件不存在: ${ERROR_LOG} (无错误)"
    fi
else
    echo "❌ 日志目录不存在: ${LOG_DIR}"
fi

echo ""
echo "=== 系统资源使用 ==="

# 检查系统负载
LOAD_AVG=$(uptime | awk -F'load average:' '{print $2}')
echo "系统负载:${LOAD_AVG}"

# 检查内存使用
if command -v free > /dev/null 2>&1; then
    echo "内存使用:"
    free -h | sed 's/^/  /'
elif command -v vm_stat > /dev/null 2>&1; then
    echo "内存使用 (macOS):"
    vm_stat | head -5 | sed 's/^/  /'
fi

echo ""
echo "=== 可用命令 ==="
echo "启动程序: ./start.sh"
echo "停止程序: ./stop.sh"
echo "查看日志: tail -f ${LOG_FILE}"
echo "查看错误: tail -f ${ERROR_LOG}"
echo "测试API: curl http://localhost:${PORT}/api/health"