FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o flight-sorter .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/flight-sorter .
ENV PORT=8080
EXPOSE 8080
CMD ["./flight-sorter"]
