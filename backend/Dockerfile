FROM golang:1.27.1 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o backend .

FROM scratch
COPY --from=builder /app/backend /backend
EXPOSE 8080
ENTRYPOINT ["/backend"]
