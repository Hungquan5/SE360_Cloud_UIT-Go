# syntax=docker/dockerfile:1
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/driver-svc ./cmd/driver-svc

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=build /bin/driver-svc /app/driver-svc
ENV APP_PORT=8080
EXPOSE 8080
USER nonroot
ENTRYPOINT ["/app/driver-svc"]
