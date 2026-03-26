package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"uniswap-campus-marketplace/middleware"
	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/services"
)

type listingRepoForUploadTest struct {
	getByIDFn  func(ctx context.Context, id int64) (*models.Listing, error)
	createFn   func(ctx context.Context, listing *models.Listing) (*models.Listing, error)
	getAllFn   func(ctx context.Context, search string) ([]models.Listing, error)
	updateByID func(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error)
	deleteByID func(ctx context.Context, id int64) (int64, error)
}

func (r *listingRepoForUploadTest) Create(ctx context.Context, listing *models.Listing) (*models.Listing, error) {
	return r.createFn(ctx, listing)
}
func (r *listingRepoForUploadTest) GetAll(ctx context.Context, search string) ([]models.Listing, error) {
	return r.getAllFn(ctx, search)
}
func (r *listingRepoForUploadTest) GetByID(ctx context.Context, id int64) (*models.Listing, error) {
	return r.getByIDFn(ctx, id)
}
func (r *listingRepoForUploadTest) UpdateByID(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error) {
	return r.updateByID(ctx, id, listing)
}
func (r *listingRepoForUploadTest) DeleteByID(ctx context.Context, id int64) (int64, error) {
	return r.deleteByID(ctx, id)
}

type imageRepoForUploadTest struct {
	createManyFn func(ctx context.Context, listingID int64, imageURLs []string) ([]models.ListingImage, error)
}

func (r *imageRepoForUploadTest) CreateMany(ctx context.Context, listingID int64, imageURLs []string) ([]models.ListingImage, error) {
	return r.createManyFn(ctx, listingID, imageURLs)
}

type parserForUploadTest struct{}

func (p *parserForUploadTest) ParseToken(tokenString string) (int64, error) {
	return 1, nil
}

func TestExtractListingIDFromImagesPath(t *testing.T) {
	id, ok := extractListingIDFromImagesPath("/api/listings/12/images")
	if !ok || id != 12 {
		t.Fatalf("expected id parse success")
	}

	if _, ok := extractListingIDFromImagesPath("/api/listings/abc/images"); ok {
		t.Fatalf("expected invalid id parse failure")
	}
}

func TestUploadListingImagesUnauthorized(t *testing.T) {
	h := NewUploadHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/listings/1/images", nil)
	rec := httptest.NewRecorder()
	h.UploadListingImages(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestUploadListingImagesMissingFiles(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.Close()

	listingSvc := services.NewListingImageService(
		&listingRepoForUploadTest{
			getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
				return &models.Listing{ID: id, UserID: 1}, nil
			},
		},
		&imageRepoForUploadTest{},
	)

	h := NewUploadHandler(listingSvc)
	req := httptest.NewRequest(http.MethodPost, "/api/listings/1/images", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer tkn")
	rec := httptest.NewRecorder()

	middleware.Auth(&parserForUploadTest{})(http.HandlerFunc(h.UploadListingImages)).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestUploadListingImagesSuccess(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("chdir temp dir: %v", err)
	}
	defer func() { _ = os.Chdir(wd) }()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	file1, err := writer.CreateFormFile("files", "img1.jpg")
	if err != nil {
		t.Fatalf("create form file1: %v", err)
	}
	_, _ = file1.Write([]byte("img-one"))

	file2, err := writer.CreateFormFile("files", "img2.jpg")
	if err != nil {
		t.Fatalf("create form file2: %v", err)
	}
	_, _ = file2.Write([]byte("img-two"))
	_ = writer.Close()

	imageRepo := &imageRepoForUploadTest{
		createManyFn: func(ctx context.Context, listingID int64, imageURLs []string) ([]models.ListingImage, error) {
			out := make([]models.ListingImage, 0, len(imageURLs))
			for i, u := range imageURLs {
				out = append(out, models.ListingImage{
					ID:        int64(i + 1),
					ListingID: listingID,
					ImageURL:  u,
				})
			}
			return out, nil
		},
	}

	listingSvc := services.NewListingImageService(
		&listingRepoForUploadTest{
			getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
				return &models.Listing{ID: id, UserID: 1}, nil
			},
		},
		imageRepo,
	)

	h := NewUploadHandler(listingSvc)
	req := httptest.NewRequest(http.MethodPost, "/api/listings/5/images", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	middleware.Auth(&parserForUploadTest{})(http.HandlerFunc(h.UploadListingImages)).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var payload struct {
		Success bool                  `json:"success"`
		Data    []models.ListingImage `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if !payload.Success || len(payload.Data) != 2 {
		t.Fatalf("unexpected upload response: %+v", payload)
	}

	if _, err := os.Stat(filepath.Join("uploads")); err != nil {
		t.Fatalf("expected uploads folder to be created: %v", err)
	}
}
