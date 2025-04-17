package main

import (
	"embed"
	"fmt"
	"net/http"
	"time"
)

//go:embed index.html
var content embed.FS

func sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, _ := w.(http.Flusher)

	for {
		select {
		case <-r.Context().Done():
			return
		default:
			fmt.Fprintf(w, "data: %s\n\n", time.Now().Format(time.RFC3339))
			flusher.Flush()
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	// 埋め込んだHTMLを配信
	http.Handle("/", http.FileServer(http.FS(content)))
	http.HandleFunc("/sse", sseHandler)

	fmt.Println("Server started at :8080")
	http.ListenAndServe("localhost:8080", nil)
}
