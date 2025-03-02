package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/alandavd/asci/internal/core/ports"
	"github.com/alandavd/asci/pkg"
)

type ConvertResponse struct {
	ASCII string `json:"ascii"`
	Error string `json:"error,omitempty"`
}

func main() {
	converter := asci.NewConverter()

	http.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			sendJSONResponse(w, ConvertResponse{Error: "Failed to parse form"}, http.StatusBadRequest)
			return
		}

		// Get the file from form data
		file, _, err := r.FormFile("image")
		if err != nil {
			sendJSONResponse(w, ConvertResponse{Error: "No image file provided"}, http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Parse options
		options := []ports.ConvertOption{}

		if width := r.FormValue("width"); width != "" {
			if w, err := strconv.Atoi(width); err == nil {
				options = append(options, asci.WithWidth(w))
			}
		}

		if height := r.FormValue("height"); height != "" {
			if h, err := strconv.Atoi(height); err == nil {
				options = append(options, asci.WithHeight(h))
			}
		}

		if charset := r.FormValue("charset"); charset != "" {
			options = append(options, asci.WithCharset(charset))
		}

		if colored := r.FormValue("colored"); colored == "true" {
			options = append(options, asci.WithColor(true))
		}

		if inverted := r.FormValue("inverted"); inverted == "true" {
			options = append(options, asci.WithInverted(true))
		}

		// Convert image
		result, err := converter.Convert(file, options...)
		if err != nil {
			sendJSONResponse(w, ConvertResponse{Error: fmt.Sprintf("Conversion failed: %v", err)}, http.StatusInternalServerError)
			return
		}

		sendJSONResponse(w, ConvertResponse{ASCII: result}, http.StatusOK)
	})

	// Serve a simple HTML form for testing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>ASCII Art Converter</title>
			<style>
				body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
				form { display: flex; flex-direction: column; gap: 10px; }
				pre { white-space: pre; font-family: monospace; background: #f5f5f5; padding: 10px; }
				.result { margin-top: 20px; }
			</style>
		</head>
		<body>
			<h1>ASCII Art Converter</h1>
			<form id="convertForm">
				<div>
					<label for="image">Image:</label>
					<input type="file" id="image" name="image" accept="image/*" required>
				</div>
				<div>
					<label for="width">Width:</label>
					<input type="number" id="width" name="width" value="80">
				</div>
				<div>
					<label for="height">Height:</label>
					<input type="number" id="height" name="height" value="0">
				</div>
				<div>
					<label for="charset">Charset:</label>
					<input type="text" id="charset" name="charset" value=" .:-=+*#%@">
				</div>
				<div>
					<label for="colored">Colored:</label>
					<input type="checkbox" id="colored" name="colored">
				</div>
				<div>
					<label for="inverted">Inverted:</label>
					<input type="checkbox" id="inverted" name="inverted">
				</div>
				<button type="submit">Convert</button>
			</form>
			<div class="result">
				<pre id="output"></pre>
			</div>
			<script>
				document.getElementById('convertForm').onsubmit = async (e) => {
					e.preventDefault();
					const formData = new FormData(e.target);
					try {
						const response = await fetch('/convert', {
							method: 'POST',
							body: formData
						});
						const data = await response.json();
						if (data.error) {
							document.getElementById('output').textContent = 'Error: ' + data.error;
						} else {
							document.getElementById('output').textContent = data.ascii;
						}
					} catch (err) {
						document.getElementById('output').textContent = 'Error: ' + err.message;
					}
				};
			</script>
		</body>
		</html>
		`
		io.WriteString(w, html)
	})

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func sendJSONResponse(w http.ResponseWriter, response ConvertResponse, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
