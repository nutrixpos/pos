export class PurchaseOrderItem {
	item_id: string;
	material_id: string;
	material_name: string;
	unit: string;
	company: string;
	sku: string;
	quantity: number;
	received_quantity: number;
	purchase_price: number;
	total_price: number;
	expiration_date: Date;
	entry_ids: string[];

	constructor() {
		this.item_id = ""
		this.material_id = ""
		this.material_name = ""
		this.unit = ""
		this.company = ""
		this.sku = ""
		this.quantity = 0
		this.received_quantity = 0
		this.purchase_price = 0
		this.total_price = 0
		this.expiration_date = null as any
		this.entry_ids = []
	}
}

export class PurchaseOrder {
	id: string;
	display_id: string;
	supplier: string;
	status: string;
	notes: string;
	items: PurchaseOrderItem[];
	total: number;
	auto_receive: boolean;
	created_at: Date;
	created_by: string;
	received_at: Date;
	received_by: string;

	constructor() {
		this.id = ""
		this.display_id = ""
		this.supplier = ""
		this.status = "open"
		this.notes = ""
		this.items = []
		this.total = 0
		this.auto_receive = false
		this.created_at = null as any
		this.created_by = ""
		this.received_at = null as any
		this.received_by = ""
	}
}

export class GRNItem {
	item_id: string;
	material_id: string;
	material_name: string;
	unit: string;
	company: string;
	sku: string;
	entry_id: string;
	quantity: number;
	purchase_price: number;
	expiration_date: Date;

	constructor() {
		this.item_id = ""
		this.material_id = ""
		this.material_name = ""
		this.unit = ""
		this.company = ""
		this.sku = ""
		this.entry_id = ""
		this.quantity = 0
		this.purchase_price = 0
		this.expiration_date = null as any
	}
}

export class GRN {
	id: string;
	display_id: string;
	purchase_order_id: string;
	purchase_order_display_id: string;
	supplier: string;
	items: GRNItem[];
	received_at: Date;
	received_by: string;

	constructor() {
		this.id = ""
		this.display_id = ""
		this.purchase_order_id = ""
		this.purchase_order_display_id = ""
		this.supplier = ""
		this.items = []
		this.received_at = null as any
		this.received_by = ""
	}
}