# Stage 1: Build
FROM golang:1.24.4-alpine AS build
WORKDIR /go/src/app
COPY . ./
RUN go mod download
RUN apk add --no-cache build-base
RUN CGO_ENABLED=1 GOOS=linux go build -o ./pos

# Stage 2: Final
FROM alpine
RUN apk add --no-cache chromium font-noto font-noto-arabic
# Receipt rendering needs a headless Chromium browser. The container runs as
# root, where Chromium's sandbox cannot initialize.
ENV PRINTING_NO_SANDBOX=true
WORKDIR /app
COPY --from=build /go/src/app/pos .
COPY --from=build /go/src/app/assets ./assets/
COPY --from=build /go/src/app/config.example.yaml ./config.yaml
EXPOSE 8000
CMD ["./pos"]