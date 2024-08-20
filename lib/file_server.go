package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func fileServer(root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Clean the requested path to prevent directory traversal attacks
		filePath := filepath.Join(root, filepath.Clean(r.URL.Path))

		// Check if the path is a directory
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		if fileInfo.IsDir() {
			// Try to serve index files in order if they exist
			indexFiles := []string{"index.html", "index.htm", "index.php"}
			for _, indexFile := range indexFiles {
				indexFilePath := filepath.Join(filePath, indexFile)
				if _, err := os.Stat(indexFilePath); err == nil {
					serveFile(w, r, indexFilePath)
					return
				}
			}

			// Generate and serve directory listing if no index file is found
			serveDirectoryListing(w, filePath)
			return
		}

		// Serve the file if it's not a directory
		serveFile(w, r, filePath)
	}
}

// Get the content type based on the file extension
func getContentType(ext string) string {
	switch ext {
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".php":
		return "application/x-httpd-php"
	case ".json":
		return "application/json"
	case ".xml":
		return "application/xml"
	case ".svg":
		return "image/svg+xml"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".pdf":
		return "application/pdf"
	case ".webm":
		return "video/webm"
	case ".mp3":
		return "audio/mpeg"
	case ".webp":
		return "image/webp"
	case ".mp4":
		return "video/mp4"
	default:
		return "application/octet-stream"
	}
}

// Serve a file with appropriate Content-Type header
func serveFile(w http.ResponseWriter, r *http.Request, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Get the file info
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Error retrieving file info", http.StatusInternalServerError)
		return
	}

	// Detect the MIME type and set the Content-Type header
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	contentType := getContentType(filepath.Ext(fileInfo.Name()))
	w.Header().Set("Content-Type", contentType)

	// Reset the file read pointer to the beginning of the file
	file.Seek(0, 0)

	// Serve the content
	if filepath.Ext(fileInfo.Name()) == ".html" || filepath.Ext(fileInfo.Name()) == ".htm" {
		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}

		// Replace the last instance of </body> with the script
		contentStr := string(content)
		lastIndex := strings.LastIndex(contentStr, "</body>")
		if lastIndex != -1 {
			contentStr = contentStr[:lastIndex] + reloadScript() + contentStr[lastIndex:]
		}
		fmt.Fprint(w, contentStr)
	} else {
		http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
	}
}

// Generate and serve an HTML directory listing
func serveDirectoryListing(w http.ResponseWriter, dirPath string) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		http.Error(w, "Failed to read directory", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><body><h1>Directory listing for %s</h1><ul>", dirPath)

	// List the files and directories
	for _, file := range files {
		name := file.Name()
		if file.IsDir() {
			name += "/"
		}
		link := filepath.Join("/", name)
		fmt.Fprintf(w, "<li><a href=\"%s\">%s</a></li>", link, name)
	}

	fmt.Fprintln(w, "</ul>"+reloadScript()+"</body></html>")
}
