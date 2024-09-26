import axios from 'axios';


export class MaterialEntry {
	_id:               string;
	purchase_quantity: number;
	purchase_price:    number;
	quantity:         number
	company:          string;
    cost: number
	sku:              string;
    label: string;

    constructor(){
        this._id = ""
        this.purchase_quantity = 0
        this.purchase_price = 0
        this.quantity = 0
        this.company = ""
        this.sku = ""
        this.cost = 0
        this.label = ""
    }
}


export class Material {
	_id:               string;
	name:             string;
	entries: MaterialEntry[]
	quantity:         number;
    unit: string;
    label: string;

    constructor(){
        this._id = ""
        this.name = ""
        this.quantity = 0
        this.entries = []
        this.unit = ""
        this.label = ""
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

                return itemmaterial

            })

            this.ready = product.ready

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
                    this.can_change_ready_toggle = false

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


    FromItemData(orderItem: OrderItem){

        this.Id = orderItem.Id
        this.product = orderItem.product
        this.materials = orderItem.materials
        this.ready = orderItem.ready
        this.is_consume_from_ready = orderItem.is_consume_from_ready
        this.can_change_ready_toggle = orderItem.can_change_ready_toggle
        this.sub_items = orderItem.sub_items
        this.comment = orderItem.comment
        this.quantity = orderItem.quantity
        this.price = orderItem.price
    }


    async RefreshProductData(){
        await axios.get("http://localhost:8000/api/recipetree?id="+this.product.id).then((response) => {

            const materials : Material[] = []
            this.product = response.data

            response.data.materials.forEach((material: any) => {

                const new_material = new Material()
                new_material._id = material._id
                new_material.quantity = material.quantity
                new_material.name = material.name
                new_material.unit = material.unit

                material.entries.forEach((entry: any) => {

                    const new_entry: MaterialEntry = new MaterialEntry()
                    new_entry._id = entry._id
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
            this.materials.forEach((material) => {
                this.product.materials.forEach((product_material) => {
                    if (material.material._id == product_material._id){

                        product_material.entries.forEach(pm => {
                            if (material.entry._id == pm._id){
                                material.entry = pm
                            }
                        })
                        material.material = product_material
                        material.entries = product_material.entries
                    }
                })
            })
        })
        
        this.sub_items.forEach((sub_item,index) => {
            const new_item = new OrderItem()
            new_item.FromItemData(sub_item)
            new_item.RefreshProductData();

            this.sub_items[index] = new_item
        })
    }


    async UpdateMaterialEntryCost(materialIndex: number){
        await axios.get(`http://localhost:8000/api/materialcost?material_id=${this.materials[materialIndex].material._id}&entry_id=${this.materials[materialIndex].entry._id}&quantity=${this.materials[materialIndex].quantity}`).then((response) => {

            this.materials[materialIndex].entry.cost = response.data
          
        })  
    }

    PushMaterial(material: Material) {

        material.entries.forEach(e => {
            e.label = e.company + " - " + e.quantity + " " + material.unit
        })

        const new_material = new OrderItemMaterial(material)
        this.materials.push(new_material)
        this.UpdateMaterialEntryCost(this.materials.length - 1)
    }


    RemoveMaterialByIndex(materialIndex: number){
        this.materials.splice(materialIndex,1)   
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
        await axios.get("http://localhost:8000/api/recipetree?id="+this.product.id).then((response) => {

            const materials : Material[] = []
            const subrecipes: OrderItem[] = []
            this.product = response.data
            this.price = response.data.price
            this.sub_items = this.FillSubitems()

            response.data.materials.forEach((material: any) => {

                const new_material = new Material()
                new_material._id = material._id
                new_material.quantity = material.quantity
                new_material.name = material.name
                new_material.unit = material.unit

                material.entries.forEach((entry: any) => {

                    const new_entry: MaterialEntry = new MaterialEntry()
                    new_entry._id = entry._id
                    new_entry.company = entry.company
                    new_entry.label = entry.company + " - " + entry.quantity + " " + material.unit
                    new_entry.purchase_price = entry.purchase_price
                    new_entry.quantity = entry.quantity
                    new_entry.purchase_quantity = entry.purchase_quantity
                    new_material.entries.push(new_entry)
                })
                
                materials.push(new_material)
            })

            // response.data.sub_recipes?.forEach((sub_product: Product) => {

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
        })  
    }

}

// export class OrderItem {

//     recipe_name: string;
//     Id: string | undefined;
//     Ready: number;
//     isConsumeFromReady: boolean;
//     canChangeReadyToggle: boolean;
//     Materials: Array<Material>
//     SubRecipes: Array<OrderItem>;
//     Quantity: number;
//     Comment: string;
//     Design: RecipeDesign


//     constructor(recipe?: RecipeDesign){

//         if (recipe != undefined)
//         {

//             this.Design = recipe
//             this.Ready = recipe.ready
//             this.isConsumeFromReady = false
//             this.canChangeReadyToggle = false
//             this.SubRecipes = []
//             this.materials = []
//             this.Quantity = recipe.quantity
//             this.recipe_name = recipe.recipe_name
//             this.Id = recipe.recipe_id
//             this.Comment = recipe.comment

//             // recipe.components.forEach((component) => {

//             //     const new_component = new Component()
//             //     new_component.component_id = component.component_id
//             //     new_component.defaultquantity = component.defaultquantity
//             //     new_component.name = component.name
//             //     new_component.unit = component.unit

//             //     component.entries.forEach(entry => {

//             //         const new_entry: ComponentEntry = new ComponentEntry()
//             //         new_entry._id = entry._id
//             //         new_entry.company = entry.company
//             //         new_entry.label = entry.company + " - " + entry.quantity + " " + entry.unit
//             //         new_entry.price = entry.price
//             //         new_entry.quantity = entry.quantity
//             //         new_entry.unit = entry.unit
//             //         new_entry.PurchaseQuantity = entry.PurchaseQuantity
//             //         new_component.entries.push(new_entry)
//             //     })
                
//             //     this.Design.components.push(new_component)
//             // })

//             if (recipe.quantity != undefined){
//                 this.Quantity = recipe.quantity
                
//                 if (this.Ready >= recipe.quantity ){

//                     // allow user to edit the ready toggle and add enable the toggle so that it consumes from the ready quantity
//                     this.isConsumeFromReady = false
//                     this.canChangeReadyToggle = false

//                 }else {
//                     this.canChangeReadyToggle = false
//                     this.isConsumeFromReady = false
//                 }
//             }else {
//                 this.canChangeReadyToggle = false
//                 this.isConsumeFromReady = false
//                 this.Quantity = 0
//             }


//             if (recipe.sub_recipes != null)
//                 recipe.sub_recipes?.forEach(subrecipe => {
                
//                     const newSubRecipe = new OrderItem(subrecipe)
//                     this.SubRecipes.push(newSubRecipe)
//                 })
//         } else {
//             this.recipe_name = ""
//             this.Id = ""
//             this.Ready = 0
//             this.isConsumeFromReady = false
//             this.canChangeReadyToggle = false
//             this.SubRecipes = []
//             this.materials = []
//             this.Quantity =1
//             this.Comment = ""
//             this.Design = new RecipeDesign()
//             this.Design.components = []
//         }
//     }


//     ReloadDefaults() {
//         axios.get("http://localhost:8000/api/recipetree?id="+this.Id,).then((response) => {

//             const components : Component[] = []
//             const subrecipes: OrderItem[] = []
//             this.Id = response.data.recipe_id
//             this.recipe_name = response.data.recipe_name
//             this.Design = response.data

//             response.data.components.forEach((component: Component) => {

//                 const new_component = new Component()
//                 new_component.component_id = component.component_id
//                 new_component.defaultquantity = component.defaultquantity
//                 new_component.name = component.name
//                 new_component.unit = component.unit

//                 component.entries.forEach((entry:ComponentEntry) => {

//                     const new_entry: ComponentEntry = new ComponentEntry()
//                     new_entry._id = entry._id
//                     new_entry.company = entry.company
//                     new_entry.label = entry.company + " - " + entry.quantity + " " + entry.unit
//                     new_entry.price = entry.price
//                     new_entry.quantity = entry.quantity
//                     new_entry.unit = entry.unit
//                     new_entry.PurchaseQuantity = entry.PurchaseQuantity
//                     new_component.entries.push(new_entry)
//                 })
                
//                 components.push(new_component)
//             })

//             response.data.sub_recipes?.forEach((subrecipe: RecipeDesign) => {

//                 const new_subrecipe = new OrderItem(subrecipe)
//                 subrecipes.push(new_subrecipe)
//             })

//             this.Design.components = components
//             this.SubRecipes = subrecipes
//         })  
//     }

// }