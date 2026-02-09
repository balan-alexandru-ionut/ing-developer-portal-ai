package routes

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"html"
	"net/http"
	"path/filepath"
	"strings"
)

type FileEntry struct {
	FilePath string `json:"filePath"`
	Code     string `json:"code"`
}

type Payload struct {
	Time  string      `json:"time"`
	Files []FileEntry `json:"files"`
}

// GET /api/download
func generateFilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpErrorJSON(w, http.StatusMethodNotAllowed, "use GET")
		return
	}

	// Parse JSON from memory
	var payload Payload
	if err := json.Unmarshal(GeneratedJSON, &payload); err != nil {
		httpErrorJSON(w, http.StatusBadRequest, "invalid global JSON: "+err.Error())
		return
	}

	// In‑memory ZIP buffer
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, fe := range payload.Files {

		// Sanitize and normalize paths
		cleanRel, err := sanitizeRelativePath(fe.FilePath)
		if err != nil {
			httpErrorJSON(w, http.StatusBadRequest, "invalid filePath: "+err.Error())
			return
		}

		// Ensure directories inside ZIP
		zipPath := filepath.ToSlash(cleanRel)

		// HTML decode code (because your JSON contains &lt; etc.)
		code := html.UnescapeString(fe.Code)

		// Create file inside the ZIP
		f, err := zipWriter.Create(zipPath)
		if err != nil {
			httpErrorJSON(w, http.StatusInternalServerError, "failed creating zip entry: "+err.Error())
			return
		}

		_, err = f.Write([]byte(code))
		if err != nil {
			httpErrorJSON(w, http.StatusInternalServerError, "failed writing zip entry: "+err.Error())
			return
		}
	}

	// Finalize zip
	if err := zipWriter.Close(); err != nil {
		httpErrorJSON(w, http.StatusInternalServerError, "failed finalizing zip: "+err.Error())
		return
	}

	// Send as downloadable file
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=generated_project.zip")
	w.Header().Set("Content-Length", string(len(buf.Bytes())))

	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

// --- HELPERS -------------------------------------------------------

func sanitizeRelativePath(p string) (string, error) {
	p = strings.TrimSpace(p)
	p = filepath.ToSlash(p)

	if strings.HasPrefix(p, "/") || strings.HasPrefix(p, "\\") {
		return "", errors.New("absolute paths not allowed")
	}
	p = strings.TrimPrefix(p, "./")
	p = filepath.Clean(p)

	if p == ".." || strings.HasPrefix(p, "../") || strings.Contains(p, "/../") {
		return "", errors.New("path traversal not allowed")
	}
	return filepath.FromSlash(p), nil
}

func httpErrorJSON(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
