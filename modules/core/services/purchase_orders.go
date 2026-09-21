package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/customerrors"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PurchaseOrderService provides methods to manage purchase orders and their
// automated goods received notes (GRNs). Receiving a purchase order pushes the
// received quantities into the materials inventory and records them in the
// inventory item history.
type PurchaseOrderService struct {
	Logger   logger.ILogger
	Config   config.Config
	Settings models.Settings
}

// GetPurchaseOrdersParams are the pagination/filter parameters for listing purchase orders.
type GetPurchaseOrdersParams struct {
	PageNumber int
	PageSize   int
	Search     string
	Status     string
}

// ReceiveItem describes how much of a purchase order item should be received.
type ReceiveItem struct {
	ItemId   string  `json:"item_id" bson:"item_id" mapstructure:"item_id"`
	Quantity float64 `json:"quantity" bson:"quantity" mapstructure:"quantity"`
}

// GetPurchaseOrdersParams are the pagination/filter parameters for listing GRNs.
type GetGRNsParams struct {
	PageNumber      int
	PageSize        int
	Search          string
	PurchaseOrderId string
}

func (ps *PurchaseOrderService) getNextSequence(ctx context.Context, name string) (seq int64, err error) {
	client, err := common.GetDatabaseClient(ps.Logger, &ps.Config)
	if err != nil {
		return 0, err
	}

	collection := client.Database(ps.Config.Databases[0].Database).Collection("counters")
	filter := bson.M{"name": name}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var doc struct {
		Seq int64 `bson:"seq"`
	}
	err = collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc)
	if err != nil {
		return 0, err
	}

	return doc.Seq, nil
}

// CreatePurchaseOrder validates and stores a purchase order. When AutoReceive is
// true the goods are received immediately, generating a GRN and adding the
// quantities to inventory.
func (ps *PurchaseOrderService) CreatePurchaseOrder(po models.PurchaseOrder, user_id string) (models.PurchaseOrder, error) {
	if len(po.Items) == 0 {
		return po, fmt.Errorf("purchase order must have at least one item")
	}

	client, err := common.GetDatabaseClient(ps.Logger, &ps.Config)
	if err != nil {
		return po, err
	}

	ctx := context.Background()

	db := client.Database(ps.Config.Databases[0].Database)
	materialsCollection := db.Collection("materials")

	total := 0.0
	for i := range po.Items {
		item := &po.Items[i]

		if item.MaterialId == "" {
			return po, fmt.Errorf("item material is required")
		}
		if item.Quantity <= 0 {
			return po, fmt.Errorf("item quantity must be greater than 0")
		}
		if item.PurchasePrice < 0 {
			return po, fmt.Errorf("item purchase price cannot be negative")
		}
		item.ItemId = primitive.NewObjectID().Hex()
		item.ReceivedQuantity = 0
		item.EntryIds = []string{}

		var material models.Material
		err = materialsCollection.FindOne(ctx, bson.M{"id": item.MaterialId}).Decode(&material)
		if err != nil {
			return po, fmt.Errorf("material %s not found", item.MaterialId)
		}

		item.MaterialName = material.Name
		item.Unit = material.Unit
		total += item.Quantity * item.PurchasePrice
	}

	po.Id = primitive.NewObjectID().Hex()

	seq, err := ps.getNextSequence(ctx, "purchase_order")
	if err != nil {
		return po, err
	}
	po.DisplayId = fmt.Sprintf("PO-%04d", seq)

	po.Status = models.PurchaseOrderStatusOpen
	if po.CreatedAt.IsZero() {
		po.CreatedAt = time.Now()
	}
	po.CreatedBy = user_id
	po.Total = total

	collection := db.Collection("purchase_orders")
	_, err = collection.InsertOne(ctx, po)
	if err != nil {
		return po, err
	}

	if po.AutoReceive {
		_, err = ps.ReceivePurchaseOrder(po.Id, user_id, nil)
		if err != nil {
			return po, err
		}

		err = collection.FindOne(ctx, bson.M{"id": po.Id}).Decode(&po)
		if err != nil {
			return po, err
		}
	}

	return po, nil
}

// GetPurchaseOrders returns a paginated list of purchase orders.
func (ps *PurchaseOrderService) GetPurchaseOrders(params GetPurchaseOrdersParams) (pos []models.PurchaseOrder, totalRecords int64, err error) {
	pos = make([]models.PurchaseOrder, 0)

	client, err := common.GetDatabaseClient(ps.Logger, &ps.Config)
	if err != nil {
		return pos, totalRecords, err
	}

	ctx := context.Background()

	filter := bson.M{}
	if params.Status != "" {
		filter["status"] = params.Status
	}
	if params.Search != "" {
		regex := bson.M{"$regex": primitive.Regex{Pattern: params.Search, Options: "i"}}
		filter["$or"] = []bson.M{
			{"supplier": regex},
			{"display_id": regex},
		}
	}

	if params.PageSize <= 0 {
		params.PageSize = 50
	}
	if params.PageNumber < 1 {
		params.PageNumber = 1
	}

	collection := client.Database(ps.Config.Databases[0].Database).Collection("purchase_orders")

	totalRecords, err = collection.CountDocuments(ctx, filter)
	if err != nil {
		return pos, totalRecords, err
	}

	findOptions := options.Find().
		SetSkip(int64((params.PageNumber - 1) * params.PageSize)).
		SetLimit(int64(params.PageSize)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return pos, totalRecords, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &pos); err != nil {
		return pos, totalRecords, err
	}

	return pos, totalRecords, nil
}

// GetPurchaseOrder returns a single purchase order by id.
func (ps *PurchaseOrderService) GetPurchaseOrder(purchase_order_id string) (po models.PurchaseOrder, err error) {
	client, err := common.GetDatabaseClient(ps.Logger, &ps.Config)
	if err != nil {
		return po, err
	}

	ctx := context.Background()

	collection := client.Database(ps.Config.Databases[0].Database).Collection("purchase_orders")
	err = collection.FindOne(ctx, bson.M{"id": purchase_order_id}).Decode(&po)
	if err != nil {
		return po, err
	}

	return po, nil
}

// CancelPurchaseOrder cancels an open purchase order. Received orders cannot be cancelled.
func (ps *PurchaseOrderService) CancelPurchaseOrder(purchase_order_id string) (err error) {
	client, err := common.GetDatabaseClient(ps.Logger, &ps.Config)
	if err != nil {
		return err
	}

	ctx := context.Background()

	collection := client.Database(ps.Config.Databases[0].Database).Collection("purchase_orders")

	var po models.PurchaseOrder
	err = collection.FindOne(ctx, bson.M{"id": purchase_order_id}).Decode(&po)
	if err != nil {
		return err
	}

	if po.Status != models.PurchaseOrderStatusOpen {
		return fmt.Errorf("only open purchase orders can be cancelled")
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"id": purchase_order_id, "status": models.PurchaseOrderStatusOpen},
		bson.M{"$set": bson.M{"status": models.PurchaseOrderStatusCancelled}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("only open purchase orders can be cancelled")
	}

	return nil
}

// ReceivePurchaseOrder receives goods against an open purchase order. It
// generates a GRN, pushes the received quantities as new material entries into
// the materials inventory and writes a history log for each received material.
// Passing an empty received_items receives all remaining quantities.
func (ps *PurchaseOrderService) ReceivePurchaseOrder(purchase_order_id string, user_id string, received_items []ReceiveItem) (grn models.GRN, err error) {
	client, err := common.GetDatabaseClient(ps.Logger, &ps.Config)
	if err != nil {
		return grn, err
	}

	ctx := context.Background()

	db := client.Database(ps.Config.Databases[0].Database)
	poCollection := db.Collection("purchase_orders")
	materialsCollection := db.Collection("materials")
	logsCollection := db.Collection("logs")
	grnsCollection := db.Collection("grns")

	var po models.PurchaseOrder
	err = poCollection.FindOne(ctx, bson.M{"id": purchase_order_id}).Decode(&po)
	if err != nil {
		return grn, err
	}

	if po.Status == models.PurchaseOrderStatusReceived {
		return grn, fmt.Errorf("purchase order %s is already fully received", po.DisplayId)
	}
	if po.Status == models.PurchaseOrderStatusCancelled {
		return grn, fmt.Errorf("purchase order %s is cancelled", po.DisplayId)
	}

	requested := map[string]float64{}
	for _, item := range received_items {
		if item.Quantity < 0 {
			return grn, fmt.Errorf("received quantity cannot be negative")
		}
		requested[item.ItemId] = item.Quantity
	}

	explicit := len(received_items) > 0

	seq, err := ps.getNextSequence(ctx, "grn")
	if err != nil {
		return grn, err
	}

	now := time.Now()

	grn = models.GRN{
		Id:                     primitive.NewObjectID().Hex(),
		DisplayId:              fmt.Sprintf("GRN-%04d", seq),
		PurchaseOrderId:        po.Id,
		PurchaseOrderDisplayId: po.DisplayId,
		Supplier:               po.Supplier,
		ReceivedAt:             now,
		ReceivedBy:             user_id,
		Items:                  []models.GRNItem{},
	}

	type pendingReceive struct {
		item      *models.PurchaseOrderItem
		toReceive float64
		entry     models.MaterialEntry
		logDoc    bson.M
	}

	pending := []pendingReceive{}

	for i := range po.Items {
		poItem := &po.Items[i]

		remaining := poItem.Quantity - poItem.ReceivedQuantity
		if remaining <= 0 {
			continue
		}

		toReceive := remaining
		if explicit {
			quantity, ok := requested[poItem.ItemId]
			if !ok {
				continue
			}
			toReceive = quantity
		}
		if toReceive <= 0 {
			continue
		}
		if toReceive > remaining {
			toReceive = remaining
		}

		entry := models.MaterialEntry{
			Id:               primitive.NewObjectID().Hex(),
			PurchaseQuantity: toReceive,
			PurchasePrice:    poItem.PurchasePrice,
			Quantity:         toReceive,
			Company:          poItem.Company,
			SKU:              poItem.SKU,
			ExpirationDate:   poItem.ExpirationDate,
		}

		grn.Items = append(grn.Items, models.GRNItem{
			ItemId:         poItem.ItemId,
			MaterialId:     poItem.MaterialId,
			MaterialName:   poItem.MaterialName,
			Unit:           poItem.Unit,
			Company:        poItem.Company,
			SKU:            poItem.SKU,
			EntryId:        entry.Id,
			Quantity:       toReceive,
			PurchasePrice:  poItem.PurchasePrice,
			ExpirationDate: poItem.ExpirationDate,
		})

		pending = append(pending, pendingReceive{
			item:      poItem,
			toReceive: toReceive,
			entry:     entry,
			logDoc: bson.M{
				"id":                        primitive.NewObjectID().Hex(),
				"type":                      models.LogTypeMaterialGRNReceive,
				"date":                      now,
				"user_id":                   user_id,
				"component_id":              poItem.MaterialId,
				"material_id":               poItem.MaterialId,
				"entry_id":                  entry.Id,
				"quantity":                  toReceive,
				"company":                   poItem.Company,
				"price":                     poItem.PurchasePrice,
				"grn_id":                    grn.Id,
				"grn_display_id":            grn.DisplayId,
				"purchase_order_id":         po.Id,
				"purchase_order_display_id": po.DisplayId,
			},
		})
	}

	if len(pending) == 0 {
		return grn, fmt.Errorf("no items to receive for purchase order %s", po.DisplayId)
	}

	materialIds := make([]string, 0, len(pending))
	seen := map[string]bool{}
	for _, p := range pending {
		if !seen[p.item.MaterialId] {
			seen[p.item.MaterialId] = true
			materialIds = append(materialIds, p.item.MaterialId)
		}
	}

	existingMaterials, err := materialsCollection.CountDocuments(ctx, bson.M{"id": bson.M{"$in": materialIds}})
	if err != nil {
		return grn, err
	}
	if int(existingMaterials) != len(materialIds) {
		return grn, fmt.Errorf("one or more materials no longer exist")
	}

	originalReceived := make(map[string]float64, len(po.Items))
	for _, item := range po.Items {
		originalReceived[item.ItemId] = item.ReceivedQuantity
	}

	for i := range pending {
		pending[i].item.ReceivedQuantity += pending[i].toReceive
		pending[i].item.EntryIds = append(pending[i].item.EntryIds, pending[i].entry.Id)
	}

	fully_received := true
	for _, item := range po.Items {
		if item.ReceivedQuantity < item.Quantity {
			fully_received = false
			break
		}
	}

	status := models.PurchaseOrderStatusPartial
	if fully_received {
		status = models.PurchaseOrderStatusReceived
	}

	guards := make([]bson.M, 0, len(po.Items))
	for _, item := range po.Items {
		guards = append(guards, bson.M{"items": bson.M{"$elemMatch": bson.M{
			"item_id":           item.ItemId,
			"received_quantity": originalReceived[item.ItemId],
		}}})
	}

	filter := bson.M{
		"id":     po.Id,
		"status": bson.M{"$nin": []string{models.PurchaseOrderStatusReceived, models.PurchaseOrderStatusCancelled}},
		"$and":   guards,
	}

	update := bson.M{
		"$set": bson.M{
			"status":      status,
			"received_at": now,
			"received_by": user_id,
			"items":       po.Items,
		},
	}

	result, err := poCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return grn, err
	}
	if result.MatchedCount == 0 {
		return grn, customerrors.ErrPurchaseOrderModified
	}

	for _, p := range pending {
		_, err = materialsCollection.UpdateOne(ctx, bson.M{"id": p.item.MaterialId}, bson.M{"$push": bson.M{"entries": p.entry}})
		if err != nil {
			return grn, err
		}

		_, err = logsCollection.InsertOne(ctx, p.logDoc)
		if err != nil {
			return grn, err
		}
	}

	_, err = grnsCollection.InsertOne(ctx, grn)
	if err != nil {
		return grn, err
	}

	return grn, nil
}

// GetGRNs returns a paginated list of goods received notes.
func (ps *PurchaseOrderService) GetGRNs(params GetGRNsParams) (grns []models.GRN, totalRecords int64, err error) {
	grns = make([]models.GRN, 0)

	client, err := common.GetDatabaseClient(ps.Logger, &ps.Config)
	if err != nil {
		return grns, totalRecords, err
	}

	ctx := context.Background()

	filter := bson.M{}
	if params.PurchaseOrderId != "" {
		filter["purchase_order_id"] = params.PurchaseOrderId
	}
	if params.Search != "" {
		regex := bson.M{"$regex": primitive.Regex{Pattern: params.Search, Options: "i"}}
		filter["$or"] = []bson.M{
			{"supplier": regex},
			{"display_id": regex},
			{"purchase_order_display_id": regex},
		}
	}

	if params.PageSize <= 0 {
		params.PageSize = 50
	}
	if params.PageNumber < 1 {
		params.PageNumber = 1
	}

	collection := client.Database(ps.Config.Databases[0].Database).Collection("grns")

	totalRecords, err = collection.CountDocuments(ctx, filter)
	if err != nil {
		return grns, totalRecords, err
	}

	findOptions := options.Find().
		SetSkip(int64((params.PageNumber - 1) * params.PageSize)).
		SetLimit(int64(params.PageSize)).
		SetSort(bson.M{"received_at": -1})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return grns, totalRecords, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &grns); err != nil {
		return grns, totalRecords, err
	}

	return grns, totalRecords, nil
}

// GetGRN returns a single goods received note by id.
func (ps *PurchaseOrderService) GetGRN(grn_id string) (grn models.GRN, err error) {
	client, err := common.GetDatabaseClient(ps.Logger, &ps.Config)
	if err != nil {
		return grn, err
	}

	ctx := context.Background()

	collection := client.Database(ps.Config.Databases[0].Database).Collection("grns")
	err = collection.FindOne(ctx, bson.M{"id": grn_id}).Decode(&grn)
	if err != nil {
		return grn, err
	}

	return grn, nil
}
