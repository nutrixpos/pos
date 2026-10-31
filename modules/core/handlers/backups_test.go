package handlers

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/internal/testutil"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/nutrixpos/pos/modules/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func seedBackupHandlerData(t *testing.T, env *testutil.TestEnv, uploads string) {
	t.Helper()

	seed(t, env, "materials", models.Material{
		Id:      "mat-1",
		Name:    "Flour",
		Entries: []models.MaterialEntry{{Id: "e-1", Quantity: 2}},
	})
	seed(t, env, "recipes", models.Product{
		Id:        "p-1",
		Name:      "Pizza",
		ImageURL:  "pizza.png",
		Materials: []models.Material{{Id: "mat-1", Quantity: 1}},
	})
	seed(t, env, "categories", models.Category{Id: "c-1", Name: "Mains", Products: []models.Product{{Id: "p-1"}}})
	seed(t, env, "sales", models.SalesPerDay{Date: "2026-01-01", Orders: []models.SalesPerDayOrder{{Id: "o-1"}}})
	seed(t, env, "orders", models.Order{Id: "o-1", DisplayId: "A-1"})

	require.NoError(t, os.WriteFile(filepath.Join(uploads, "pizza.png"), []byte("image"), 0o644))
}

func buildBackupHandlerZip(t *testing.T, env *testutil.TestEnv) string {
	t.Helper()

	svc := services.BackupService{Logger: env.Logger, Config: env.Config}
	pkg, err := svc.BuildExport(context.Background())
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "backups.zip")
	require.NoError(t, svc.WriteBackupZIP(context.Background(), &pkg, path))
	return path
}

func multipartBackupRequest(t *testing.T, target, zipPath string, fields map[string]string) *http.Request {
	t.Helper()

	data, err := os.ReadFile(zipPath)
	require.NoError(t, err)

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	part, err := writer.CreateFormFile("file", "backups.zip")
	require.NoError(t, err)
	_, err = part.Write(data)
	require.NoError(t, err)

	for key, value := range fields {
		require.NoError(t, writer.WriteField(key, value))
	}
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, target, &buffer)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestBackupExportHandler(t *testing.T) {
	env := newHandlerEnv(t)
	env.Config.UploadsPath = t.TempDir()
	env.Config.BackupsPath = t.TempDir()

	seed(t, env, "materials", models.Material{Id: "mat-1", Name: "Flour", Entries: []models.MaterialEntry{{Id: "e-1", Quantity: 2}}})
	seed(t, env, "recipes", models.Product{Id: "p-1", Name: "Pizza"})
	seed(t, env, "categories", models.Category{Id: "c-1", Name: "Mains"})

	router := mux.NewRouter()
	router.HandleFunc("/api/backups/export", ExportBackup(env.Config, env.Logger)).Methods(http.MethodGet)

	rec := doRequest(t, router, http.MethodGet, "/api/backups/export", "")
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/zip", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Header().Get("Content-Disposition"), "nutrix-backup-")
	assert.True(t, bytes.HasPrefix(rec.Body.Bytes(), []byte("PK")))
}

func TestBackupPreviewAndImportHandlers(t *testing.T) {
	env := newHandlerEnv(t)
	env.Config.UploadsPath = t.TempDir()
	env.Config.BackupsPath = t.TempDir()

	seedBackupHandlerData(t, env, env.Config.UploadsPath)
	zipPath := buildBackupHandlerZip(t, env)

	router := mux.NewRouter()
	router.HandleFunc("/api/backups/import/preview", PreviewBackupImport(env.Config, env.Logger)).Methods(http.MethodPost)
	router.HandleFunc("/api/backups/import", ImportBackup(env.Config, env.Logger)).Methods(http.MethodPost)
	router.HandleFunc("/api/backups/files/{filename}", DownloadBackupFile(env.Config, env.Logger)).Methods(http.MethodGet)

	previewReq := multipartBackupRequest(t, "/api/backups/import/preview", zipPath, nil)
	previewRec := serve(t, router, previewReq)
	require.Equal(t, http.StatusOK, previewRec.Code)

	var preview services.BackupPreview
	decodeData(t, previewRec, &preview)
	assert.Equal(t, 1, preview.Sections.Materials)
	assert.Equal(t, 1, preview.Sections.Products)
	assert.Equal(t, 1, preview.Sections.Categories)
	assert.Equal(t, 1, preview.Sections.Sales)
	assert.Equal(t, 1, preview.Sections.Orders)
	assert.True(t, preview.HasStock)
	assert.True(t, preview.HasImages)

	importReq := multipartBackupRequest(t, "/api/backups/import", zipPath, map[string]string{
		"sections":       "materials,products,categories,sales",
		"include_stock":  "true",
		"include_images": "true",
	})
	importRec := serve(t, router, importReq)
	require.Equal(t, http.StatusOK, importRec.Code)

	var result services.BackupImportResult
	decodeData(t, importRec, &result)
	assert.Equal(t, 1, result.Imported.Materials)
	assert.Equal(t, 1, result.Imported.Products)
	assert.Equal(t, 1, result.Imported.Categories)
	assert.Equal(t, 1, result.Imported.Sales)
	assert.Equal(t, 1, result.Imported.Orders)
	require.NotEmpty(t, result.BackupFile)

	downloadRec := doRequest(t, router, http.MethodGet, "/api/backups/files/"+result.BackupFile, "")
	require.Equal(t, http.StatusOK, downloadRec.Code)
	assert.True(t, bytes.HasPrefix(downloadRec.Body.Bytes(), []byte("PK")))

	badNameRec := doRequest(t, router, http.MethodGet, "/api/backups/files/not-a-zip.txt", "")
	assert.Equal(t, http.StatusBadRequest, badNameRec.Code)

	missingRec := doRequest(t, router, http.MethodGet, "/api/backups/files/missing.zip", "")
	assert.Equal(t, http.StatusNotFound, missingRec.Code)
}

func TestBackupHandlerErrorPaths(t *testing.T) {
	env := newHandlerEnv(t)
	env.Config.UploadsPath = t.TempDir()
	env.Config.BackupsPath = t.TempDir()

	router := mux.NewRouter()
	router.HandleFunc("/api/backups/import/preview", PreviewBackupImport(env.Config, env.Logger)).Methods(http.MethodPost)
	router.HandleFunc("/api/backups/import", ImportBackup(env.Config, env.Logger)).Methods(http.MethodPost)

	rec := doRequest(t, router, http.MethodPost, "/api/backups/import/preview", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	badZip := filepath.Join(t.TempDir(), "bad.zip")
	require.NoError(t, os.WriteFile(badZip, []byte("not a zip"), 0o644))

	req := multipartBackupRequest(t, "/api/backups/import/preview", badZip, nil)
	rec = serve(t, router, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	seed(t, env, "materials", models.Material{Id: "mat-keep", Name: "Keep"})
	zipPath := buildBackupHandlerZip(t, env)

	req = multipartBackupRequest(t, "/api/backups/import", zipPath, map[string]string{
		"sections": "products",
	})
	rec = serve(t, router, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Validation failed before the backup/wipe, so existing data is untouched.
	var material models.Material
	require.NoError(t, env.Client.Database(env.Config.Databases[0].Database).Collection("materials").
		FindOne(context.Background(), map[string]interface{}{"id": "mat-keep"}).Decode(&material))
}
