package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// ----------------------------------------
// STRUCTS
// ----------------------------------------

type File struct {
	Name string
	URL  string
}

type DownloadRequest struct {
	URL string `json:"url"`
}

// ----------------------------------------
// FUNCS
// ----------------------------------------

func main() {
	// Handlers
	http.Handle("/", http.HandlerFunc(indexHandler))
	http.Handle("/files", http.HandlerFunc(mp3fileHandler))
	http.Handle("/api/download", http.HandlerFunc(downloadHandler))
	http.Handle("/style.css", http.FileServer(http.Dir("web")))
	http.HandleFunc("/download/", handleDownload)

	// Start the server.
	fmt.Println("Server listening on port 8080")
	log.Panic(
		http.ListenAndServe(":8080", nil),
	)
}

// ----------------------------------------
// HANDLERS
// ----------------------------------------

func mp3fileHandler(w http.ResponseWriter, r *http.Request) {
	// Get mp3 files
	files, err := ioutil.ReadDir("./destination")
	if err != nil {
		stringErr := fmt.Sprintf("%s", err)
		http.Error(w, "{\"error\": \"Could not obtain MP3s: "+stringErr+"\"}", http.StatusInternalServerError)
		return
	}

	mp3Files := make([]File, 0)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".mp3") {
			fileURL := "/download/" + file.Name()
			mp3Files = append(mp3Files, File{Name: file.Name(), URL: fileURL})
		}
	}

	// Parse the HTML template
	tmpl, err := template.ParseFiles("web/file_index.html")
	if err != nil {
		http.Error(w, "Oopsie woopsie", http.StatusInternalServerError)
		log.Fatal(err)
	}

	// Execute the template with the file list
	err = tmpl.Execute(w, mp3Files)
	if err != nil {
		http.Error(w, "Oopsie woopsie", http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the template file
	tmpl, err := template.ParseFiles("web/index.html")
	if err != nil {
		stringErr := fmt.Sprintf("%s", err)
		http.Error(w, "{\"error\": \"Could not template HTML: "+stringErr+"\"}", http.StatusInternalServerError)
		return
	}

	// Execute the template with the data
	err = tmpl.Execute(w, nil)
	if err != nil {
		stringErr := fmt.Sprintf("%s", err)
		http.Error(w, "{\"error\": \"Could not execute template HTML: "+stringErr+"\"}", http.StatusInternalServerError)
		return
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Triggering downloadHandler...")

	pattern := `^https://soundcloud\.com/`

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

		// Check if the URL is a SoundCloud URL
		match, err := regexp.MatchString(pattern, targetUrl)
		if !match || err != nil {
			stringErr := "Not a SoundCloud URL"
			fmt.Println("[x] Error matching regex: " + stringErr)
			http.Error(w, "{\"error\": \"Error matching regex: "+stringErr+"\"}", http.StatusInternalServerError)
			return
		}

		// Actually download the thing
		cmd := exec.Command("scdl", "-l", targetUrl, "-a", "--original-name", "--path", "./destination")

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {
			stringErr := fmt.Sprintf("%s", err)

			if stderr.Len() > 0 {
				stringErr = stderr.String()
			}

			fmt.Println("[x] Error downloading: " + stringErr)
			http.Error(w, "{\"error\": \""+stringErr+"\"}", http.StatusInternalServerError)
		} else {
			fmt.Printf("[i] scdl OK")
		}

		w.Header().Set("Content-Type", "application/json")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.URL.Path)
	folderPath := "./destination"

	filePath := filepath.Join(folderPath, filename)

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set the appropriate headers for the file download
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)

	// Copy the file content to the response writer
	_, err = io.Copy(w, file)
	if err != nil {
		log.Fatal(err)
	}
}
