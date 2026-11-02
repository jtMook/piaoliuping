#!/bin/bash

# 漂流瓶开发环境一键管理脚本
# 用法: ./dev.sh [start|stop|restart|status]

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 项目路径
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
SERVER_DIR="$PROJECT_ROOT/piaoliu-v2-server"
FRONTEND_DIR="$PROJECT_ROOT/piaoliu-v2-miniprogram"

# PID 文件
SERVER_PID_FILE="$PROJECT_ROOT/.server.pid"
FRONTEND_PID_FILE="$PROJECT_ROOT/.frontend.pid"

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 启动后端服务
start_server() {
    log_info "启动后端服务..."
    cd "$SERVER_DIR"
    
    # 启动 Docker 容器（MySQL + Redis）
    docker compose up -d mysql redis
    
    # 等待 MySQL 就绪
    log_info "等待 MySQL 启动..."
    for i in {1..30}; do
        if docker compose exec -T mysql mysqladmin ping -h localhost -uroot -proot_password --silent 2>/dev/null; then
            log_info "MySQL 已就绪"
            break
        fi
        if [ $i -eq 30 ]; then
            log_error "MySQL 启动超时"
            exit 1
        fi
        sleep 1
    done
    
    # 启动 Go 后端服务
    nohup go run cmd/server/main.go > "$PROJECT_ROOT/server.log" 2>&1 &
    echo $! > "$SERVER_PID_FILE"
    log_info "后端服务已启动 (PID: $(cat $SERVER_PID_FILE))"
    
    # 等待后端服务就绪
    log_info "等待后端服务启动..."
    for i in {1..20}; do
        if lsof -i :8080 -P 2>/dev/null | grep LISTEN > /dev/null 2>&1; then
            log_info "后端服务已就绪"
            break
        fi
        if [ $i -eq 20 ]; then
            log_warn "后端服务可能启动失败，请查看日志: tail -f $PROJECT_ROOT/server.log"
        fi
        sleep 1
    done
}

# 停止后端服务
stop_server() {
    log_info "停止后端服务..."
    
    if [ -f "$SERVER_PID_FILE" ]; then
        PID=$(cat "$SERVER_PID_FILE")
        if kill -0 "$PID" 2>/dev/null; then
            kill "$PID"
            log_info "后端服务已停止 (PID: $PID)"
        else
            log_warn "后端服务未运行 (PID: $PID)"
        fi
        rm -f "$SERVER_PID_FILE"
    else
        log_warn "未找到后端服务 PID 文件"
    fi
    
    # 停止 Docker 容器
    cd "$SERVER_DIR"
    docker compose down
    log_info "Docker 容器已停止"
}

# 启动前端服务
start_frontend() {
    log_info "启动前端服务..."
    cd "$FRONTEND_DIR"
    
    nohup npm run dev:h5 > "$PROJECT_ROOT/frontend.log" 2>&1 &
    echo $! > "$FRONTEND_PID_FILE"
    log_info "前端服务已启动 (PID: $(cat $FRONTEND_PID_FILE))"
    
    # 等待前端服务就绪
    log_info "等待前端服务启动..."
    for i in {1..15}; do
        if curl -s http://localhost:5173 > /dev/null 2>&1; then
            log_info "前端服务已就绪"
            break
        fi
        if [ $i -eq 15 ]; then
            log_warn "前端服务可能启动失败，请查看日志: tail -f $PROJECT_ROOT/frontend.log"
        fi
        sleep 1
    done
}

# 停止前端服务
stop_frontend() {
    log_info "停止前端服务..."
    
    if [ -f "$FRONTEND_PID_FILE" ]; then
        PID=$(cat "$FRONTEND_PID_FILE")
        if kill -0 "$PID" 2>/dev/null; then
            kill "$PID"
            log_info "前端服务已停止 (PID: $PID)"
        else
            log_warn "前端服务未运行 (PID: $PID)"
        fi
        rm -f "$FRONTEND_PID_FILE"
    else
        log_warn "未找到前端服务 PID 文件"
    fi
}

# 查看服务状态
check_status() {
    echo ""
    echo "========== 服务状态 =========="
    echo ""
    
    # 后端状态
    if [ -f "$SERVER_PID_FILE" ]; then
        PID=$(cat "$SERVER_PID_FILE")
        if kill -0 "$PID" 2>/dev/null; then
            echo -e "后端服务: ${GREEN}运行中${NC} (PID: $PID)"
        else
            echo -e "后端服务: ${RED}已停止${NC} (PID 文件存在但进程不存在)"
        fi
    else
        echo -e "后端服务: ${YELLOW}未启动${NC}"
    fi
    
    # 前端状态
    if [ -f "$FRONTEND_PID_FILE" ]; then
        PID=$(cat "$FRONTEND_PID_FILE")
        if kill -0 "$PID" 2>/dev/null; then
            echo -e "前端服务: ${GREEN}运行中${NC} (PID: $PID)"
        else
            echo -e "前端服务: ${RED}已停止${NC} (PID 文件存在但进程不存在)"
        fi
    else
        echo -e "前端服务: ${YELLOW}未启动${NC}"
    fi
    
    # Docker 状态
    echo ""
    echo "Docker 容器状态:"
    cd "$SERVER_DIR"
    docker compose ps 2>/dev/null || echo "Docker 未运行"
    
    echo ""
    echo "========== 访问地址 =========="
    echo "后端 API:  http://localhost:8080"
    echo "前端 H5:   http://localhost:5173"
    echo "日志文件:  $PROJECT_ROOT/server.log"
    echo "           $PROJECT_ROOT/frontend.log"
    echo ""
}

# 主逻辑
case "$1" in
    start)
        echo ""
        echo "========== 启动漂流瓶开发环境 =========="
        echo ""
        start_server
        start_frontend
        echo ""
        echo -e "${GREEN}========== 启动完成 ==========${NC}"
        check_status
        ;;
    
    stop)
        echo ""
        echo "========== 停止漂流瓶开发环境 =========="
        echo ""
        stop_frontend
        stop_server
        echo ""
        echo -e "${GREEN}========== 停止完成 ==========${NC}"
        ;;
    
    restart)
        echo ""
        echo "========== 重启漂流瓶开发环境 =========="
        echo ""
        stop_frontend
        stop_server
        sleep 2
        start_server
        start_frontend
        echo ""
        echo -e "${GREEN}========== 重启完成 ==========${NC}"
        check_status
        ;;
    
    status)
        check_status
        ;;
    
    *)
        echo ""
        echo "漂流瓶开发环境管理脚本"
        echo ""
        echo "用法: $0 [start|stop|restart|status]"
        echo ""
        echo "命令说明:"
        echo "  start    启动后端和前端服务"
        echo "  stop     停止后端和前端服务"
        echo "  restart  重启后端和前端服务"
        echo "  status   查看服务状态"
        echo ""
        echo "示例:"
        echo "  $0 start    # 启动开发环境"
        echo "  $0 stop     # 停止开发环境"
        echo "  $0 status   # 查看服务状态"
        echo ""
        exit 1
        ;;
esac
