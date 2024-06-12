FROM golang:1.21-alpine as build-base
RUN apk --no-cache add tzdata

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -v main.go

RUN go build -o my-app main.go

#======================

FROM alpine:3.17

COPY --from=build-base /app/my-app /app/my-app
COPY --from=build-base /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ="Asia/Bangkok"
CMD ["/app/my-app"]