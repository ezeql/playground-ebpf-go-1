// Toy service: /hello calls /work over HTTP, and /work uses MongoDB.
// No OTel SDK: OBI instruments it with eBPF, Mongo operations included.
package main

import (
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	// Connect doesn't dial yet: if Mongo is slow to start, the first requests return 500.
	// zstd on purpose: OBI still sees the operations because for Go it hooks the driver, not the wire bytes.
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://mongo:27017/?compressors=zstd"))
	if err != nil {
		log.Fatal(err)
	}
	hits := client.Database("demo").Collection("hits")

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://localhost:8080/work")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	// One insert and one find in Mongo, plus random latency and errors for the RED metrics.
	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		if _, err := hits.InsertOne(r.Context(), bson.M{"at": time.Now()}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := hits.FindOne(r.Context(), bson.M{}).Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		time.Sleep(time.Duration(rand.IntN(300)) * time.Millisecond)
		if rand.IntN(10) == 0 {
			http.Error(w, "boom", http.StatusInternalServerError)
			return
		}
		io.WriteString(w, "ok\n")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
