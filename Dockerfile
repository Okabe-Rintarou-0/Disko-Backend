FROM golang:1.19 as BUILDER

WORKDIR /build

COPY . .

RUN GOPROXY=https://goproxy.cn CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ./bin/app .

FROM alpine

WORKDIR /build

COPY --from=builder /build .

EXPOSE 8888

RUN chmod +x ./bin/app

ENTRYPOINT ["./bin/app"]