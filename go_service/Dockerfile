FROM golang:1.24-alpine as builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o greeting-service .
WORKDIR /app/counter
RUN go build -o counter-service .

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache libc6-compat

COPY --from=builder /app/greeting-service .
COPY --from=builder /app/counter/counter-service .

CMD sh -c "./greeting-service & ./counter-service"