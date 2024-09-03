export class ComponentEntry {
    _id:              string;
    quantity         :number;
    company          :string;
    unit             :string;
    price            :number;
    label: string;
    PurchaseQuantity: number;

    constructor(){
        this._id = ""
        this.quantity = 0
        this.company = ""
        this.unit = ""
        this.price = 0
        this.label = ""
        this.PurchaseQuantity = 0
    }
}

export class Component {
    component_id: string;
    name: string;
    unit: string;
    defaultquantity: number;
    entries: Array<ComponentEntry>;
    

    constructor(){
        this.component_id = ""
        this.name = ""
        this.unit = ""
        this.defaultquantity = 0
        this.entries = []
    }
}

export interface Recipe {
    recipe_name: string,
    recipe_id: string,
    ready: number,
    quantity: number,
    components: Array<Component>,
    sub_recipes: Array<Recipe>,
}


export class ComponentSelection {
    Name: string;
    ComponentId: string;
    Entry: ComponentEntry;
    Unit: string;
    Quantity: number;

    constructor(){
        this.Name = ""
        this.ComponentId = ""
        this.Unit = ""
        this.Quantity = 0
        this.Entry = new ComponentEntry()
    }
}

export class RecipeSelections {
    recipe_name: string;
    Id: string | undefined;
    Ready: number;
    Components: Array<Component>; 
    isConsumeFromReady: boolean;
    canChangeReadyToggle: boolean;
    Selections: Array<ComponentSelection>
    SubRecipes: Array<RecipeSelections>;
    Quantity: number;


    constructor(recipe?: Recipe){

        if (recipe != undefined)
        {

            this.Ready = recipe.ready
            this.Components = []
            this.isConsumeFromReady = false
            this.canChangeReadyToggle = false
            this.SubRecipes = []
            this.Selections = []
            this.Quantity = recipe.quantity
            this.recipe_name = recipe.recipe_name
            this.Id = recipe.recipe_id

            recipe.components.forEach((component) => {

                const new_component = new Component()
                new_component.component_id = component.component_id
                new_component.defaultquantity = component.defaultquantity
                new_component.name = component.name
                new_component.unit = component.unit

                component.entries.forEach(entry => {

                    const new_entry: ComponentEntry = new ComponentEntry()
                    new_entry._id = entry._id
                    new_entry.company = entry.company
                    new_entry.label = entry.company + " - " + entry.quantity + " " + entry.unit
                    new_entry.price = entry.price
                    new_entry.quantity = entry.quantity
                    new_entry.unit = entry.unit
                    new_entry.PurchaseQuantity = entry.PurchaseQuantity
                    new_component.entries.push(new_entry)
                })
                
                this.Components.push(new_component)
            })

            if (recipe.quantity != undefined){
                this.Quantity = recipe.quantity
                
                if (this.Ready >= recipe.quantity ){
                    this.isConsumeFromReady = true
                    this.canChangeReadyToggle = true

                }else {
                    this.canChangeReadyToggle = false
                    this.isConsumeFromReady = false
                }
            }else {
                this.canChangeReadyToggle = false
                this.isConsumeFromReady = false
                this.Quantity = 0
            }


            if (recipe.sub_recipes != null)
                recipe.sub_recipes?.forEach(subrecipe => {
                
                    const newSubRecipe = new RecipeSelections(subrecipe)
                    this.SubRecipes.push(newSubRecipe)
                })
        } else {
            this.recipe_name = ""
            this.Id = ""
            this.Ready = 0
            this.Components = []
            this.isConsumeFromReady = false
            this.canChangeReadyToggle = false
            this.SubRecipes = []
            this.Selections = []
            this.Quantity =1
        }
    }


}