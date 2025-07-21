FROM golang:1.24.5-alpine AS build
WORKDIR /decode_jwt
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -ldflags '-s -w' .

FROM alpine:3.21
COPY --from=build /decode_jwt/decode_jwt /opt/decode_jwt
ENTRYPOINT ["/opt/decode_jwt"]

