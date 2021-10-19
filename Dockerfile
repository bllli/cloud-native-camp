FROM golang:1.16-alpine AS build
ENV GO111MODULE=on GOPROXY=https://goproxy.cn
COPY go.mod /app/go.mod
WORKDIR /app
RUN go mod download
COPY . /app
RUN export CGO_ENABLED=0 && go build -o demo-server

FROM scratch
EXPOSE 8888
COPY --from=build /app/demo-server /usr/local/bin/demo-server
ENTRYPOINT ["demo-server"]