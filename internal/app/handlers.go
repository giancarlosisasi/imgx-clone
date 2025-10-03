package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	_  = iota             // ignore first value (0)
	KB = 1 << (10 * iota) // 1 << 10
	MB                    // 1 << 20
	GB                    // 1 << 30
	TB                    // 1 << 40
	PB                    // 1 << 50
)

const (
	MaxUploadSize = 10 * MB
	MaxMemory     = 5 * MB
)

func (app *Application) getImageHandler(w http.ResponseWriter, r *http.Request) {
	// get the id
	filename := r.PathValue("filename")

	if filename == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}

	// sanitize filename
	filename = filepath.Base(filename) // removes any paths components

	if filename == "." || filename == ".." {
		log.Error().Msg("Filename is a directory (someone is trying to attack us?)")
		http.Error(w, "invalid filename", http.StatusBadRequest)
		return
	}

	if filename == "" || strings.HasPrefix(filename, ".") {
		log.Error().Msg("Filename is empty or its a binary")
		http.Error(w, "invalid filename", http.StatusBadRequest)
		return
	}

	if strings.Contains(filename, "\x00") {
		log.Error().Msg("Filename contains null bytes")
		http.Error(w, "invalid filename", http.StatusBadRequest)
		return
	}

	fp := filepath.Join(app.config.ImageUploadDir, filename)

	absPath, err := filepath.Abs(fp)
	if err != nil {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	absUploadDir, err := filepath.Abs(app.config.ImageUploadDir)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// make sure the file path is inside our upload dir folder
	// this help us to prevent directory-transversal attacks
	if !strings.HasPrefix(absPath, absUploadDir) {
		log.Error().Msg("Filename is not inside our upload directory")
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	fileInfo, err := os.Stat(fp)
	if os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	if fileInfo.IsDir() {
		log.Error().Msg("Filename is a directory")
		http.Error(w, "invalid file", http.StatusBadRequest)
		return
	}

	if fileInfo.Size() > 100<<20 { // 100mb limit
		http.Error(w, "file too large", http.StatusRequestEntityTooLarge)
		return
	}

	// handle unicode filenames (utf-8 encoded)
	encodedFilename := url.QueryEscape(filename)

	// attachment forces download in the browser (opens save as dialog)
	// inline display the image in the browser
	// w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", encodedFilename))
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", encodedFilename))

	// set cache header for performance
	w.Header().Set("Cache-Control", "public, max-age=86400") // 24 hours

	http.ServeFile(w, r, fp)
}

func (app *Application) uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)

	// parse multiform data using a memory limit
	if err := r.ParseMultipartForm(MaxMemory); err != nil {
		http.Error(w, "File too large or invalid form", http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll() // clean up tmp files

	// get file
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// validate that its an image (check mime type and it's actual content)
	if !isAllowedFileType(file, header.Filename) {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// generate a random uuid as identifier
	safeFilename := generateSafeFilename(header.Filename)

	if err := os.MkdirAll(app.config.ImageUploadDir, 0755); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	destPath := filepath.Join(app.config.ImageUploadDir, safeFilename)

	dst, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	if err := copyWithContext(ctx, dst, file); err != nil {
		os.Remove(destPath) // cleanup on error
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	description := r.FormValue("description")

	w.Header().Set("Content-Type", "application/json")

	app.commitHeadersAndWriteStatus(w, http.StatusCreated)

	fmt.Fprintf(w, `{"message":"File uploaded", "filename":"/api/v1/images/%s","description":"%s"}`, safeFilename, description)
}

func isAllowedFileType(file io.ReadSeeker, filename string) bool {
	// read first 512 bytes for MIME detection
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return false
	}

	// reset file pointer
	_, err = file.Seek(0, 0)
	if err != nil {
		return false
	}

	contentType := http.DetectContentType(buffer)

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/avif": true,
		"image/webp": true,
	}

	return allowedTypes[contentType]
}

func generateSafeFilename(original string) string {
	ext := filepath.Ext(original)

	ext = strings.ToLower(ext)
	if len(ext) > 10 {
		ext = ".bin"
	}

	// generate a unique name using UUID
	uniqueName := uuid.New().String()

	return uniqueName + ext
}

// copyWithContext copies data with context cancellation support
func copyWithContext(ctx context.Context, dst io.Writer, src io.Reader) error {
	// use a buffer for efficient copying
	buf := make([]byte, 32*1024)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			nr, err := src.Read(buf)
			if nr > 0 {
				nw, ew := dst.Write(buf[0:nr])
				if ew != nil {
					return ew
				}

				if nr != nw {
					return io.ErrShortWrite
				}
			}
			if err != nil {
				if err == io.EOF {
					return nil
				}

				return err
			}
		}
	}
}
