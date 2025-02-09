FROM golang:1.23-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN apk add --no-cache build-base

RUN CGO_ENABLED=1 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest

ENV GO111MODULE=on
RUN CGO_ENABLED=1 go build -gcflags "all=-N -l" -o pos .

EXPOSE 8000 4000
CMD [ "/go/bin/dlv", "exec", "/app/pos", "--listen=:4000", "--headless=true", "--log=true", "--accept-multiclient", "--api-version=2" ]