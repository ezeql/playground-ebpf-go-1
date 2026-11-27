FROM golang:1.27-alpine AS build
RUN apk add --no-cache clang llvm libbpf-dev linux-headers
WORKDIR /src
COPY . .
# bpf2go compiles count.c to eBPF bytecode and generates the Go code that loads it.
RUN go generate && CGO_ENABLED=0 go build -o /uprobe

FROM scratch
COPY --from=build /uprobe /uprobe
ENTRYPOINT ["/uprobe"]
