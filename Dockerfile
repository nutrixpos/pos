# Stage 1: Build
FROM golang:alpine AS build
WORKDIR /go/src/app
COPY . ./
RUN go mod download
RUN apk add --no-cache build-base
RUN CGO_ENABLED=1 GOOS=linux go build -o ./pos

# Stage 2: Final
FROM alpine
WORKDIR /app
COPY --from=build /go/src/app/pos .
COPY --from=build /go/src/app/assets .
COPY --from=build /go/src/app/config.yaml .
COPY --from=build /go/src/app/zitadel-key.json .
EXPOSE 8080
CMD ["./pos"]
