#!/bin/bash

# 确保 swag 命令已安装
if ! command -v swag &> /dev/null; then
    echo "Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# 生成 swagger 文档
echo "Generating Swagger documentation..."
swag init -g main.go --parseDependency --parseInternal

echo "Swagger documentation generated successfully!"