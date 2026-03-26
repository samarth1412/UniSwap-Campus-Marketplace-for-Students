package handlers

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"uniswap-campus-marketplace/repository"
	"uniswap-campus-marketplace/services"
)

type UploadHandler struct {
	listingImageService *services.ListingImageService
}

func NewUploadHandler(listingImageService *services.ListingImageService) *UploadHandler {
	return &UploadHandler{listingImageService: listingImageService}
}

func (h *UploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if _, ok := userIDFromContext(r); !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	url, err := saveMultipartFile(file, header, "", 0)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save file")
		return
	}

	writeSuccess(w, http.StatusCreated, map[string]string{
		"url": url,
	})
}

func (h *UploadHandler) UploadListingImages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := userIDFromContext(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	listingID, ok := extractListingIDFromImagesPath(r.URL.Path)
	if !ok {
		writeError(w, http.StatusNotFound, "resource not found")
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	headers := r.MultipartForm.File["files"]
	if len(headers) == 0 {
		writeError(w, http.StatusBadRequest, "files are required")
		return
	}

	urls := make([]string, 0, len(headers))
	for idx, header := range headers {
		file, err := header.Open()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to process file")
			return
		}

		url, err := saveMultipartFile(file, header, fmt.Sprintf("listing_%d", listingID), idx)
		_ = file.Close()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to save file")
			return
		}
		urls = append(urls, url)
	}

	images, err := h.listingImageService.AddImages(r.Context(), userID, listingID, urls)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrValidation):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrForbidden):
			writeError(w, http.StatusForbidden, "forbidden")
		case errors.Is(err, repository.ErrListingNotFound):
			writeError(w, http.StatusNotFound, "listing not found")
		default:
			writeError(w, http.StatusInternalServerError, "failed to save listing images")
		}
		return
	}

	writeSuccess(w, http.StatusCreated, images)
}

func saveMultipartFile(file multipart.File, header *multipart.FileHeader, prefix string, index int) (string, error) {
	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, 0o755); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = ".bin"
	}

	base := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), index, ext)
	if prefix != "" {
		base = fmt.Sprintf("%s_%s", prefix, base)
	}

	dstPath := filepath.Join(uploadsDir, base)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return "/" + filepath.ToSlash(dstPath), nil
}

func extractListingIDFromImagesPath(path string) (int64, bool) {
	trimmed := strings.Trim(strings.TrimPrefix(path, "/api/listings/"), "/")
	parts := strings.Split(trimmed, "/")
	if len(parts) != 2 || parts[1] != "images" {
		return 0, false
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}
