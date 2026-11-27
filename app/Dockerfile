FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
# No -ldflags "-s -w": OBI needs the Go binary's symbols for its uprobes.
RUN CGO_ENABLED=0 go build -o /app .

FROM scratch
COPY --from=build /app /app
ENTRYPOINT ["/app"]
