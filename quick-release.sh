#!/bin/bash

# 快速发布脚本 - 一键提交并创建新的补丁版本
# 自动递增最后一位版本号 (例: v0.1.2 -> v0.1.3)

set -e

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}🚀 快速发布工具${NC}"
echo "=================================="

# 检查是否在git仓库中
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}❌ 错误: 当前目录不是Git仓库${NC}"
    exit 1
fi

# 获取当前最新版本
echo -e "${BLUE}📋 获取当前版本...${NC}"
CURRENT_VERSION=$(git tag --list --sort=-version:refname | head -n1)

if [ -z "$CURRENT_VERSION" ]; then
    echo -e "${YELLOW}⚠️  没有找到现有tag，将创建 v0.1.0${NC}"
    NEW_VERSION="v0.1.0"
else
    echo -e "   当前版本: ${YELLOW}$CURRENT_VERSION${NC}"
    
    # 解析版本号并递增补丁版本
    VERSION_NUM=${CURRENT_VERSION#v}
    IFS='.' read -ra PARTS <<< "$VERSION_NUM"
    MAJOR=${PARTS[0]:-0}
    MINOR=${PARTS[1]:-1}
    PATCH=${PARTS[2]:-0}
    
    # 递增补丁版本
    PATCH=$((PATCH + 1))
    NEW_VERSION="v${MAJOR}.${MINOR}.${PATCH}"
fi

echo -e "   新版本: ${GREEN}$NEW_VERSION${NC}"
echo

# 检查是否有未提交的更改
if ! git diff-index --quiet HEAD --; then
    echo -e "${BLUE}📝 检测到未提交的更改:${NC}"
    git status --porcelain
    echo
    
    # 自动提交所有更改
    echo -e "${BLUE}💾 提交所有更改...${NC}"
    git add .
    git commit -m "Release $NEW_VERSION: 自动提交版本更新"
    echo -e "${GREEN}✅ 更改已提交${NC}"
else
    echo -e "${GREEN}✅ 工作目录干净${NC}"
fi

# 创建tag
echo -e "${BLUE}🏷️  创建tag: $NEW_VERSION${NC}"
git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION: 自动版本发布"

# 推送到远程
echo -e "${BLUE}📤 推送到远程仓库...${NC}"
git push origin main
git push origin "$NEW_VERSION"

# 完成
echo
echo -e "${GREEN}🎉 发布完成!${NC}"
echo "=================================="
echo -e "版本: ${CURRENT_VERSION:-'无'} → ${GREEN}$NEW_VERSION${NC}"
echo
echo -e "${YELLOW}📦 使用新版本:${NC}"
echo "go get github.com/yalks/wallet@$NEW_VERSION"
echo
