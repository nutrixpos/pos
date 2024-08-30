export interface ComponentEntry {
    _id:              string | undefined,
    purchaseQuantity: number,
    quantity         :number,
    company          :string,
    unit             :string,
    price            :number,
    label: string
}

export interface Component {
    component_id: string,
    name: string,
    unit: string,
    defaultquantity: number,
    entries: Array<ComponentEntry>,
}

export interface Recipe {
    recipe_name: string,
    recipe_id: string,
    ready: number,
    default_quantity: number,
    components: Array<Component>,
    sub_recipes: Array<Recipe>,
}


export interface ComponentEntrySelection {
    name: string,
    ComponentId: string | undefined,
    EntryId: string | undefined,
    unit: string,
    quantity: number,
}

export class RecipeSelections {
    recipe_name?: string;
    Id?: string | undefined;
    Ready: number;
    Components: Array<Component>; 
    isConsumeFromReady: boolean;
    canChangeReadyToggle: boolean;
    Selections?: Array<ComponentEntrySelection>
    SubRecipes: Array<RecipeSelections>;
    Quantity: string | null;


    constructor(recipe: Recipe){

        console.log(recipe)

        this.Ready = recipe.ready
        this.Components = recipe.components
        this.isConsumeFromReady = false
        this.canChangeReadyToggle = false
        this.SubRecipes = []
        this.Selections = []
        this.Quantity = null
        this.recipe_name = recipe.recipe_name
        this.Id = recipe.recipe_id

        if (recipe.default_quantity != undefined){
            this.Quantity = recipe.default_quantity.toString()
            
            if (this.Ready >= recipe.default_quantity ){
                this.isConsumeFromReady = true
                this.canChangeReadyToggle = true

            }else {
                this.canChangeReadyToggle = false
                this.isConsumeFromReady = false
            }
        }else {
            this.canChangeReadyToggle = false
            this.isConsumeFromReady = false
            this.Quantity = ""
        }


        if (recipe.sub_recipes != null)
            recipe.sub_recipes?.forEach(subrecipe => {
            
                const newSubRecipe = new RecipeSelections(subrecipe)
                this.SubRecipes.push(newSubRecipe)
            })
    }

}