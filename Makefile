# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=food-roulette-server
BINARY_UNIX=$(BINARY_NAME)_unix

# Main package path
MAIN_PATH=./cmd/server

.PHONY: all build clean test coverage help run dev

all: test build

## build: 編譯應用程式
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

## clean: 清理編譯產物
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

## test: 執行測試
test:
	$(GOTEST) -v ./...

## coverage: 執行測試並產生覆蓋率報告
coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

## run: 編譯並執行應用程式
run: build
	./$(BINARY_NAME)

## dev: 開發模式執行（直接執行原始碼）
dev:
	$(GOCMD) run $(MAIN_PATH)/main.go

## deps: 下載相依套件
deps:
	$(GOMOD) download
	$(GOMOD) tidy

## build-linux: 為 Linux 平台交叉編譯
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PATH)

## docker-build: 建立 Docker 映像
docker-build:
	docker build -t food-roulette-backend .

## migrate-up: 執行資料庫遷移（向上）
migrate-up:
	migrate -path migrations -database "postgres://postgres:password@localhost/food_roulette?sslmode=disable" up

## migrate-down: 執行資料庫遷移（向下）
migrate-down:
	migrate -path migrations -database "postgres://postgres:password@localhost/food_roulette?sslmode=disable" down

## migrate-create: 建立新的遷移檔案
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations $$name

## help: 顯示說明
help: Makefile
	@echo "可用指令："
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'