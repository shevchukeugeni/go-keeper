FROM golang:1.22-alpine AS builder

RUN apk --update --no-cache add build-base

WORKDIR  $GOPATH/src/go-keeper

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /go/bin/server ./cmd/server
RUN GOOS=windows go build -ldflags "-X 'main.buildVersion=0.0.1' -X 'main.buildTime=$(date)'" -o /go/bin/build/client_win.exe ./cmd/client
RUN GOOS=linux go build -ldflags "-X 'main.buildVersion=0.0.1' -X 'main.buildTime=$(date)'" -o /go/bin/build/client_linux ./cmd/client
RUN GOOS=darwin go build -ldflags "-X 'main.buildVersion=0.0.1' -X 'main.buildTime=$(date)'" -o /go/bin/build/client_darwin ./cmd/client

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /go/bin/ /usr/bin/
COPY --from=builder /go/bin/build /app/build

EXPOSE 8080

CMD ["/usr/bin/server"]