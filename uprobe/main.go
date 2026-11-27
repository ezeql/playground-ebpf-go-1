// Mini-OBI: attaches a uprobe to net/http.serverHandler.ServeHTTP inside the app
// and reads, from user space, the counter the eBPF program increments in the kernel.
package main

//go:generate go tool bpf2go bpf count.c

import (
	"log"
	"time"

	"github.com/cilium/ebpf/link"
)

func main() {
	var objs bpfObjects
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Fatal(err)
	}

	// With pid: "service:app", PID 1 is the app, so /proc/1/exe is its binary.
	ex, err := link.OpenExecutable("/proc/1/exe")
	if err != nil {
		log.Fatal(err)
	}
	// The program and the uprobe stay loaded for as long as this process runs.
	if _, err := ex.Uprobe("net/http.serverHandler.ServeHTTP", objs.Count, nil); err != nil {
		log.Fatal(err)
	}

	for range time.Tick(time.Second) {
		var n uint64
		if err := objs.Hits.Lookup(uint32(0), &n); err != nil {
			log.Fatal(err)
		}
		log.Printf("HTTP requests served: %d", n)
	}
}
