package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

type DownloadRequest struct {
	URL string `json:"url"`
}

func main() {
	// Handlers

	// Vue thing must be compiled using yarn run build first
	fs := http.FileServer(http.Dir("./sc-frontend/dist"))
	http.Handle("/", fs)
	http.Handle("/api/download", http.HandlerFunc(downloadHandler))

	// Start the server.
	fmt.Println("Server listening on port 3000")
	log.Panic(
		http.ListenAndServe(":3000", nil),
	)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Triggering downloadHandler...")

	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			fmt.Println("[x] Request was not JSON")
			http.Error(w, "{\"error\": \"Must be JSON\"}", http.StatusUnsupportedMediaType)
		}

		// Parse the request body as JSON
		var requestData DownloadRequest
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			stringErr := fmt.Sprintf("%s", err)
			fmt.Println("[x] Request is invalid JSON")
			http.Error(w, "{\"error\": \"Invalid JSON: +"+stringErr+"\"}", http.StatusBadRequest)
			return
		}

		// Check if the "url" key is present
		targetUrl := requestData.URL
		if targetUrl == "" {
			fmt.Println("[x] Request did not contain 'url' key in JSON")
			http.Error(w, "{\"error\": \"Missing 'url' key in request body\"}", http.StatusBadRequest)
			return
		}

		cmd := exec.Command("scdl -l " + targetUrl + " -a")
		err = cmd.Run()
		if err != nil {
			stringErr := fmt.Sprintf("%s", err)
			fmt.Println("[x] Error downloading: " + stringErr)
			http.Error(w, "{\"error\": \"Error downloading: "+stringErr+"\"}", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "{\"message\": \"Hello, world!\"}")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
