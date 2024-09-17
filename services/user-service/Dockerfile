FROM docker.io/golang:1.22 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go vet -v ./src/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

FROM gcr.io/distroless/static-debian12
COPY --from=build /app/main /main
WORKDIR /app
EXPOSE 50051 8000
ENTRYPOINT ["/main"]