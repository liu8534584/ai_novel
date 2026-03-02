.PHONY: dev-back dev-front build install

# 启动后端服务
dev-back:
	go run main.go

# 启动前端服务
dev-front:
	cd frontend && npm run dev

# 安装所有依赖
install:
	go mod tidy
	cd frontend && npm install

# 编译项目
build:
	go build -o ai_novel_app main.go
	cd frontend && npm run build
