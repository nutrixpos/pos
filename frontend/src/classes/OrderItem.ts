import axios from 'axios';
import zitadelAuth from "@/services/zitadelAuth";


export class MaterialEntry {
	id:               string;
	purchase_quantity: number;
	purchase_price:    number;
	quantity:         number
	company:          string;
    cost: number
	sku:              string;
    label: string;

    constructor(){
        this.id = ""
        this.purchase_quantity = 0
        this.purchase_price = 0
        this.quantity = 0
        this.company = ""
        this.sku = ""
        this.cost = 0
        this.label = ""
    }
}


export class MaterialSettings {
    stock_alert_treshold: number
    constructor(){
        this.stock_alert_treshold = 0
    }
}

export class Material {
	id:               string;
	name:             string;
	entries: MaterialEntry[]
	quantity:         number;
    unit: string;
    label: string;
    settings: MaterialSettings;

    constructor(){
        this.id = ""
        this.name = ""
        this.quantity = 0
        this.entries = []
        this.unit = ""
        this.label = ""
        this.settings = {
            stock_alert_treshold:0
        }
    }
}


export class ProductEntry {
    id:               string
	purchase_quantity: number
	purchase_price:    number
	quantity:         number
	company:          string
	unit:             string
	sku:              string

    constructor(){
        this.id = ""
        this.purchase_quantity = 0
        this.purchase_price = 0
        this.quantity = 0
        this.company = ""
        this.unit = ""
        this.sku = ""
    }
}

export class Product {
    id:               string;
	name:             string     ;
	materials:        Material[];
	sub_products:      Product[];
	entries: ProductEntry[];
	price:            number ;
	image_url:         string  ;
	measuring_unit:    string  ;
	quantity:         number ;
    ready: number;

    constructor(){
        this.id = ""
        this.name = ""
        this.materials = []
        this.sub_products = []
        this.entries = []
        this.price = 0
        this.image_url = ""
        this.measuring_unit = ""
        this.quantity = 0
        this.ready = 0
    }
}


export class OrderItemMaterial {
	material:       Material;
	entry: MaterialEntry; 
	quantity:       number;
    entries: MaterialEntry[]
    isQuantityValid: boolean = true;
    
    constructor(material ?:Material){

        if (material != undefined){
            this.material = material
            this.entry = material.entries[0]
            this.entries = material.entries
        }else {
            this.material = new Material()
            this.entry = new MaterialEntry()
            this.entries = []
        }

        this.quantity = 0

    }
}

export class OrderItem {

    Id: string | undefined;
    product:            Product;
	materials:          OrderItemMaterial[];
    ready: number;
    price: number;
	is_consume_from_ready: boolean;
    can_change_ready_toggle: boolean;
	sub_items:           OrderItem[];
	quantity:           number;
	comment:            string;
    isValid: boolean;


    constructor(product?: Product){

        if (product != undefined)
        {

            this.product = product
            this.quantity = product.quantity
            this.price = product.price
            this.isValid = true;
            this.ready = product.ready

            this.materials = product.materials.map( (material,index) => {
                
                material.entries.forEach((entry :MaterialEntry, entry_index) => {
                    entry.label = entry.company + " - " + entry.quantity + " " + material.unit
                    product.materials[index].entries[entry_index].label = entry.label
                    this.product.materials[index].entries[entry_index].label = entry.label
                })

                const itemmaterial = <OrderItemMaterial>{
                    entry: material.entries[0],
                    entries: material.entries,
                    material: material,
                    quantity: material.quantity
                }


                this.ValidateMaterialQuantity(index)

                return itemmaterial

            })


            if (this.product.sub_products != undefined){
                this.sub_items = product.sub_products.map((p) => {
                    const new_sub_product :OrderItem = new OrderItem(p)
                    new_sub_product.quantity = p.quantity
                    return new_sub_product
    
                })
            }else {
                this.sub_items = []
            }

            this.comment = ""
           

            if (product.quantity != undefined){
                this.quantity = product.quantity
                
                if (this.ready >= product.quantity ){

                    // allow user to edit the ready toggle and add enable the toggle so that it consumes from the ready quantity
                    this.is_consume_from_ready = false
                    this.can_change_ready_toggle = true

                }else {
                    this.can_change_ready_toggle = false
                    this.is_consume_from_ready = false
                }
            }else {
                this.can_change_ready_toggle = false
                this.is_consume_from_ready = false
                this.quantity = 1
            }

        } else {
            this.product = new Product ()
            this.materials = []
            this.comment = ""
            this.Id = ""
            this.is_consume_from_ready = false;
            this.ready = 0;
            this.sub_items = []
            this.quantity = 1
            this.can_change_ready_toggle = false
            this.price = 0
            this.isValid = true
        }
    }


    async FromItemData(orderItem: OrderItem){

        this.Id = orderItem.product.id
        await this.RefreshReadyNumber()


        this.product = orderItem.product
        this.materials = orderItem.materials

        this.materials.forEach((material: OrderItemMaterial,index: number) => {
            this.ValidateMaterialQuantity(index)
        })


        this.is_consume_from_ready = orderItem.is_consume_from_ready
        // this.can_change_ready_toggle = orderItem.can_change_ready_toggle


        for (let i=0;i<orderItem.sub_items.length;i++){
            const new_sub_item = new OrderItem()
            await new_sub_item.FromItemData(orderItem.sub_items[i])
            this.sub_items.push(new_sub_item)
        }


        this.comment = orderItem.comment
        this.quantity = orderItem.quantity
        this.price = orderItem.price ? orderItem.price : orderItem.product.price

        
        if (this.quantity != undefined){
            
            if (this.ready >= this.quantity ){

                // allow user to edit the ready toggle and add enable the toggle so that it consumes from the ready quantity
                this.can_change_ready_toggle = true

            }else {
                this.can_change_ready_toggle = false
                this.is_consume_from_ready = false
            }
        }else {
            this.can_change_ready_toggle = false
            this.is_consume_from_ready = false
            this.quantity = 1
        }
    }


    SetProductId(id: string){
        this.Id = id
        this.product.id = id
    }

    async RefreshReadyNumber(){

        await axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/products/${this.Id}`,
            {
                headers: {
                    'Authorization': `Bearer ${zitadelAuth.oidcAuth.accessToken}`
                }
            }
        ).then((response) => { 

            this.ready = response.data.data.ready

         })
    }


    async RefreshProductData(){
        await axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/products/${this.product.id}/recipetree`,
            {
                headers: {
                    'Authorization': `Bearer ${zitadelAuth.oidcAuth.accessToken}`
                }
            }
        ).then((response) => {

            const materials : Material[] = []
            this.product = response.data.data

            response.data.data.materials.forEach((material: any) => {

                const new_material = new Material()
                new_material.id = material.id
                new_material.quantity = material.quantity
                new_material.name = material.name
                new_material.unit = material.unit

                material.entries.forEach((entry: any) => {

                    const new_entry: MaterialEntry = new MaterialEntry()
                    new_entry.id = entry.id
                    new_entry.company = entry.company
                    new_entry.label = entry.company + " - " + entry.quantity + " " + material.unit
                    new_entry.purchase_price = entry.purchase_price
                    new_entry.quantity = entry.quantity
                    new_entry.purchase_quantity = entry.purchase_quantity
                    new_material.entries.push(new_entry)
                })
                materials.push(new_material)
                
            })
            this.product.materials = materials
            this.materials.forEach((material,materialIndex) => {
                this.product.materials.forEach((product_material) => {
                    if (material.material.id == product_material.id){

                        product_material.entries.forEach(pm => {
                            if (material.entry.id == pm.id){
                                material.entry = pm
                            }
                        })
                        material.material = product_material
                        material.entries = product_material.entries
                    }
                })

                this.ValidateMaterialQuantity(materialIndex)

            })
        })


        for (let i=0;i<this.sub_items.length;i++){
            const new_item = new OrderItem()
            await new_item.FromItemData(this.sub_items[i])
            await new_item.RefreshProductData();

            this.sub_items[i] = new_item   
        }
    }


    async UpdateMaterialEntryCost(materialIndex: number){
        if (this.materials[materialIndex].entry != undefined){
            await axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/materials/${this.materials[materialIndex].material.id}/entries/${this.materials[materialIndex].entry.id}/cost?quantity=${this.materials[materialIndex].quantity}`, {
                headers: {
                    'Authorization': `Bearer ${zitadelAuth.oidcAuth.accessToken}`
                }
            }).then((response) => {
    
                this.materials[materialIndex].entry.cost = response.data.data
                this.ValidateMaterialQuantity(materialIndex)
              
            })
        }
    }

    ValidateItem(){
        let valid = true


        if (this.is_consume_from_ready){
            if (this.ready >= this.quantity){
                this.isValid = valid
                return
            }
        }


        this.materials.forEach((material) => {
            if (!material.isQuantityValid)
                valid = false
        })


        if (!this.ValidateSubItems())
            valid = false

        
        this.isValid = valid
    }

    ValidateSubItems(): boolean {

        let valid = true

        this.sub_items.forEach(item => {
            if (!item.isValid){
                valid = false
            }
                
        })

        return valid
    }

    ValidateAllMaterials(){
        for (let i=0;i<this.materials.length;i++){
            this.ValidateMaterialQuantity(i)
        }

        for (let i=0;i<this.sub_items.length;i++){
            this.sub_items[i].ValidateAllMaterials()
        }
    }

    ValidateMaterialQuantity(materialIndex: number){

        if (this.materials == undefined) {
            this.isValid = false
            return
        }

        let material_entry_exists_in_product = false

        this.product.materials.forEach((material) => {
            if (material.id == this.materials[materialIndex].material.id){
                material.entries.forEach((entry) => {
                    if (entry.id == this.materials[materialIndex].entry.id){
                        material_entry_exists_in_product = true
                    }
                })
            }
        })

        if (!material_entry_exists_in_product){
            this.materials[materialIndex].isQuantityValid = false
            this.ValidateItem()
            return
        }


        if (this.materials[materialIndex].entry == undefined){
            this.materials[materialIndex].isQuantityValid = false
            this.ValidateItem()
            return
        }

        if (this.materials[materialIndex].quantity > this.materials[materialIndex].entry.quantity){
            this.materials[materialIndex].isQuantityValid = false
        }else {
            this.materials[materialIndex].isQuantityValid = true
        }

        this.ValidateItem()
    }

    async PushMaterial(material: Material) {

        material.entries.forEach(e => {
            e.label = e.company + " - " + e.quantity + " " + material.unit
        })

        const new_material = new OrderItemMaterial(material)
        this.materials.push(new_material)
        this.ValidateMaterialQuantity(this.materials.length - 1)
        await this.UpdateMaterialEntryCost(this.materials.length - 1)
    }


    RemoveMaterialByIndex(materialIndex: number){
        this.materials.splice(materialIndex,1)
        this.ValidateItem()
    }


    FillSubitems(): OrderItem[]{

        const items: OrderItem[] = []

        this.product.sub_products?.forEach((sub_product) => {
            const sub_item = new OrderItem(sub_product)
            sub_item.quantity = sub_product.quantity
            items.push(sub_item)
        })

        return items;

    }


    async ReloadDefaults() {
        await axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/products/${this.product.id}/recipetree`,
        {
            headers: {
                'Authorization': `Bearer ${zitadelAuth.oidcAuth.accessToken}`
            }
        }).then((response) => {

            const materials : Material[] = []
            const subrecipes: OrderItem[] = []
            this.product = response.data.data
            this.price = response.data.data.price
            this.ready = response.data.data.ready
            this.sub_items = this.FillSubitems()
            this.Id = response.data.data.id

            if (this.ready >= this.quantity){
                this.can_change_ready_toggle = true
            }

            response.data.data.materials.forEach((material: any) => {

                const new_material = new Material()
                new_material.id = material.id
                new_material.quantity = material.quantity
                new_material.name = material.name
                new_material.unit = material.unit

                material.entries.forEach((entry: any) => {

                    const new_entry: MaterialEntry = new MaterialEntry()
                    new_entry.id = entry.id
                    new_entry.company = entry.company
                    new_entry.label = entry.company + " - " + entry.quantity + " " + material.unit
                    new_entry.purchase_price = entry.purchase_price
                    new_entry.quantity = entry.quantity
                    new_entry.purchase_quantity = entry.purchase_quantity
                    new_material.entries.push(new_entry)
                })
                
                materials.push(new_material)
            })

            // response.data.data.sub_recipes?.forEach((sub_product: Product) => {

            //     const sub_order_item = new OrderItem()
            //     sub_order_item.Product.Id = sub_product.Id
            //     sub_order_item.Quantity = sub_product.Quantity
            //     sub_order_item.ReloadDefaults()
            //     subrecipes.push(sub_order_item)
            // })

            this.product.materials = materials
            this.materials = materials.map( material => {

                return <OrderItemMaterial>{
                    entry: material.entries[0],
                    entries: material.entries,
                    material: material,
                    quantity: material.quantity
                }

            })

            this.materials.forEach((material,materialIndex) => {
                this.ValidateMaterialQuantity(materialIndex)
            })
        })  
    }

}