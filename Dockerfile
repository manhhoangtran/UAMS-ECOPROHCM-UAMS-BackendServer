## Build
FROM golang:1.17.2-alpine AS builder

ARG APP_NAME=dms-be

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o ./$APP_NAME

# Run
FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/$APP_NAME .
COPY --from=builder /app/certs /certs
RUN apk add --no-cache tzdata
RUN cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime
EXPOSE 8080

CMD /app/dms-be