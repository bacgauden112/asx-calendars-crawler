# Bước 1: Xây dựng ứng dụng
FROM golang:1.22.4 AS builder

# Thiết lập thư mục làm việc
WORKDIR /app

# Sao chép mã nguồn vào container
COPY . .

# Sao chép thư mục templates vào image

RUN export GIN_MODE=release

# Tải về các module cần thiết
RUN go mod download

# Biên dịch ứng dụng
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Bước 2: Tạo image cuối cùng
FROM alpine:latest  

RUN export GIN_MODE=release
# Cài đặt các gói cần thiết
RUN apk --no-cache add ca-certificates
RUN apk add chromium
WORKDIR /root/

# Sao chép file thực thi từ bước builder
COPY --from=builder /app/main .
COPY --from=builder /app/templates/* ./templates/

# Expose port
EXPOSE 8080

# Chạy ứng dụng
CMD ["./main"]