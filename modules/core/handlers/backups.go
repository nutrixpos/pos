package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/services"
)

// ExportBackup returns a ZIP archive containing the whole backup, stock
// entries, product images, sales, and orders.
func ExportBackup(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		backupSvc := services.BackupService{Logger: logger, Config: config}

		pkg, err := backupSvc.BuildExport(r.Context())
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmp, err := os.CreateTemp("", "nutrix-backup-export-*.zip")
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpPath := tmp.Name()
		_ = tmp.Close()
		defer func() { _ = os.Remove(tmpPath) }()

		if err := backupSvc.WriteBackupZIP(r.Context(), &pkg, tmpPath); err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filename := fmt.Sprintf("nutrix-backup-%s.zip", time.Now().UTC().Format("20060102-150405"))
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
		http.ServeFile(w, r, tmpPath)
	}
}

// PreviewBackupImport validates an uploaded archive and returns counts without
// changing the database.
func PreviewBackupImport(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path, cleanup, err := saveUploadedBackup(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer cleanup()

		backupSvc := services.BackupService{Logger: logger, Config: config}
		preview, err := backupSvc.PreviewZip(r.Context(), path)
		if err != nil {
			logger.Error(err.Error())
			status := http.StatusInternalServerError
			if errors.Is(err, services.ErrInvalidBackup) {
				status = http.StatusBadRequest
			}
			http.Error(w, err.Error(), status)
			return
		}

		writeBackupJSON(w, http.StatusOK, JSONApiOkResponse{Data: preview})
	}
}

// ImportBackup replaces the backup data with an uploaded archive. A full backup is
// created before any data is deleted.
func ImportBackup(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path, cleanup, err := saveUploadedBackup(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer cleanup()

		opts := services.BackupImportOptions{
			Sections:      parseBackupSections(r.FormValue("sections")),
			IncludeStock:  parseBackupBool(r.FormValue("include_stock")),
			IncludeImages: parseBackupBool(r.FormValue("include_images")),
		}

		backupSvc := services.BackupService{Logger: logger, Config: config}
		result, err := backupSvc.ImportZip(r.Context(), path, opts)
		if err != nil {
			logger.Error(err.Error())
			status := http.StatusInternalServerError
			if errors.Is(err, services.ErrInvalidBackup) {
				status = http.StatusBadRequest
			}

			writeBackupJSON(w, status, struct {
				Error      string `json:"error"`
				BackupFile string `json:"backup_file,omitempty"`
			}{
				Error:      err.Error(),
				BackupFile: result.BackupFile,
			})
			return
		}

		writeBackupJSON(w, http.StatusOK, JSONApiOkResponse{Data: result})
	}
}

// DownloadBackupFile serves a backup file created before an import.
func DownloadBackupFile(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := mux.Vars(r)["filename"]
		if filename == "" || filename != filepath.Base(filename) || !strings.HasSuffix(strings.ToLower(filename), ".zip") {
			http.Error(w, "invalid backup file name", http.StatusBadRequest)
			return
		}

		backupSvc := services.BackupService{Logger: logger, Config: config}
		path := filepath.Join(backupSvc.BackupsDir(), filename)

		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			http.Error(w, "backup not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
		http.ServeFile(w, r, path)
	}
}

func saveUploadedBackup(w http.ResponseWriter, r *http.Request) (string, func(), error) {
	r.Body = http.MaxBytesReader(w, r.Body, services.MaxBackupUploadBytes)

	if err := r.ParseMultipartForm(8 << 20); err != nil {
		return "", nil, err
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return "", nil, err
	}
	defer func() { _ = file.Close() }()

	tmp, err := os.CreateTemp("", "nutrix-backup-upload-*.zip")
	if err != nil {
		return "", nil, err
	}

	path := tmp.Name()
	cleanup := func() { _ = os.Remove(path) }

	if _, err := io.Copy(tmp, file); err != nil {
		_ = tmp.Close()
		cleanup()
		return "", nil, err
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return "", nil, err
	}

	return path, cleanup, nil
}

func parseBackupSections(raw string) map[string]bool {
	sections := make(map[string]bool)
	for _, section := range strings.Split(raw, ",") {
		section = strings.ToLower(strings.TrimSpace(section))
		if section == "" {
			continue
		}
		sections[section] = true
	}
	return sections
}

func parseBackupBool(raw string) bool {
	value, err := strconv.ParseBool(raw)
	return err == nil && value
}

func writeBackupJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
