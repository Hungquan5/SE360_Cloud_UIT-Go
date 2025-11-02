# 1) Build
FROM golang:1.22 AS builder 
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/user-svc ./cmd/user-svc

# 2) Run
FROM gcr.io/distroless/base-debian12
WORKDIR /srv
ENV ADDR=:8080
EXPOSE 8080
COPY --from=builder /usr/local/bin/user-svc /usr/local/bin/user-svc
# (optional) copy migrations if you want the container to run them itself, otherwise run externally
CMD ["/usr/local/bin/user-svc"]
