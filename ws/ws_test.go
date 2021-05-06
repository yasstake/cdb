package ws

import (
	"net/http"
	"testing"
)

func TestStartServer(t *testing.T) {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.HandleFunc("/api/tran", TransactionHandler)
	http.HandleFunc("/html/info", HtmlInfoHandler)

	server.ListenAndServe()
}
