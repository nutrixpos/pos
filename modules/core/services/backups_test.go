package services

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/internal/testutil"
	"github.com/nutrixpos/pos/modules/core/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func newBackupTestService(t *testing.T) (*testutil.TestEnv, *BackupService) {
	t.Helper()

	env := newServiceEnv(t)
	env.Config.UploadsPath = t.TempDir()
	env.Config.BackupsPath = t.TempDir()

	return env, &BackupService{Logger: env.Logger, Config: env.Config}
}

func seedBackupTestData(t *testing.T, env *testutil.TestEnv) {
	t.Helper()

	seedServiceDoc(t, env, "materials", models.Material{
		Id:   "mat-1",
		Name: "Flour",
		Unit: "kg",
		Entries: []models.MaterialEntry{{
			Id:               "entry-1",
			Quantity:         5,
			PurchaseQuantity: 10,
			PurchasePrice:    100,
		}},
	})
	seedServiceDoc(t, env, "materials", models.Material{Id: "mat-2", Name: "Cheese", Unit: "kg"})

	seedServiceDoc(t, env, "recipes", models.Product{
		Id:          "p-1",
		Name:        "Pizza",
		Price:       25,
		ImageURL:    "pizza.png",
		Materials:   []models.Material{{Id: "mat-1", Name: "Flour", Quantity: 2}},
		SubProducts: []models.Product{{Id: "p-2", Quantity: 1}},
	})
	seedServiceDoc(t, env, "recipes", models.Product{Id: "p-2", Name: "Dough", Price: 5})

	seedServiceDoc(t, env, "categories", models.Category{
		Id:       "cat-1",
		Name:     "Mains",
		Products: []models.Product{{Id: "p-1"}},
	})

	seedServiceDoc(t, env, "sales", models.SalesPerDay{
		Date:       "2026-01-01",
		Costs:      10,
		TotalSales: 25,
		Orders: []models.SalesPerDayOrder{{
			Id:    "order-1",
			Order: models.Order{Id: "order-1", DisplayId: "A-5", Cost: 10, SalePrice: 25},
		}},
	})
	seedServiceDoc(t, env, "orders", models.Order{
		Id:          "order-1",
		DisplayId:   "A-5",
		SubmittedAt: time.Now(),
		State:       "finished",
	})

	require.NoError(t, os.WriteFile(filepath.Join(env.Config.UploadsPath, "pizza.png"), []byte("image-bytes"), 0o644))
}

func writeBackupTestZip(t *testing.T, path string, pkg *BackupPackage, images map[string][]byte) {
	t.Helper()

	file, err := os.Create(path)
	require.NoError(t, err)
	defer func() { _ = file.Close() }()

	writer := zip.NewWriter(file)

	for name, data := range images {
		entry, err := writer.Create(name)
		require.NoError(t, err)
		_, err = entry.Write(data)
		require.NoError(t, err)
	}

	if pkg != nil {
		data, err := json.Marshal(pkg)
		require.NoError(t, err)
		entry, err := writer.Create("backup.json")
		require.NoError(t, err)
		_, err = entry.Write(data)
		require.NoError(t, err)
	}

	require.NoError(t, writer.Close())
}

func writeBackupTestRawZip(t *testing.T, path string, entries map[string][]byte) {
	t.Helper()

	file, err := os.Create(path)
	require.NoError(t, err)
	defer func() { _ = file.Close() }()

	writer := zip.NewWriter(file)
	for name, data := range entries {
		entry, err := writer.Create(name)
		require.NoError(t, err)
		_, err = entry.Write(data)
		require.NoError(t, err)
	}
	require.NoError(t, writer.Close())
}

func countServiceDocs(t *testing.T, env *testutil.TestEnv, collection string, filter bson.M) int64 {
	t.Helper()

	count, err := env.Client.Database(env.Config.Databases[0].Database).Collection(collection).CountDocuments(context.Background(), filter)
	require.NoError(t, err)
	return count
}

func TestBackupExportImportRoundTrip(t *testing.T) {
	env, svc := newBackupTestService(t)
	seedBackupTestData(t, env)

	pkg, err := svc.BuildExport(context.Background())
	require.NoError(t, err)
	require.Len(t, pkg.Materials, 2)
	require.Len(t, pkg.Products, 2)
	require.Len(t, pkg.Categories, 1)
	require.Len(t, pkg.Sales, 1)
	require.Len(t, pkg.Orders, 1)

	var flourQuantity float64
	for _, material := range pkg.Materials {
		if material.Id == "mat-1" {
			flourQuantity = material.Quantity
		}
	}
	assert.Equal(t, 5.0, flourQuantity)

	exportPath := filepath.Join(t.TempDir(), "export.zip")
	require.NoError(t, svc.WriteBackupZIP(context.Background(), &pkg, exportPath))

	preview, err := svc.PreviewZip(context.Background(), exportPath)
	require.NoError(t, err)
	assert.Equal(t, 2, preview.Sections.Materials)
	assert.Equal(t, 1, preview.Sections.MaterialEntries)
	assert.Equal(t, 2, preview.Sections.Products)
	assert.Equal(t, 1, preview.Sections.Categories)
	assert.Equal(t, 1, preview.Sections.Sales)
	assert.Equal(t, 1, preview.Sections.Orders)
	assert.Equal(t, 1, preview.Sections.Images)
	assert.True(t, preview.HasStock)
	assert.True(t, preview.HasImages)
	assert.Equal(t, 2, preview.Current.Materials)

	// Add data that must be removed by the replace import.
	seedServiceDoc(t, env, "materials", models.Material{Id: "mat-extra", Name: "Extra"})
	seedServiceDoc(t, env, "recipes", models.Product{Id: "p-extra", Name: "Extra product"})
	seedServiceDoc(t, env, "categories", models.Category{Id: "cat-extra", Name: "Extra category"})
	seedServiceDoc(t, env, "sales", models.SalesPerDay{Date: "2026-01-02"})
	seedServiceDoc(t, env, "orders", models.Order{Id: "order-extra", DisplayId: "A-9"})

	// Existing files that are not referenced by imported products must survive.
	keepPath := filepath.Join(env.Config.UploadsPath, "keep.png")
	require.NoError(t, os.WriteFile(keepPath, []byte("keep"), 0o644))

	result, err := svc.ImportZip(context.Background(), exportPath, BackupImportOptions{
		Sections: map[string]bool{
			backupSectionMaterials:  true,
			backupSectionProducts:   true,
			backupSectionCategories: true,
			backupSectionSales:      true,
		},
		IncludeStock:  true,
		IncludeImages: true,
	})
	require.NoError(t, err)
	assert.Equal(t, 2, result.Imported.Materials)
	assert.Equal(t, 1, result.Imported.MaterialEntries)
	assert.Equal(t, 2, result.Imported.Products)
	assert.Equal(t, 1, result.Imported.Categories)
	assert.Equal(t, 1, result.Imported.Sales)
	assert.Equal(t, 1, result.Imported.Orders)
	assert.Equal(t, 1, result.Imported.Images)
	require.NotEmpty(t, result.BackupFile)

	backupPath := filepath.Join(svc.BackupsDir(), result.BackupFile)
	_, err = os.Stat(backupPath)
	require.NoError(t, err)

	backupPreview, err := svc.PreviewZip(context.Background(), backupPath)
	require.NoError(t, err)
	assert.Equal(t, 3, backupPreview.Sections.Materials)
	assert.Equal(t, 2, backupPreview.Sections.Categories)
	assert.Equal(t, 2, backupPreview.Sections.Sales)
	assert.Equal(t, 2, backupPreview.Sections.Orders)

	assert.Equal(t, int64(2), countServiceDocs(t, env, "materials", bson.M{}))
	assert.Equal(t, int64(2), countServiceDocs(t, env, "recipes", bson.M{}))
	assert.Equal(t, int64(1), countServiceDocs(t, env, "categories", bson.M{}))
	assert.Equal(t, int64(1), countServiceDocs(t, env, "sales", bson.M{}))
	assert.Equal(t, int64(1), countServiceDocs(t, env, "orders", bson.M{}))

	var material models.Material
	require.NoError(t, env.Client.Database(env.Config.Databases[0].Database).Collection("materials").
		FindOne(context.Background(), bson.M{"id": "mat-1"}).Decode(&material))
	require.Len(t, material.Entries, 1)
	assert.Equal(t, 5.0, material.Entries[0].Quantity)

	imageData, err := os.ReadFile(filepath.Join(env.Config.UploadsPath, "pizza.png"))
	require.NoError(t, err)
	assert.Equal(t, "image-bytes", string(imageData))

	_, err = os.Stat(keepPath)
	require.NoError(t, err)

	var settings models.Settings
	require.NoError(t, env.Client.Database(env.Config.Databases[0].Database).Collection("settings").
		FindOne(context.Background(), bson.M{}).Decode(&settings))
	require.Len(t, settings.Orders.Queues, 1)
	assert.Equal(t, uint32(6), settings.Orders.Queues[0].Next)
}

func TestBackupImportWithoutStockOrImages(t *testing.T) {
	env, svc := newBackupTestService(t)
	seedBackupTestData(t, env)

	pkg, err := svc.BuildExport(context.Background())
	require.NoError(t, err)

	exportPath := filepath.Join(t.TempDir(), "export.zip")
	require.NoError(t, svc.WriteBackupZIP(context.Background(), &pkg, exportPath))

	require.NoError(t, os.Remove(filepath.Join(env.Config.UploadsPath, "pizza.png")))

	result, err := svc.ImportZip(context.Background(), exportPath, BackupImportOptions{
		Sections: map[string]bool{backupSectionMaterials: true},
	})
	require.NoError(t, err)
	assert.Equal(t, 2, result.Imported.Materials)
	assert.Equal(t, 0, result.Imported.MaterialEntries)
	assert.Equal(t, 0, result.Imported.Images)

	var material models.Material
	require.NoError(t, env.Client.Database(env.Config.Databases[0].Database).Collection("materials").
		FindOne(context.Background(), bson.M{"id": "mat-1"}).Decode(&material))
	assert.Empty(t, material.Entries)
	assert.Equal(t, 0.0, material.Quantity)

	assert.Equal(t, int64(0), countServiceDocs(t, env, "recipes", bson.M{}))
	assert.Equal(t, int64(0), countServiceDocs(t, env, "categories", bson.M{}))

	// Sales & Orders are outside the backup and are not selected, so they stay.
	assert.Equal(t, int64(1), countServiceDocs(t, env, "sales", bson.M{}))
	assert.Equal(t, int64(1), countServiceDocs(t, env, "orders", bson.M{}))

	_, err = os.Stat(filepath.Join(env.Config.UploadsPath, "pizza.png"))
	assert.True(t, os.IsNotExist(err))
}

func TestBackupSalesOnlyImportKeepsCatalog(t *testing.T) {
	env, svc := newBackupTestService(t)
	seedBackupTestData(t, env)

	pkg, err := svc.BuildExport(context.Background())
	require.NoError(t, err)

	exportPath := filepath.Join(t.TempDir(), "export.zip")
	require.NoError(t, svc.WriteBackupZIP(context.Background(), &pkg, exportPath))

	result, err := svc.ImportZip(context.Background(), exportPath, BackupImportOptions{
		Sections: map[string]bool{backupSectionSales: true},
	})
	require.NoError(t, err)
	assert.Equal(t, 0, result.Imported.Materials)
	assert.Equal(t, 0, result.Imported.Products)
	assert.Equal(t, 0, result.Imported.Categories)
	assert.Equal(t, 1, result.Imported.Sales)
	assert.Equal(t, 1, result.Imported.Orders)

	// Only Sales & Orders were selected, so the catalog must stay untouched.
	assert.Equal(t, int64(2), countServiceDocs(t, env, "materials", bson.M{}))
	assert.Equal(t, int64(2), countServiceDocs(t, env, "recipes", bson.M{}))
	assert.Equal(t, int64(1), countServiceDocs(t, env, "categories", bson.M{}))
	assert.Equal(t, int64(1), countServiceDocs(t, env, "sales", bson.M{}))
	assert.Equal(t, int64(1), countServiceDocs(t, env, "orders", bson.M{}))
}

func TestBackupImportMissingImageWarning(t *testing.T) {
	env, svc := newBackupTestService(t)
	seedBackupTestData(t, env)

	// Export without the image so the archive cannot restore it.
	require.NoError(t, os.Remove(filepath.Join(env.Config.UploadsPath, "pizza.png")))

	pkg, err := svc.BuildExport(context.Background())
	require.NoError(t, err)

	exportPath := filepath.Join(t.TempDir(), "export.zip")
	require.NoError(t, svc.WriteBackupZIP(context.Background(), &pkg, exportPath))

	result, err := svc.ImportZip(context.Background(), exportPath, BackupImportOptions{
		Sections: map[string]bool{
			backupSectionMaterials: true,
			backupSectionProducts:  true,
		},
		IncludeImages: true,
	})
	require.NoError(t, err)
	assert.Equal(t, 0, result.Imported.Images)
	assert.NotEmpty(t, result.Warnings)

	found := false
	for _, warning := range result.Warnings {
		if strings.Contains(warning, "missing image") {
			found = true
		}
	}
	assert.True(t, found)
}

func TestBackupImportDependencies(t *testing.T) {
	pkg := &BackupPackage{
		Format:     BackupFormat,
		Version:    BackupVersion,
		Materials:  []models.Material{{Id: "m"}},
		Products:   []models.Product{{Id: "p"}},
		Categories: []models.Category{{Id: "c"}},
		Sales:      []models.SalesPerDay{},
		Orders:     []models.Order{},
	}

	tests := []struct {
		name string
		opts BackupImportOptions
	}{
		{"no sections", BackupImportOptions{}},
		{"unknown section", BackupImportOptions{Sections: map[string]bool{"other": true}}},
		{"products without materials", BackupImportOptions{Sections: map[string]bool{backupSectionProducts: true}}},
		{"categories without products", BackupImportOptions{Sections: map[string]bool{backupSectionCategories: true}}},
		{"stock without materials", BackupImportOptions{Sections: map[string]bool{backupSectionProducts: true, backupSectionMaterials: false}, IncludeStock: true}},
		{"images without products", BackupImportOptions{Sections: map[string]bool{backupSectionMaterials: true}, IncludeImages: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBackupImportSelection(tt.opts, pkg)
			require.Error(t, err)
			assert.True(t, errors.Is(err, ErrInvalidBackup))
		})
	}

	err := validateBackupImportSelection(BackupImportOptions{
		Sections: map[string]bool{
			backupSectionMaterials:  true,
			backupSectionProducts:   true,
			backupSectionCategories: true,
			backupSectionSales:      true,
		},
		IncludeStock:  true,
		IncludeImages: true,
	}, pkg)
	require.NoError(t, err)

	err = validateBackupImportSelection(BackupImportOptions{
		Sections: map[string]bool{backupSectionMaterials: true},
	}, &BackupPackage{Format: BackupFormat, Version: BackupVersion})
	require.Error(t, err)
}

func TestBackupReadZipValidation(t *testing.T) {
	_, svc := newBackupTestService(t)

	encode := func(pkg *BackupPackage) []byte {
		data, err := json.Marshal(pkg)
		require.NoError(t, err)
		return data
	}

	tests := []struct {
		name    string
		entries map[string][]byte
	}{
		{"missing backup file", map[string][]byte{"images/pizza.png": []byte("x")}},
		{"invalid json", map[string][]byte{"backup.json": []byte("{")}},
		{"wrong format", map[string][]byte{"backup.json": encode(&BackupPackage{Format: "other", Version: BackupVersion})}},
		{"wrong version", map[string][]byte{"backup.json": encode(&BackupPackage{Format: BackupFormat, Version: 99})}},
		{
			"duplicate material id",
			map[string][]byte{"backup.json": encode(&BackupPackage{
				Format: BackupFormat, Version: BackupVersion,
				Materials: []models.Material{{Id: "same"}, {Id: "same"}},
			})},
		},
		{
			"duplicate order id",
			map[string][]byte{"backup.json": encode(&BackupPackage{
				Format: BackupFormat, Version: BackupVersion,
				Orders: []models.Order{{Id: "same"}, {Id: "same"}},
			})},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "backups.zip")
			writeBackupTestRawZip(t, path, tt.entries)
			_, err := svc.PreviewZip(context.Background(), path)
			require.Error(t, err)
			assert.True(t, errors.Is(err, ErrInvalidBackup))
		})
	}
}

func TestBackupReadZipRejectsDuplicateImage(t *testing.T) {
	_, svc := newBackupTestService(t)

	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	pkg := &BackupPackage{Format: BackupFormat, Version: BackupVersion, Materials: []models.Material{}, Products: []models.Product{}, Categories: []models.Category{}, Sales: []models.SalesPerDay{}, Orders: []models.Order{}}
	data, err := json.Marshal(pkg)
	require.NoError(t, err)

	entry, err := writer.Create("backup.json")
	require.NoError(t, err)
	_, err = entry.Write(data)
	require.NoError(t, err)

	for _, name := range []string{"images/pizza.png", "images/pizza.png"} {
		entry, err = writer.Create(name)
		require.NoError(t, err)
		_, err = entry.Write([]byte("x"))
		require.NoError(t, err)
	}
	require.NoError(t, writer.Close())

	path := filepath.Join(t.TempDir(), "backups.zip")
	require.NoError(t, os.WriteFile(path, buffer.Bytes(), 0o644))

	_, err = svc.PreviewZip(context.Background(), path)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidBackup))
}

func TestBackupImportBackupFailureDoesNotWipe(t *testing.T) {
	env, svc := newBackupTestService(t)
	seedBackupTestData(t, env)

	pkg := &BackupPackage{
		Format:     BackupFormat,
		Version:    BackupVersion,
		Materials:  []models.Material{{Id: "replacement", Name: "Replacement"}},
		Products:   []models.Product{},
		Categories: []models.Category{},
		Sales:      []models.SalesPerDay{},
		Orders:     []models.Order{},
	}
	path := filepath.Join(t.TempDir(), "backups.zip")
	writeBackupTestZip(t, path, pkg, nil)

	blocker := filepath.Join(t.TempDir(), "blocker")
	require.NoError(t, os.WriteFile(blocker, []byte("x"), 0o644))
	svc.Config.BackupsPath = blocker

	_, err := svc.ImportZip(context.Background(), path, BackupImportOptions{
		Sections: map[string]bool{backupSectionMaterials: true},
	})
	require.Error(t, err)

	assert.Equal(t, int64(2), countServiceDocs(t, env, "materials", bson.M{}))
}

func TestBackupDirFallback(t *testing.T) {
	uploads := t.TempDir()
	svc := &BackupService{Config: config.Config{UploadsPath: uploads}}

	assert.Equal(t, filepath.Join(filepath.Dir(uploads), "backups"), svc.BackupsDir())

	backups := t.TempDir()
	svc.Config.BackupsPath = backups
	assert.Equal(t, backups, svc.BackupsDir())
}
