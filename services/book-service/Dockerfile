FROM docker.io/golang:1.22 AS build
WORKDIR /app
COPY ../../go.mod ../../go.sum ./
RUN go mod download && go mod verify
COPY ../../ ./
WORKDIR /app/services/book-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .
RUN chmod +x /main

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /main /app/main
EXPOSE 8001
ENTRYPOINT ["/app/main"]