#!/bin/bash

# 版本发布脚本 - 自动提交代码并创建新版本tag
# 使用方法: 
#   ./release.sh patch   # 增加补丁版本 (0.1.2 -> 0.1.3)
#   ./release.sh minor   # 增加次版本 (0.1.2 -> 0.2.0)  
#   ./release.sh major   # 增加主版本 (0.1.2 -> 1.0.0)
#   ./release.sh         # 默认增加补丁版本

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查git状态
check_git_status() {
    print_info "检查Git状态..."
    
    # 检查是否在git仓库中
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        print_error "当前目录不是Git仓库"
        exit 1
    fi
    
    # 检查是否有未跟踪的文件或未提交的更改
    if ! git diff-index --quiet HEAD --; then
        print_warning "检测到未提交的更改"
        git status --porcelain
        echo
        read -p "是否要提交这些更改? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            return 0
        else
            print_error "请先提交或撤销更改后再运行此脚本"
            exit 1
        fi
    else
        print_success "工作目录干净，没有未提交的更改"
    fi
}

# 获取当前最新版本
get_latest_version() {
    print_info "获取当前最新版本..."
    
    # 获取最新的tag，按版本号排序
    local latest_tag=$(git tag --list --sort=-version:refname | head -n1)
    
    if [ -z "$latest_tag" ]; then
        print_warning "没有找到现有的tag，将从v0.1.0开始"
        echo "v0.1.0"
    else
        print_info "当前最新版本: $latest_tag"
        echo "$latest_tag"
    fi
}

# 解析版本号
parse_version() {
    local version=$1
    # 移除v前缀
    version=${version#v}
    
    # 分割版本号
    IFS='.' read -ra VERSION_PARTS <<< "$version"
    
    MAJOR=${VERSION_PARTS[0]:-0}
    MINOR=${VERSION_PARTS[1]:-1}
    PATCH=${VERSION_PARTS[2]:-0}
}

# 递增版本号
increment_version() {
    local bump_type=$1
    
    case $bump_type in
        "major")
            MAJOR=$((MAJOR + 1))
            MINOR=0
            PATCH=0
            ;;
        "minor")
            MINOR=$((MINOR + 1))
            PATCH=0
            ;;
        "patch"|*)
            PATCH=$((PATCH + 1))
            ;;
    esac
    
    NEW_VERSION="v${MAJOR}.${MINOR}.${PATCH}"
}

# 生成提交消息
generate_commit_message() {
    local version=$1
    echo "Release $version: 版本更新和功能改进"
}

# 提交更改
commit_changes() {
    local commit_message=$1
    
    print_info "添加所有更改到暂存区..."
    git add .
    
    print_info "提交更改..."
    git commit -m "$commit_message"
    
    print_success "更改已提交"
}

# 创建并推送tag
create_and_push_tag() {
    local version=$1
    local commit_message=$1
    
    print_info "创建tag: $version"
    git tag -a "$version" -m "$commit_message"
    
    print_info "推送代码到远程仓库..."
    git push origin main
    
    print_info "推送tag到远程仓库..."
    git push origin "$version"
    
    print_success "Tag $version 已成功创建并推送到远程仓库"
}

# 显示发布信息
show_release_info() {
    local old_version=$1
    local new_version=$2
    
    echo
    print_success "🎉 发布完成!"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo -e "  ${BLUE}旧版本:${NC} $old_version"
    echo -e "  ${GREEN}新版本:${NC} $new_version"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo
    echo "现在可以通过以下方式使用新版本:"
    echo -e "  ${YELLOW}go get github.com/yalks/wallet@$new_version${NC}"
    echo
}

# 主函数
main() {
    local bump_type=${1:-patch}
    
    print_info "开始版本发布流程..."
    print_info "版本递增类型: $bump_type"
    echo
    
    # 检查参数
    if [[ ! "$bump_type" =~ ^(patch|minor|major)$ ]]; then
        print_error "无效的版本递增类型: $bump_type"
        echo "支持的类型: patch, minor, major"
        exit 1
    fi
    
    # 检查git状态
    check_git_status
    
    # 获取当前版本
    local current_version=$(get_latest_version)
    
    # 解析版本号
    parse_version "$current_version"
    
    # 递增版本号
    increment_version "$bump_type"
    
    print_info "版本更新: $current_version -> $NEW_VERSION"
    echo
    
    # 确认发布
    read -p "确认要发布版本 $NEW_VERSION 吗? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_warning "发布已取消"
        exit 0
    fi
    
    # 生成提交消息
    local commit_message=$(generate_commit_message "$NEW_VERSION")
    
    # 如果有未提交的更改，先提交
    if ! git diff-index --quiet HEAD --; then
        commit_changes "$commit_message"
    fi
    
    # 创建并推送tag
    create_and_push_tag "$NEW_VERSION" "$commit_message"
    
    # 显示发布信息
    show_release_info "$current_version" "$NEW_VERSION"
}

# 显示帮助信息
show_help() {
    echo "版本发布脚本"
    echo
    echo "使用方法:"
    echo "  $0 [patch|minor|major]"
    echo
    echo "参数说明:"
    echo "  patch  - 递增补丁版本 (默认) 例: 0.1.2 -> 0.1.3"
    echo "  minor  - 递增次版本        例: 0.1.2 -> 0.2.0"
    echo "  major  - 递增主版本        例: 0.1.2 -> 1.0.0"
    echo
    echo "示例:"
    echo "  $0        # 创建补丁版本"
    echo "  $0 patch  # 创建补丁版本"
    echo "  $0 minor  # 创建次版本"
    echo "  $0 major  # 创建主版本"
}

# 检查是否请求帮助
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    show_help
    exit 0
fi

# 运行主函数
main "$@"
