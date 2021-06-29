package api

import (
	"fmt"
	"net/http"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, time.Now().Format("2006-01-02 15:04:05"))
}
