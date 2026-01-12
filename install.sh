#!/bin/bash

# aicli 安装脚本
# 支持 Linux 和 macOS

set -e

# 颜色输出
RED='\033[0.31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 版本信息
VERSION=${VERSION:-"latest"}
REPO="studyzy/aicli"

# 检测操作系统和架构
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case $OS in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        *)
            echo -e "${RED}错误: 不支持的操作系统 $OS${NC}"
            exit 1
            ;;
    esac
    
    case $ARCH in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            echo -e "${RED}错误: 不支持的架构 $ARCH${NC}"
            exit 1
            ;;
    esac
    
    echo -e "${GREEN}检测到平台: $OS/$ARCH${NC}"
}

# 检查依赖
check_dependencies() {
    if ! command -v curl &> /dev/null; then
        echo -e "${RED}错误: curl 未安装${NC}"
        exit 1
    fi
    
    if ! command -v tar &> /dev/null; then
        echo -e "${RED}错误: tar 未安装${NC}"
        exit 1
    fi
}

# 下载并安装
install_aicli() {
    # 临时目录
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"
    
    echo -e "${YELLOW}正在下载 aicli...${NC}"
    
    # 构建下载 URL
    if [ "$VERSION" = "latest" ]; then
        DOWNLOAD_URL="https://github.com/$REPO/releases/latest/download/aicli-${OS}-${ARCH}.tar.gz"
    else
        DOWNLOAD_URL="https://github.com/$REPO/releases/download/${VERSION}/aicli-${OS}-${ARCH}.tar.gz"
    fi
    
    echo "下载地址: $DOWNLOAD_URL"
    
    # 下载文件
    if ! curl -L -o aicli.tar.gz "$DOWNLOAD_URL"; then
        echo -e "${RED}错误: 下载失败${NC}"
        echo -e "${YELLOW}提示: 请检查网络连接或手动从源码编译${NC}"
        exit 1
    fi
    
    echo -e "${YELLOW}正在解压...${NC}"
    tar -xzf aicli.tar.gz
    
    # 安装目录
    INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
    
    echo -e "${YELLOW}正在安装到 $INSTALL_DIR...${NC}"
    
    # 检查是否需要 sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv aicli "$INSTALL_DIR/"
    else
        echo -e "${YELLOW}需要 sudo 权限安装到 $INSTALL_DIR${NC}"
        sudo mv aicli "$INSTALL_DIR/"
    fi
    
    # 添加执行权限
    chmod +x "$INSTALL_DIR/aicli"
    
    # 清理
    cd -
    rm -rf "$TMP_DIR"
    
    echo -e "${GREEN}✓ aicli 安装成功！${NC}"
}

# 创建示例配置
create_example_config() {
    CONFIG_FILE="$HOME/.aicli.json"
    
    if [ -f "$CONFIG_FILE" ]; then
        echo -e "${YELLOW}配置文件已存在: $CONFIG_FILE${NC}"
        return
    fi
    
    read -p "是否创建示例配置文件？(y/n): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        cat > "$CONFIG_FILE" << 'EOF'
{
  "version": "1.0",
  "llm": {
    "provider": "openai",
    "api_key": "your-api-key-here",
    "model": "gpt-4",
    "timeout": 10
  },
  "execution": {
    "auto_confirm": false,
    "timeout": 30
  },
  "safety": {
    "enable_checks": true,
    "require_confirmation": true
  },
  "history": {
    "enabled": true,
    "max_entries": 1000
  }
}
EOF
        chmod 600 "$CONFIG_FILE"
        echo -e "${GREEN}✓ 配置文件已创建: $CONFIG_FILE${NC}"
        echo -e "${YELLOW}请编辑配置文件并设置您的 API 密钥${NC}"
    fi
}

# 显示使用说明
show_usage() {
    echo ""
    echo -e "${GREEN}=== aicli 安装完成 ===${NC}"
    echo ""
    echo "使用方法:"
    echo "  aicli \"列出当前目录文件\""
    echo "  aicli --help"
    echo ""
    echo "配置文件:"
    echo "  ~/.aicli.json"
    echo ""
    echo "更多信息:"
    echo "  https://github.com/$REPO"
    echo ""
}

# 主函数
main() {
    echo -e "${GREEN}=== aicli 安装程序 ===${NC}"
    echo ""
    
    detect_platform
    check_dependencies
    install_aicli
    create_example_config
    show_usage
}

main "$@"
