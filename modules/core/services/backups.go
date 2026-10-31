package services

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/helpers"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// BackupFormat identifies a NutrixPOS backup archive.
	BackupFormat = "nutrixpos-backup"
	// BackupVersion is the supported backup archive format version.
	BackupVersion = 1
	// MaxBackupUploadBytes limits the size of an uploaded backup archive.
	MaxBackupUploadBytes int64 = 512 << 20

	backupFileName      = "backup.json"
	backupImagePrefix   = "images/"
	maxBackupJSONBytes  = int64(256 << 20)
	maxBackupImageBytes = int64(20 << 20)
	maxBackupTotalBytes = uint64(512 << 20)
	maxBackupEntries    = 10000
	maxBackupWarnings   = 100

	backupSectionMaterials  = "materials"
	backupSectionProducts   = "products"
	backupSectionCategories = "categories"
	backupSectionSales      = "sales"

	backupBatchSize = 500
)

// ErrInvalidBackup marks a malformed or unsupported backup archive. Handlers
// map this to HTTP 400 while database/write failures map to HTTP 500.
var ErrInvalidBackup = errors.New("invalid backup package")

// BackupSource describes where an archive was produced.
type BackupSource struct {
	Database string `json:"database"`
}

// BackupPackage is the JSON document stored inside a backups ZIP.
type BackupPackage struct {
	Format     string               `json:"format"`
	Version    int                  `json:"version"`
	ExportedAt time.Time            `json:"exported_at"`
	Source     BackupSource         `json:"source"`
	Materials  []models.Material    `json:"materials"`
	Products   []models.Product     `json:"products"`
	Categories []models.Category    `json:"categories"`
	Sales      []models.SalesPerDay `json:"sales"`
	Orders     []models.Order       `json:"orders"`
	Warnings   []string             `json:"warnings,omitempty"`
}

// BackupSectionCounts contains counts for one side of a preview (incoming or
// current database data).
type BackupSectionCounts struct {
	Materials       int `json:"materials"`
	MaterialEntries int `json:"material_entries"`
	Products        int `json:"products"`
	Categories      int `json:"categories"`
	Sales           int `json:"sales"`
	Orders          int `json:"orders"`
	Images          int `json:"images"`
}

// BackupPreview is returned before a destructive import so the UI can show
// exactly what will be replaced.
type BackupPreview struct {
	Format     string              `json:"format"`
	Version    int                 `json:"version"`
	ExportedAt time.Time           `json:"exported_at"`
	Source     BackupSource        `json:"source"`
	Sections   BackupSectionCounts `json:"sections"`
	Current    BackupSectionCounts `json:"current"`
	HasStock   bool                `json:"has_stock"`
	HasImages  bool                `json:"has_images"`
	Warnings   []string            `json:"warnings"`
}

// BackupImportOptions controls which sections are restored.
type BackupImportOptions struct {
	Sections      map[string]bool
	IncludeStock  bool
	IncludeImages bool
}

// Has reports whether the named section was selected.
func (o BackupImportOptions) Has(section string) bool {
	return o.Sections[section]
}

// BackupImportResult summarizes a completed import.
type BackupImportResult struct {
	BackupFile string              `json:"backup_file"`
	Imported   BackupSectionCounts `json:"imported"`
	Warnings   []string            `json:"warnings"`
}

// BackupService reads and writes backup archives.
type BackupService struct {
	Logger logger.ILogger
	Config config.Config
}

func (cs *BackupService) database() (*mongo.Database, error) {
	client, err := common.GetDatabaseClient(cs.Logger, &cs.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(cs.Config.Databases[0].Database), nil
}

func (cs *BackupService) salesCollectionName() string {
	if name := cs.Config.Databases[0].Tables["sales"]; name != "" {
		return name
	}
	return "sales"
}

// BackupsDir returns the directory where pre-import backups are stored. When
// backups_path is unset it is a sibling of uploads_path, so it is never served
// by the public uploads file server.
func (cs *BackupService) BackupsDir() string {
	if cs.Config.BackupsPath != "" {
		return helpers.ResolveOsEnvPath(cs.Config.BackupsPath)
	}

	uploads := helpers.ResolveOsEnvPath(cs.Config.UploadsPath)
	if uploads == "" {
		uploads = "."
	}
	return filepath.Join(filepath.Dir(uploads), "backups")
}

// BuildExport reads every backup section from the database.
func (cs *BackupService) BuildExport(ctx context.Context) (BackupPackage, error) {
	db, err := cs.database()
	if err != nil {
		return BackupPackage{}, err
	}

	pkg := BackupPackage{
		Format:     BackupFormat,
		Version:    BackupVersion,
		ExportedAt: time.Now().UTC(),
		Source:     BackupSource{Database: cs.Config.Databases[0].Database},
		Materials:  make([]models.Material, 0),
		Products:   make([]models.Product, 0),
		Categories: make([]models.Category, 0),
		Sales:      make([]models.SalesPerDay, 0),
		Orders:     make([]models.Order, 0),
	}

	if err := findAllDocuments(ctx, db.Collection("materials"), bson.D{{Key: "name", Value: 1}}, &pkg.Materials); err != nil {
		return pkg, err
	}
	if err := findAllDocuments(ctx, db.Collection("recipes"), bson.D{{Key: "name", Value: 1}}, &pkg.Products); err != nil {
		return pkg, err
	}
	if err := findAllDocuments(ctx, db.Collection("categories"), bson.D{{Key: "name", Value: 1}}, &pkg.Categories); err != nil {
		return pkg, err
	}
	if err := findAllDocuments(ctx, db.Collection(cs.salesCollectionName()), bson.D{{Key: "date", Value: 1}}, &pkg.Sales); err != nil {
		return pkg, err
	}
	if err := findAllDocuments(ctx, db.Collection("orders"), bson.D{{Key: "submitted_at", Value: 1}}, &pkg.Orders); err != nil {
		return pkg, err
	}

	if pkg.Materials == nil {
		pkg.Materials = make([]models.Material, 0)
	}
	if pkg.Products == nil {
		pkg.Products = make([]models.Product, 0)
	}
	if pkg.Categories == nil {
		pkg.Categories = make([]models.Category, 0)
	}
	if pkg.Sales == nil {
		pkg.Sales = make([]models.SalesPerDay, 0)
	}
	if pkg.Orders == nil {
		pkg.Orders = make([]models.Order, 0)
	}

	for i := range pkg.Materials {
		var quantity float64
		for _, entry := range pkg.Materials[i].Entries {
			quantity += entry.Quantity
		}
		pkg.Materials[i].Quantity = quantity
	}

	uploads := helpers.ResolveOsEnvPath(cs.Config.UploadsPath)
	for _, product := range pkg.Products {
		if product.ImageURL == "" {
			continue
		}
		base := filepath.Base(product.ImageURL)
		if _, err := os.Stat(filepath.Join(uploads, base)); err != nil {
			addBackupWarning(&pkg.Warnings, "product %q references missing image file %q", product.Name, base)
		}
	}

	return pkg, nil
}

func findAllDocuments(ctx context.Context, collection *mongo.Collection, sortDoc bson.D, out interface{}) error {
	opts := options.Find()
	if len(sortDoc) > 0 {
		opts.SetSort(sortDoc)
	}

	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return err
	}
	defer func() { _ = cursor.Close(ctx) }()

	return cursor.All(ctx, out)
}

// WriteBackupZIP writes the package and its referenced images to dstPath.
func (cs *BackupService) WriteBackupZIP(ctx context.Context, pkg *BackupPackage, dstPath string) error {
	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}

	cleanup := true
	defer func() {
		if cleanup {
			_ = out.Close()
			_ = os.Remove(dstPath)
		}
	}()

	zw := zip.NewWriter(out)

	uploads := helpers.ResolveOsEnvPath(cs.Config.UploadsPath)
	images := make(map[string]struct{})
	for _, product := range pkg.Products {
		if product.ImageURL == "" {
			continue
		}

		base := filepath.Base(product.ImageURL)
		if base == "" || base == "." || base == string(filepath.Separator) {
			continue
		}
		if _, exists := images[base]; exists {
			continue
		}
		images[base] = struct{}{}

		src, err := os.Open(filepath.Join(uploads, base))
		if err != nil {
			continue
		}

		dst, err := zw.Create(backupImagePrefix + base)
		if err != nil {
			_ = src.Close()
			return err
		}
		if _, err := io.Copy(dst, src); err != nil {
			_ = src.Close()
			return err
		}
		_ = src.Close()
	}

	data, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return err
	}

	entry, err := zw.Create(backupFileName)
	if err != nil {
		return err
	}
	if _, err := entry.Write(data); err != nil {
		return err
	}

	if err := zw.Close(); err != nil {
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}

	cleanup = false
	return nil
}

// Backup writes a full snapshot of the current data to the backups directory
// and returns its file name. It is called before any destructive import.
func (cs *BackupService) Backup(ctx context.Context) (string, error) {
	pkg, err := cs.BuildExport(ctx)
	if err != nil {
		return "", err
	}

	dir := cs.BackupsDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	name := fmt.Sprintf("backup-%s.zip", time.Now().UTC().Format("20060102-150405.000000000"))
	dst := filepath.Join(dir, name)
	if err := cs.WriteBackupZIP(ctx, &pkg, dst); err != nil {
		return "", err
	}

	return name, nil
}

// PreviewZip reads an archive and returns counts for the incoming data and the
// current database data without changing anything.
func (cs *BackupService) PreviewZip(ctx context.Context, zipPath string) (BackupPreview, error) {
	pkg, images, err := cs.readBackupZip(zipPath)
	if err != nil {
		return BackupPreview{}, err
	}

	warnings, err := validateBackupPackage(pkg)
	if err != nil {
		return BackupPreview{}, err
	}
	warnings = append(warnings, pkg.Warnings...)

	current, err := cs.CurrentCounts(ctx)
	if err != nil {
		return BackupPreview{}, err
	}

	hasStock := false
	for _, material := range pkg.Materials {
		if len(material.Entries) > 0 {
			hasStock = true
			break
		}
	}

	return BackupPreview{
		Format:     pkg.Format,
		Version:    pkg.Version,
		ExportedAt: pkg.ExportedAt,
		Source:     pkg.Source,
		Sections:   backupCountsFromPackage(pkg, len(images)),
		Current:    current,
		HasStock:   hasStock,
		HasImages:  len(images) > 0,
		Warnings:   warnings,
	}, nil
}

// ImportZip replaces the backup data (and Sales & Orders when selected) with the
// archive contents. A full backup is created before the first delete.
func (cs *BackupService) ImportZip(ctx context.Context, zipPath string, opts BackupImportOptions) (BackupImportResult, error) {
	result := BackupImportResult{
		Imported: BackupSectionCounts{},
		Warnings: make([]string, 0),
	}

	pkg, images, err := cs.readBackupZip(zipPath)
	if err != nil {
		return result, err
	}

	warnings, err := validateBackupPackage(pkg)
	if err != nil {
		return result, err
	}
	warnings = append(warnings, pkg.Warnings...)

	if err := validateBackupImportSelection(opts, pkg); err != nil {
		return result, err
	}

	backup, err := cs.Backup(ctx)
	if err != nil {
		return result, fmt.Errorf("failed to create backup before import: %w", err)
	}
	result.BackupFile = backup

	if err := cs.replaceBackup(ctx, pkg, images, opts, &result, &warnings); err != nil {
		return result, err
	}

	result.Warnings = warnings
	return result, nil
}

func (cs *BackupService) replaceBackup(ctx context.Context, pkg *BackupPackage, images map[string][]byte, opts BackupImportOptions, result *BackupImportResult, warnings *[]string) error {
	db, err := cs.database()
	if err != nil {
		return err
	}

	wipes := make([]string, 0, 5)
	if opts.Has(backupSectionMaterials) {
		wipes = append(wipes, "materials", "recipes", "categories")
	}
	if opts.Has(backupSectionSales) {
		wipes = append(wipes, cs.salesCollectionName(), "orders")
	}

	for _, name := range wipes {
		if _, err := db.Collection(name).DeleteMany(ctx, bson.M{}); err != nil {
			return fmt.Errorf("failed to clear %s: %w", name, err)
		}
	}

	if opts.Has(backupSectionMaterials) {
		docs := make([]interface{}, 0, len(pkg.Materials))
		for i := range pkg.Materials {
			material := pkg.Materials[i]
			if !opts.IncludeStock {
				material.Entries = []models.MaterialEntry{}
				material.Quantity = 0
			} else {
				var quantity float64
				for _, entry := range material.Entries {
					quantity += entry.Quantity
				}
				material.Quantity = quantity
				result.Imported.MaterialEntries += len(material.Entries)
			}
			docs = append(docs, material)
		}

		if err := insertBackupDocs(ctx, db.Collection("materials"), docs); err != nil {
			return fmt.Errorf("failed to restore materials: %w", err)
		}
		result.Imported.Materials = len(docs)
	}

	if opts.Has(backupSectionProducts) {
		docs := make([]interface{}, 0, len(pkg.Products))
		for _, product := range pkg.Products {
			docs = append(docs, product)
		}

		if err := insertBackupDocs(ctx, db.Collection("recipes"), docs); err != nil {
			return fmt.Errorf("failed to restore products: %w", err)
		}
		result.Imported.Products = len(docs)

		if opts.IncludeImages {
			result.Imported.Images = cs.restoreBackupImages(pkg, images, warnings)
			for _, product := range pkg.Products {
				if product.ImageURL == "" {
					continue
				}
				base := filepath.Base(product.ImageURL)
				if _, ok := images[base]; !ok {
					addBackupWarning(warnings, "product %q references missing image %q in the archive", product.Name, base)
				}
			}
		}
	}

	if opts.Has(backupSectionCategories) {
		docs := make([]interface{}, 0, len(pkg.Categories))
		for _, category := range pkg.Categories {
			docs = append(docs, category)
		}

		if err := insertBackupDocs(ctx, db.Collection("categories"), docs); err != nil {
			return fmt.Errorf("failed to restore categories: %w", err)
		}
		result.Imported.Categories = len(docs)
	}

	if opts.Has(backupSectionSales) {
		orderDocs := make([]interface{}, 0, len(pkg.Orders))
		for _, order := range pkg.Orders {
			orderDocs = append(orderDocs, order)
		}
		if err := insertBackupDocs(ctx, db.Collection("orders"), orderDocs); err != nil {
			return fmt.Errorf("failed to restore orders: %w", err)
		}

		salesDocs := make([]interface{}, 0, len(pkg.Sales))
		for _, sale := range pkg.Sales {
			salesDocs = append(salesDocs, sale)
		}
		if err := insertBackupDocs(ctx, db.Collection(cs.salesCollectionName()), salesDocs); err != nil {
			return fmt.Errorf("failed to restore sales: %w", err)
		}

		result.Imported.Orders = len(orderDocs)
		result.Imported.Sales = len(salesDocs)

		if err := cs.advanceOrderDisplayCounters(ctx, pkg.Orders); err != nil {
			addBackupWarning(warnings, "could not advance order display id counters: %v", err)
		}
	}

	return nil
}

func insertBackupDocs(ctx context.Context, collection *mongo.Collection, docs []interface{}) error {
	for start := 0; start < len(docs); start += backupBatchSize {
		end := start + backupBatchSize
		if end > len(docs) {
			end = len(docs)
		}
		if _, err := collection.InsertMany(ctx, docs[start:end]); err != nil {
			return err
		}
	}
	return nil
}

func (cs *BackupService) restoreBackupImages(pkg *BackupPackage, images map[string][]byte, warnings *[]string) int {
	uploads := helpers.ResolveOsEnvPath(cs.Config.UploadsPath)
	if uploads == "" {
		uploads = "."
	}
	if err := os.MkdirAll(uploads, 0o755); err != nil {
		addBackupWarning(warnings, "could not create uploads directory: %v", err)
		return 0
	}

	written := make(map[string]struct{})
	restored := 0
	for _, product := range pkg.Products {
		if product.ImageURL == "" {
			continue
		}

		base := filepath.Base(product.ImageURL)
		if base == "" || base == "." || base == string(filepath.Separator) {
			continue
		}
		if _, exists := written[base]; exists {
			continue
		}
		written[base] = struct{}{}

		data, ok := images[base]
		if !ok {
			continue
		}
		if err := os.WriteFile(filepath.Join(uploads, base), data, 0o644); err != nil {
			addBackupWarning(warnings, "could not restore image %q: %v", base, err)
			continue
		}
		restored++
	}
	return restored
}

func (cs *BackupService) advanceOrderDisplayCounters(ctx context.Context, orders []models.Order) error {
	if len(orders) == 0 {
		return nil
	}

	orderDisplayIdMu.Lock()
	defer orderDisplayIdMu.Unlock()

	settingsSvc := SettingsService{Logger: cs.Logger, Config: cs.Config}
	settings, err := settingsSvc.GetSettings()
	if err != nil {
		return err
	}

	maximum := make(map[string]uint32)
	for _, order := range orders {
		separator := strings.LastIndex(order.DisplayId, "-")
		if separator <= 0 || separator == len(order.DisplayId)-1 {
			continue
		}

		prefix := order.DisplayId[:separator]
		number, err := strconv.Atoi(order.DisplayId[separator+1:])
		if err != nil || number < 0 {
			continue
		}

		next := uint32(number + 1)
		if next > maximum[prefix] {
			maximum[prefix] = next
		}
	}

	changed := false
	for i := range settings.Orders.Queues {
		if next, ok := maximum[settings.Orders.Queues[i].Prefix]; ok && next > settings.Orders.Queues[i].Next {
			settings.Orders.Queues[i].Next = next
			changed = true
		}
	}
	if !changed {
		return nil
	}

	db, err := cs.database()
	if err != nil {
		return err
	}

	filter := bson.M{}
	if settings.Id != "" {
		filter["id"] = settings.Id
	}
	_, err = db.Collection("settings").UpdateOne(ctx, filter, bson.M{
		"$set": bson.M{"orders.queues": settings.Orders.Queues},
	})
	return err
}

// CurrentCounts returns counts for the data currently stored in the database.
func (cs *BackupService) CurrentCounts(ctx context.Context) (BackupSectionCounts, error) {
	db, err := cs.database()
	if err != nil {
		return BackupSectionCounts{}, err
	}

	var counts BackupSectionCounts

	materialsCursor, err := db.Collection("materials").Find(ctx, bson.D{})
	if err != nil {
		return counts, err
	}
	for materialsCursor.Next(ctx) {
		var material models.Material
		if err := materialsCursor.Decode(&material); err != nil {
			_ = materialsCursor.Close(ctx)
			return counts, err
		}
		counts.Materials++
		counts.MaterialEntries += len(material.Entries)
	}
	if err := materialsCursor.Err(); err != nil {
		_ = materialsCursor.Close(ctx)
		return counts, err
	}
	_ = materialsCursor.Close(ctx)

	productsCursor, err := db.Collection("recipes").Find(ctx, bson.D{})
	if err != nil {
		return counts, err
	}
	for productsCursor.Next(ctx) {
		var product models.Product
		if err := productsCursor.Decode(&product); err != nil {
			_ = productsCursor.Close(ctx)
			return counts, err
		}
		counts.Products++
		if product.ImageURL != "" {
			counts.Images++
		}
	}
	if err := productsCursor.Err(); err != nil {
		_ = productsCursor.Close(ctx)
		return counts, err
	}
	_ = productsCursor.Close(ctx)

	categories, err := db.Collection("categories").CountDocuments(ctx, bson.M{})
	if err != nil {
		return counts, err
	}
	counts.Categories = int(categories)

	sales, err := db.Collection(cs.salesCollectionName()).CountDocuments(ctx, bson.M{})
	if err != nil {
		return counts, err
	}
	counts.Sales = int(sales)

	orders, err := db.Collection("orders").CountDocuments(ctx, bson.M{})
	if err != nil {
		return counts, err
	}
	counts.Orders = int(orders)

	return counts, nil
}

func (cs *BackupService) readBackupZip(zipPath string) (*BackupPackage, map[string][]byte, error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", ErrInvalidBackup, err)
	}
	defer func() { _ = reader.Close() }()

	if len(reader.File) > maxBackupEntries {
		return nil, nil, fmt.Errorf("%w: archive contains too many entries", ErrInvalidBackup)
	}

	var pkg *BackupPackage
	images := make(map[string][]byte)
	var total uint64

	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		name := path.Clean(strings.ReplaceAll(file.Name, "\\", "/"))

		switch {
		case name == backupFileName:
			if pkg != nil {
				return nil, nil, fmt.Errorf("%w: archive contains more than one %s", ErrInvalidBackup, backupFileName)
			}

			data, err := readZipEntry(file, maxBackupJSONBytes)
			if err != nil {
				return nil, nil, err
			}
			total += uint64(len(data))
			if total > maxBackupTotalBytes {
				return nil, nil, fmt.Errorf("%w: archive is too large when uncompressed", ErrInvalidBackup)
			}

			var decoded BackupPackage
			if err := json.Unmarshal(data, &decoded); err != nil {
				return nil, nil, fmt.Errorf("%w: invalid %s: %v", ErrInvalidBackup, backupFileName, err)
			}
			pkg = &decoded

		case strings.HasPrefix(name, backupImagePrefix):
			base := path.Base(name)
			if base == "." || base == ".." || base == "/" || strings.Contains(base, "/") || name != backupImagePrefix+base {
				return nil, nil, fmt.Errorf("%w: invalid image path %q", ErrInvalidBackup, file.Name)
			}
			if _, exists := images[base]; exists {
				return nil, nil, fmt.Errorf("%w: duplicate image %q", ErrInvalidBackup, base)
			}

			data, err := readZipEntry(file, maxBackupImageBytes)
			if err != nil {
				return nil, nil, err
			}
			total += uint64(len(data))
			if total > maxBackupTotalBytes {
				return nil, nil, fmt.Errorf("%w: archive is too large when uncompressed", ErrInvalidBackup)
			}
			images[base] = data
		}
	}

	if pkg == nil {
		return nil, nil, fmt.Errorf("%w: archive does not contain %s", ErrInvalidBackup, backupFileName)
	}

	return pkg, images, nil
}

func readZipEntry(file *zip.File, limit int64) ([]byte, error) {
	if file.UncompressedSize64 > uint64(limit) {
		return nil, fmt.Errorf("%w: entry %q is too large", ErrInvalidBackup, file.Name)
	}

	reader, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("%w: could not read entry %q: %v", ErrInvalidBackup, file.Name, err)
	}
	defer func() { _ = reader.Close() }()

	data, err := io.ReadAll(io.LimitReader(reader, limit+1))
	if err != nil {
		return nil, fmt.Errorf("%w: could not read entry %q: %v", ErrInvalidBackup, file.Name, err)
	}
	if int64(len(data)) > limit {
		return nil, fmt.Errorf("%w: entry %q is too large", ErrInvalidBackup, file.Name)
	}

	return data, nil
}

func validateBackupPackage(pkg *BackupPackage) ([]string, error) {
	if pkg.Format != BackupFormat {
		return nil, fmt.Errorf("%w: unsupported format %q", ErrInvalidBackup, pkg.Format)
	}
	if pkg.Version != BackupVersion {
		return nil, fmt.Errorf("%w: unsupported version %d", ErrInvalidBackup, pkg.Version)
	}

	warnings := make([]string, 0)

	materialIDs := make(map[string]struct{}, len(pkg.Materials))
	for i := range pkg.Materials {
		material := &pkg.Materials[i]
		if strings.TrimSpace(material.Id) == "" {
			return nil, fmt.Errorf("%w: material at index %d has no id", ErrInvalidBackup, i+1)
		}
		if _, exists := materialIDs[material.Id]; exists {
			return nil, fmt.Errorf("%w: duplicate material id %q", ErrInvalidBackup, material.Id)
		}
		materialIDs[material.Id] = struct{}{}
	}

	productIDs := make(map[string]struct{}, len(pkg.Products))
	for i := range pkg.Products {
		product := &pkg.Products[i]
		if strings.TrimSpace(product.Id) == "" {
			return nil, fmt.Errorf("%w: product at index %d has no id", ErrInvalidBackup, i+1)
		}
		if _, exists := productIDs[product.Id]; exists {
			return nil, fmt.Errorf("%w: duplicate product id %q", ErrInvalidBackup, product.Id)
		}
		productIDs[product.Id] = struct{}{}
	}

	categoryIDs := make(map[string]struct{}, len(pkg.Categories))
	for i := range pkg.Categories {
		category := &pkg.Categories[i]
		if strings.TrimSpace(category.Id) == "" {
			return nil, fmt.Errorf("%w: category at index %d has no id", ErrInvalidBackup, i+1)
		}
		if _, exists := categoryIDs[category.Id]; exists {
			return nil, fmt.Errorf("%w: duplicate category id %q", ErrInvalidBackup, category.Id)
		}
		categoryIDs[category.Id] = struct{}{}
	}

	orderIDs := make(map[string]struct{}, len(pkg.Orders))
	for i := range pkg.Orders {
		order := &pkg.Orders[i]
		if strings.TrimSpace(order.Id) == "" {
			return nil, fmt.Errorf("%w: order at index %d has no id", ErrInvalidBackup, i+1)
		}
		if _, exists := orderIDs[order.Id]; exists {
			return nil, fmt.Errorf("%w: duplicate order id %q", ErrInvalidBackup, order.Id)
		}
		orderIDs[order.Id] = struct{}{}
	}

	for _, product := range pkg.Products {
		for _, material := range product.Materials {
			if _, ok := materialIDs[material.Id]; !ok {
				addBackupWarning(&warnings, "product %q references missing material id %q", product.Name, material.Id)
			}
		}
		for _, subProduct := range product.SubProducts {
			if _, ok := productIDs[subProduct.Id]; !ok {
				addBackupWarning(&warnings, "product %q references missing subproduct id %q", product.Name, subProduct.Id)
			}
		}
	}

	for _, category := range pkg.Categories {
		for _, product := range category.Products {
			if _, ok := productIDs[product.Id]; !ok {
				addBackupWarning(&warnings, "category %q references missing product id %q", category.Name, product.Id)
			}
		}
	}

	return warnings, nil
}

func validateBackupImportSelection(opts BackupImportOptions, pkg *BackupPackage) error {
	if len(opts.Sections) == 0 {
		return fmt.Errorf("%w: no sections selected", ErrInvalidBackup)
	}

	knownSections := map[string]struct{}{
		backupSectionMaterials:  {},
		backupSectionProducts:   {},
		backupSectionCategories: {},
		backupSectionSales:      {},
	}
	for section := range opts.Sections {
		if _, ok := knownSections[section]; !ok {
			return fmt.Errorf("%w: unknown section %q", ErrInvalidBackup, section)
		}
	}

	if opts.Has(backupSectionProducts) && !opts.Has(backupSectionMaterials) {
		return fmt.Errorf("%w: importing products requires importing materials", ErrInvalidBackup)
	}
	if opts.Has(backupSectionCategories) && !opts.Has(backupSectionProducts) {
		return fmt.Errorf("%w: importing categories requires importing products", ErrInvalidBackup)
	}
	if opts.IncludeStock && !opts.Has(backupSectionMaterials) {
		return fmt.Errorf("%w: importing stock entries requires importing materials", ErrInvalidBackup)
	}
	if opts.IncludeImages && !opts.Has(backupSectionProducts) {
		return fmt.Errorf("%w: importing images requires importing products", ErrInvalidBackup)
	}

	if opts.Has(backupSectionMaterials) && pkg.Materials == nil {
		return fmt.Errorf("%w: archive has no materials section", ErrInvalidBackup)
	}
	if opts.Has(backupSectionProducts) && pkg.Products == nil {
		return fmt.Errorf("%w: archive has no products section", ErrInvalidBackup)
	}
	if opts.Has(backupSectionCategories) && pkg.Categories == nil {
		return fmt.Errorf("%w: archive has no categories section", ErrInvalidBackup)
	}
	if opts.Has(backupSectionSales) && pkg.Sales == nil && pkg.Orders == nil {
		return fmt.Errorf("%w: archive has no sales or orders section", ErrInvalidBackup)
	}

	return nil
}

func backupCountsFromPackage(pkg *BackupPackage, images int) BackupSectionCounts {
	counts := BackupSectionCounts{
		Materials:  len(pkg.Materials),
		Products:   len(pkg.Products),
		Categories: len(pkg.Categories),
		Sales:      len(pkg.Sales),
		Orders:     len(pkg.Orders),
		Images:     images,
	}
	for _, material := range pkg.Materials {
		counts.MaterialEntries += len(material.Entries)
	}
	return counts
}

func addBackupWarning(warnings *[]string, format string, args ...interface{}) {
	if len(*warnings) >= maxBackupWarnings {
		return
	}
	*warnings = append(*warnings, fmt.Sprintf(format, args...))
}
