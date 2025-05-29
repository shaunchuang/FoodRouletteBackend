# Build stage
FROM golang:1.21-alpine AS builder

# 設定工作目錄
WORKDIR /app

# 安裝 git（某些 Go 模組可能需要）
RUN apk add --no-cache git

# 複製 go mod 檔案
COPY go.mod go.sum ./

# 下載相依套件
RUN go mod download

# 複製原始碼
COPY . .

# 編譯應用程式
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Runtime stage
FROM alpine:latest

# 安裝 ca-certificates（用於 HTTPS 請求）
RUN apk --no-cache add ca-certificates tzdata

# 設定時區
ENV TZ=Asia/Taipei

# 建立非 root 使用者
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 設定工作目錄
WORKDIR /root/

# 從 builder stage 複製編譯好的執行檔
COPY --from=builder /app/main .

# 複製設定檔案範例
COPY --from=builder /app/.env.example .

# 變更檔案擁有者
RUN chown -R appuser:appgroup /root/

# 切換到非 root 使用者
USER appuser

# 暴露埠號
EXPOSE 8080

# 執行應用程式
CMD ["./main"]