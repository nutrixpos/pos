<template>
    <div v-if="!loading" class="flex flex-column m-0 p-0" style="height: 100vh;">
        <Toolbar style="border-radius: 0px;flex-shrink: 0;" class="py-1 lg:py-2">
            <template #start>
                <div @click="version_dialog_visible=true" style="text-decoration: none;color:gray">
                    <img src="@/assets/logo.png" alt="logo" style="height:25px" v-if="store.getColorMode == 'light'">
                    <div v-else class="flex justify-content-center align-items-center" style="font-size:1rem;color:white;font-family:'FontAwesome'">
                        nutrix
                    </div>
                </div>
                <div class="flex mx-2 gap-2">
                    <div  v-for="(item,index) in navbar_links" :key="index" >
                        <router-link v-if="item.authority.some(role => roles.includes(role)) || !zitadel_enabled" :to="item.link" >
                            <Button :icon="item.icon" :label="t(`${item.label.title}`,item.label.plural ? 3 : 1)"  :text="!item.focused" severity="secondary" />
                        </router-link>
                    </div>
                </div>
            </template>

            <template #center>
                <IconField iconPosition="left">
                    <InputIcon>
                        <i class="pi pi-search" />
                    </InputIcon>
                    <InputText :placeholder="t('search')" v-model="mainSearchText" @keyup.stop="mainSearchTextChanged" />
                </IconField>
                <OverlayPanel ref="mainsearch_op" class="w-5 lg:w-3" style="max-height:60vh;overflow-y: auto;">
                    <MainSearchResultView class="mt-2" @view-order-pressed="order_to_show = result; order_details_dialog=true" v-for="(result,index) in mainSearchResult" :key="index" :order="result" />
                </OverlayPanel>
            </template>

            <template #end>
                <Button  severity="secondary" size="large"  text rounded aria-label="Current" @click.stop="paylater_toggle">
                    <span class="p-button-icon pi pi-hourglass"></span>
                    <Badge :value="inProgressOrders.length" class="p-badge-warn"  />
                </Button>
                <OverlayPanel ref="current_orders_op" class="w-8 xl:w-3" style="max-height:60vh;overflow-y: auto;">
                    <h4 class="m-2" style="color:#c2c2c2">{{t('current_orders')}}</h4>
                    <MainSearchResultView class="mt-2" @view-order-pressed="order_to_show = result; order_details_dialog=true" v-for="(result,index) in inProgressOrders" :key="index" :order="result" />
                </OverlayPanel>
                <Button outlined :icon="`pi pi-${store.getColorMode == 'light' ? 'sun' : 'moon'}`" @click="toggleDarkMode()" />
                <Button icon="pi pi-bars" severity="secondary" size="large" text rounded aria-label="drawer" @click="drawer_visible=true" />
            </template>
        </Toolbar>
        <div class="grid m-0 p-0" style="height:calc(94vh - 1.5rem);flex-shrink: 0;">
            <div class="col-2 flex flex-column py-3">
                <div class="w- flex my-1" v-for="(category,index) in categories" :key="index" :style="`background-color:${store.getColorMode == 'light' ? 'white' : '#18181B' };cursor:pointer`" @click="selectedCategory = category">
                    <div :style="`width:0.5rem;background-color:${category == selectedCategory? '#FDDB00' : 'silver'}`"></div>
                    <div class="py-3 mx-3" :style="`color:${category == selectedCategory? store.getColorMode == 'light' ? '#2E4762': '#FDDB00' : store.getColorMode == 'light' ? 'black' : 'white'};font-weight:${category == selectedCategory? 'bold' : '200'}`">{{ category.name }}</div>
                </div>
            </div>
            <div class="xl:col-7 col-5 flex px-0 xl:px-2 overflow-auto">
                <Card style="width:100%;">
                    <template #content>
                        <IconField iconPosition="left" class="mb-3">
                            <InputIcon v-if="!isSearchingProduct">
                                <i class="pi pi-search" />
                            </InputIcon>
                            <InputIcon v-if="isSearchingProduct" class="pi pi-spin pi-spinner"> </InputIcon>
                            <InputText v-model="searchtext" :placeholder="t('search')+' '+t('product',3)" />
                        </IconField>
                        <div class="flex flex-wrap">
                            <MealCard  v-for="(item,index) in products" :key="index" :item="item" class="m-1 lg:m-2" style="width:9rem;" @addwithcomment="visible=true;idwithcomment=item.id; namewithcomment=item.name" @add="addItem(item)"/>
                        </div>
                    </template>
                </Card>
            </div>
            <div class="col-5 xl:col-3 flex">
                <Card class="w-12" :style="`border-color: ${is_order_valid ?  '' : 'red'};`">
                    <template #content>
                        <div class="flex flex-column" style="height:calc(94vh - 1.5rem - 36px - var(--p-card-body-padding)); flex-shrink: 0">
                            <div style="height:65%;overflow: auto;">
                                <div v-for="(item,index) in orderItems" :key="index">
                                    <div class="flex justify-content-between align-items-center">
                                        <div style="background-color:red;width:0.3rem;height:2.5rem;" class="mr-2" v-if="!item.isValid">
                                            &nbsp;
                                        </div>
                                        <p class="w-6" style="text-overflow:ellipsis"><strong>{{ item.product.name }}</strong></p>
                                        <p>{{ item.price }} {{t('egp')}}</p>
                                        <div>
                                            <Button icon="pi pi-pencil" size="small" style="width:2rem;height: 2rem;" aria-label="Edit" severity="secondary" @click="itemToEditIndex = index; edit_item_dialog=true" class="mr-1"/>
                                            <Button icon="pi pi-times" size="small" style="width:2rem;height: 2rem;" aria-label="Remove" severity="secondary" @click="orderItems.splice(index,1)" />
                                        </div>
                                    </div>
                                    <p class="m-0">{{ item.comment }}</p>
                                </div>
                            </div>
                            <div class="flex flex-column flex-wrap justify-content-between">
                                <div>
                                    <Divider />
                                    <ButtonGroup class="flex justify-content-start flex-wrap gap-1">
                                        <Button icon="pi pi-bookmark" label="Draft" @click="stashOrder" size="small" severity="secondary" v-tooltip.top="'Draft order for later interactions'" aria-label="Stash order" />
                                        <Button label="Add Discount" severity="secondary" icon="fa fa-percent" size="small"  @click="toggle_discount_popover" />
                                    </ButtonGroup>
                                    <Popover ref="discount_op">
                                        <div class="flex flex-column gap-2" style="width:25vw">
                                            <p class="mb-0">Add discount</p>
                                            <div class="flex w-full">
                                                <Slider v-model="discount_percent" class="w-9 m-1 mx-2" style="height:0.6rem;" />
                                                <p class="ml-2" style="font-size:0.8rem">{{ discount_percent.toFixed(2) }} %</p>
                                            </div>
                                            <div class="w-7 flex">
                                                <InputText v-model.number="discount" :disabled="orderItems.length == 0" placeholder="0" type="number" class="w-8 h-2rem"  />
                                                <p style="font-size:0.8rem" class="mx-2">{{ t('egp') }}</p>
                                            </div>
                                        </div>
                                    </Popover>
                                    <div class="flex justify-content-between flex-wrap align-items-center">
                                        <p style="font-size:0.8rem;" class="mb-0">{{t('subtotal')}} : </p>
                                        <p style="font-size:0.7rem;">{{ subtotal.toFixed(2) }} <span style="font-size:0.7rem">{{t('egp')}}</span></p>
                                    </div>
                                    <div class="flex justify-content-between flex-wrap align-items-center">
                                        <p class="my-0" style="font-size:0.8rem">{{t('discount')}} :</p>
                                        <p style="font-size:0.7rem;">{{ discount.toFixed(2) }} <span style="font-size:0.7rem">{{t('egp')}}</span></p>
                                    </div>
                                    <div class="flex justify-content-between flex-wrap align-items-center">
                                        <h2>{{t('total')}} : </h2>
                                        <p style="font-size:1.4rem">{{ total.toFixed(2) }} <span style="font-size:1rem">{{t('egp')}}</span></p>
                                    </div>
                                </div>
                                <div class="flex flex-column align-items-start mb-3">
                                    <div class="flex justify-content-center align-items-center mt-2 flex-wrap gap-2">
                                        <ToggleButton size="small" v-model="is_print_receipt_kitchen" onLabel="Kitchen" offLabel="Kitchen" onIcon="fa fa-print" offIcon="fa fa-print" class="w-36" aria-label="Do you confirm" />
                                        <ToggleButton size="small" v-model="is_print_receipt_client" onLabel="Client" offLabel="Client" onIcon="fa fa-print" offIcon="fa fa-print" class="w-36 mx-1" aria-label="Do you confirm" />
                                        <!-- <ToggleButton size="small" v-tooltip.top="'Auto start order and consume components from inventory'" v-model="is_auto_start_order" onLabel="Autostarting" offLabel="Autostart" onIcon="pi pi-check" offIcon="pi pi-play-circle" class="w-36 mx-1" aria-label="Do you confirm" /> -->
                                        <ToggleButton size="small" v-tooltip.top="'Auto finish order and consume components from inventory'" v-model="is_auto_finish_order" onLabel="Autofinishing" offLabel="Autofinish" onIcon="pi pi-check" offIcon="pi pi-play-circle" class="w-36 mx-1" aria-label="Do you confirm" />
                                    </div>
                                </div>
                                <Button label="Next" icon="pi pi-arrow-right" iconPos="right" :disabled="!is_order_valid || orderItems.length == 0" @click="order_additional_details_dialog=true" />
                            </div>
                        </div>
                    </template>
                </Card>
            </div>
            <Dialog v-model:visible="edit_item_dialog" modal header="Edit item" class="xs:w-12 md:w-10 lg:w-8">
                <OrderItemView v-model="orderItems[itemToEditIndex]"  />
            </Dialog>
            <Dialog v-model:visible="order_details_dialog" modal header="Order details" class="w-11 xl:w-8">
                <OrderView @updated="order_details_dialog=false" @finished="finishOrderDisplayed()" @cancelled="cancelOrderDisplayed()" @amount_collected="orderToShowAmountCollected()" :order="order_to_show" />
            </Dialog>
            <Dialog v-model:visible="visible" modal header="Add Comment" :style="{ width: '25rem' }">
                <InputText v-model="comment" placeholder="Comment" class="mb-4" />
                <div class="flex justify-content-end gap-2">
                    <Button type="button" label="Close" severity="secondary"></Button>
                    <Button type="button" label="Add" @click="addWithComment()"></Button>
                </div>
            </Dialog>
            <Dialog v-model:visible="order_additional_details_dialog" modal header="Order details" class="xs:w-12 md:w-10 lg:w-8">
                <Stepper linear value="1">
                    <StepList>
                        <Step v-for="(step,index) in order_details_steps" :key="index" :value="`${index+1}`" >{{ step.label }}</Step>
                    </StepList>
                    <StepPanels>
                        <StepPanel v-slot="{ activateCallback }" value="1">
                            <div class="flex justify-content-center flex-wrap gap-3">
                                <div class="flex flex-column align-items-center">
                                    <h2 class="mt-0">
                                        <i class="fa-regular fa-credit-card mx-2"></i>
                                         Payment
                                    </h2>
                                    <ToggleButton v-model="is_collecting_money" onIcon="fa fa-hand-holding-dollar" offIcon="fa fa-hand-holding-dollar" :offLabel="`Collect (${total.toFixed(2)} EGP)`" :onLabel="`Collecting (${(total+ ( current_order_tip || 0 )).toFixed(2)} EGP)`" class="w-15rem h-5rem lg:h-10rem sm:w-40 border-noround" aria-label="Confirmation" />
                                    <ToggleButton v-model="is_pay_later" onIcon="pi pi-hourglass" offIcon="fa fa-hourglass" offLabel="Pay later" onLabel="Paying later" class="w-15rem h-5rem lg:h-10rem sm:w-40 border-noround" aria-label="Confirmation" />
                                    <div class="flex flex-column align-items-start justify-content-start mt-4 w-full">
                                        <h4 class="mt-0">Payment source</h4>
                                        <Select v-model="payment_source" :options="payment_sources" optionLabel="name" placeholder="Payment source" />
                                    </div>

                                    <div class="flex m-0 p-0 flex-column mt-4 align-items-start justify-content-start w-full">
                                        <h4>Tips</h4>
                                        <InputText placeholder="Tip" v-model.number="current_order_tip" />
                                    </div>
                                </div>
                                <Divider layout="vertical" />
                                <div class="flex flex-column align-items-center">
                                    <h2 class="mt-0">
                                        <i class="fa fa-box-open mx-2"></i>
                                        Service Style
                                    </h2>
                                    <ToggleButton v-model="is_serve_inside" onIcon="fa fa-bowl-food" offIcon="fa fa-bowl-food" offLabel="Dine in" onLabel="Dine in" class="w-15rem h-5rem lg:h-10rem sm:w-40 border-noround" aria-label="Confirmation" />
                                    <ToggleButton v-model="is_take_away" onIcon="pi pi-box" offIcon="pi pi-box" offLabel="Takeaway" onLabel="Takeaway" class="w-15rem h-5rem lg:h-10rem sm:w-40 border-noround" aria-label="Confirmation" />
                                    <ToggleButton v-model="is_delivery" onIcon="pi pi-truck" offIcon="pi pi-truck" offLabel="Delivery" onLabel="Delivery" class="w-15rem h-5rem lg:h-10rem sm:w-40 border-noround" aria-label="Confirmation" />
                                </div>
                                <Divider layout="vertical" />
                                <div class="flex flex-column align-items-start">
                                    <h2 class="mt-0">
                                        <i class="fa-regular fa-comment mx-2"></i>
                                        Comment
                                    </h2>
                                    <Textarea v-model="order_comment" size="small" placeholder="Comment" rows="3" />
                                    <h2 class="mt-7 w-full align-items-start justify-content-start">
                                        <i class="fa fa-pen-to-square mx-2"></i>
                                        Custom
                                    </h2>
                                    <div class="flex flex-column">
                                        <div v-for="(item,index) in custom_data" :key="index" class="flex align-items-start flex-column gap-1">
                                            <span>Key:</span>
                                            <InputText v-model="custom_data[index].key" />
                                            <span>Value:</span>
                                            <InputText v-model.number="custom_data[index].value"/>
                                            <Button severity="secondary" aria-label="Remove" icon="pi pi-times" @click="custom_data.splice(index,1)" />
                                        </div>
                                        <div class="flex flex-column align-items-start mt-3 gap-1">
                                            <span>Key:</span>
                                            <InputText v-model="new_custom_data_key" />
                                            <span>Value:</span>
                                            <InputText v-model="new_custom_data_value"/>

                                            <Button class="mt-2" label="Add" @click="custom_data.push({key:new_custom_data_key,value:new_custom_data_value}); new_custom_data_key = ''; new_custom_data_value = ''" />
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div class="flex pt-6 justify-content-end">
                                <Button label="Next" icon="pi pi-arrow-right" iconPos="right" @click="activateCallback('2')" />
                            </div>
                        </StepPanel>
                        <StepPanel v-slot="{ activateCallback }" v-if="order_details_steps.length == 3" value="2">
                            <div class="flex flex-column">
                                <h2 class="mt-0">Customer</h2>
                                <div class="flex flex-wrap gap-2">
                                    <Button class="px-3" :label="t('existing_customer')" severity="secondary" icon="pi pi-user" iconPos="right" @click="pick_customer_dialog=true" />
                                    <Button class="px-3" :label="t('add_new')" severity="secondary" icon="pi pi-plus" iconPos="right" @click="add_customer_dialog=true" />
                                </div>
                                <AddCustomer @update:visible="(x) => add_customer_dialog = x" :visible="add_customer_dialog" @customer-added="add_customer_dialog=false" />
                                <DataTable showGridlines  class="w-full mt-3 px-3" :value="new_order_delivery_customer">
                                    <Column field="name" header="Name"></Column>
                                    <Column field="address" header="Address"></Column>
                                    <Column field="phone" header="Phone"></Column>
                                </DataTable>
                                <div class="flex pt-6 justify-content-end gap-2">
                                    <Button label="Back" icon="pi pi-arrow-left" iconPos="left" @click="activateCallback('1')" severity="secondary" />
                                    <Button label="Next" icon="pi pi-arrow-right" iconPos="right" @click="activateCallback('3')" />
                                </div>
                            </div>
                        </StepPanel>
                        <StepPanel  v-slot="{ activateCallback }" :value="order_details_steps.length == 3 ? '3' : '2'">
                            <div class="flex flex-column">
                                <h2 class="mt-0">
                                    <i class="pi pi-comment mx-2"></i>
                                    Recap
                                </h2>
                                <Divider />
                                <h3>Main details</h3>
                                <div class="flex align-items-start mt-3 gap-1">
                                    <span>{{t('autostarting')}}:</span>
                                    <p class="my-0"><strong> {{is_auto_start_order}} </strong></p>
                                </div>
                                <div class="flex align-items-start mt-3 gap-1">
                                    <span>Payment:</span>
                                    <p class="my-0"><strong> {{ is_pay_later ? t('pay_later') : t('now') }} </strong></p>
                                </div>
                                <div class="flex align-items-start mt-3 gap-1">
                                    <span>Payment Source:</span>
                                    <p class="my-0"><strong> {{ payment_source?.name }} </strong></p>
                                </div>
                                <div class="flex align-items-start mt-3 gap-1">
                                    <span>{{t('location')}}:</span>
                                    <p  class="my-0" v-if="is_serve_inside"><strong> {{t('dine_in')}} </strong></p>
                                    <p  class="my-0" v-if="is_delivery"><strong> {{t('delivery')}} </strong></p>
                                    <p  class="my-0" v-if="is_take_away"><strong> {{t('takeaway')}} </strong></p>
                                </div>
                                <div class="flex align-items-start mt-3 gap-1">
                                    <span>{{t('comment')}}:</span>
                                    <p class="my-0"><strong> {{order_comment}} </strong></p>
                                </div>
                                <div class="flex flex-column align-items-start mt-3 gap-1">
                                    <span>{{t('other')}}:</span>
                                    <div class="flex align-items-start mt-3 gap-1" v-for="(item,index) in custom_data" :key="index">
                                        <span>{{item.key}}:</span>
                                        <p class="my-0"><strong> {{item.value}} </strong></p>
                                    </div>
                                </div>
                                <div v-if="is_delivery">
                                    <Divider />
                                    <h3>Delivery</h3>
                                    <div class="flex flex-column">
                                        <div class="flex align-items-start mt-3 gap-1">
                                            <span>Customer:</span>
                                            <p class="my-0"><strong> {{new_order_delivery_customer[0]?.name}} </strong></p>
                                        </div>
                                        <div class="flex align-items-start mt-3 gap-1">
                                            <span>Address:</span>
                                            <p class="my-0"><strong> {{new_order_delivery_customer[0]?.address}} </strong></p>
                                        </div>
                                        <div class="flex align-items-start mt-3 gap-1">
                                            <span>Phone:</span>
                                            <p class="my-0"><strong> {{new_order_delivery_customer[0]?.phone}} </strong></p>
                                        </div>
                                    </div>
                                </div>
                                <div class="flex pt-6 justify-content-end gap-2">
                                    <Button label="Back" icon="pi pi-arrow-left" iconPos="left" @click="order_details_steps.length == 3 ? activateCallback('2') : activateCallback('1')" severity="secondary" />
                                    <Button :label="`${is_auto_start_order ? t('start') : t('submit')} ${is_collecting_money ? '( '+t('collect')+ ' '+ total.toFixed(2) + ' EGP )' : '( '+t('pay_later') + ' )'} `" :disabled="!is_order_valid" @click="submitOrder()" />
                                </div>
                            </div>
                        </StepPanel>
                    </StepPanels>
                </Stepper>
                <Dialog v-model:visible="pick_customer_dialog">
                    <PickCustomer @returnCustomer="(customer) => {new_order_delivery_customer = [customer]; pick_customer_dialog = false}" />
                </Dialog>
            </Dialog>
        </div>
    </div>
    <div style="width:100vw;height:100vh;display:flex;justify-content:center;align-items:center" v-if="loading">
      <ProgressSpinner style="width: 50px; height: 50px" strokeWidth="8" fill="transparent"
      animationDuration=".5s" aria-label="Custom ProgressSpinner" />
    </div>
    <Drawer v-model:visible="drawer_visible">
        <template #container="">
            <div class="flex flex-column h-full">
                <div class="flex gap-2 p-2 align-items-start mt-3">
                    <Avatar icon="pi pi-user" class="mr-2" size="medium" />
                    <div class="flex flex-column align-items-start justify-content-start">
                        <div class="mx-2">
                            {{ user?.name || "Anonymous" }}
                        </div>
                        <Button class="mt-1" icon="pi pi-sign-out" severity="secondary" text aria-label="Signout" :label="t('signout')" @click="proxy.$zitadel?.oidcAuth.signOut()" />
                    </div>
                </div>
                <div class="mt-2">
                    <Chip v-for="(role,index) in roles" :key="index" :label="role" style="height: 1.5rem;" class="m-1" />
                </div>
                <Divider />
                <div class="overflow-y-auto">
                    <ul class="list-none p-4 m-0">
                        <li>
                            <a
                                v-ripple
                                v-styleclass="{
                                    selector: '@next',
                                    enterFromClass: 'hidden',
                                    enterActiveClass: 'slidedown',
                                    leaveToClass: 'hidden',
                                    leaveActiveClass: 'slideup'
                                }"
                                class="flex align-items-center cursor-pointer p-4 rounded text-surface-700 hover:bg-surface-100 dark:text-surface-0 dark:hover:bg-surface-800 duration-150 transition-colors p-ripple"
                            >
                                <i class="pi pi-bell mr-2"></i>
                                <span class="font-medium">Notifications</span>
                                <i class="pi pi-chevron-down ml-auto"></i>
                            </a>
                            <ul class="list-none py-0 pl-4 pr-0 m-0 overflow-y-hidden transition-all duration-[400ms] ease-in-out">
                                <li @click.stop="notifications_toggle">
                                    <a v-ripple class="flex justify-content-center align-items-center cursor-pointer px-4 py-3 rounded text-surface-700 hover:bg-surface-100 dark:text-surface-0 dark:hover:bg-surface-800 duration-150 transition-colors p-ripple">
                                        <span class="font-medium">Success</span>
                                        <span class="inline-flex items-center justify-center ml-auto text-primary-contrast rounded-full" style="min-width: 1.5rem; height: 1.5rem">
                                            <Badge :value="notifications_severity_counter[0]" class="p-badge-success"  />
                                        </span>
                                    </a>
                                </li>
                                <li @click.stop="notifications_toggle">
                                    <a v-ripple class="flex justify-content-center align-items-center cursor-pointer px-4 py-3 rounded text-surface-700 hover:bg-surface-100 dark:text-surface-0 dark:hover:bg-surface-800 duration-150 transition-colors p-ripple">
                                        <span class="font-medium">Info</span>
                                        <span class="inline-flex items-center justify-center ml-auto text-primary-contrast rounded-full" style="min-width: 1.5rem; height: 1.5rem">
                                            <Badge :value="notifications_severity_counter[1]" class="p-badge-info"  />
                                        </span>
                                    </a>
                                </li>
                                <li @click.stop="notifications_toggle">
                                    <a v-ripple class="flex justify-content-center align-items-center cursor-pointer px-4 py-3 rounded text-surface-700 hover:bg-surface-100 dark:text-surface-0 dark:hover:bg-surface-800 duration-150 transition-colors p-ripple">
                                        <span class="font-medium">Warn</span>
                                        <span class="inline-flex items-center justify-center ml-auto text-primary-contrast rounded-full" style="min-width: 1.5rem; height: 1.5rem">
                                            <Badge :value="notifications_severity_counter[2]" class="p-badge-warn"  />
                                        </span>
                                    </a>
                                </li>
                                <li @click.stop="notifications_toggle">
                                    <a v-ripple class="flex justify-content-center align-items-center cursor-pointer px-4 py-3 rounded text-surface-700 hover:bg-surface-100 dark:text-surface-0 dark:hover:bg-surface-800 duration-150 transition-colors p-ripple">
                                        <span class="font-medium">Danger</span>
                                        <span class="inline-flex items-center justify-center ml-auto text-primary-contrast rounded-full" style="min-width: 1.5rem; height: 1.5rem">
                                            <Badge :value="notifications_severity_counter[3]" class="p-badge-danger"  />
                                        </span>
                                    </a>
                                </li>
                            </ul>
                            <OverlayPanel ref="notifications_op" class="w-3" style="max-height:60vh;overflow-y: auto;">
                                <h4 class="my-0 mx-2" style="color:#c2c2c2">{{ t('notifications') }}</h4>
                                <Button text :label="t('clear_all')" severity="secondary" @click="clearNotifications()"/>
                                <div class="flex flex-column-reverse">
                                    <NotificationView @closed="notifications.splice(index,1)" :notification="notification" v-for="(notification,index) in notifications" :key="notification.id" />
                                </div>
                            </OverlayPanel>
                        </li>
                        <li @click.stop="chats_toggle">
                            <a v-ripple class="flex items-center align-items-center cursor-pointer px-4 py-3 rounded text-surface-700 hover:bg-surface-100 dark:text-surface-0 dark:hover:bg-surface-800 duration-150 transition-colors p-ripple">
                                <i class="pi pi-comments mr-2"></i>
                                <span class="font-medium">Messages</span>
                                <span class="inline-flex items-center justify-center ml-auto text-primary-contrast rounded-full" style="min-width: 1.5rem; height: 1.5rem">
                                    <Badge class="p-badge-danger" v-if="has_new_message"  />
                                </span>
                            </a>
                            <OverlayPanel ref="chats_op" class="w-5 lg:w-4" style="max-height:90vh;overflow-y: auto;">
                                <Panel :header="t('chats')">

                                    <div style="height:40vh;overflow-y: auto;" ref="chat_container">                            
                                        <div v-for="(chat,index) in chats" :key="index" :class="`flex ${chat.user_sub == user.sub ? 'justify-content-end' : ''}`">
                                            <InlineMessage severity="success" v-if="chat.user_sub != user.sub" class="mt-2">
                                                <template #default>
                                                    <div class="px-3 flex flex-column">
                                                        <strong>{{ chat.sender_name }}</strong>
                                                        <span class="pt-2 px-2">{{ chat.message }}</span>
                                                    </div>
                                                </template>
                                            </InlineMessage>
                
                                            <InlineMessage severity="info" v-if="chat.user_sub == user.sub" class="mt-2">
                                                <template #default>
                                                    <div class="px-3 flex flex-column">
                                                        <strong>{{ chat.sender_name }}</strong>
                                                        <span class="pt-2 px-2">{{ chat.message }}</span>
                                                    </div>
                                                </template>
                                            </InlineMessage>
                                        </div>
                                    </div>

                                </Panel>

                                <InputGroup class="mt-2">
                                    <InputText v-model="chat_text" :placeholder="t('write_message')+'..'" @keyup.enter="SendChatMessage(chat_text)" />
                                    <Button icon="pi pi-send" severity="info" @click="SendChatMessage(chat_text)" />
                                </InputGroup>
                            </OverlayPanel>
                        </li>
                        <li @click.stop="stashed_toggle">
                            <a v-ripple class="flex justify-content-center align-items-center cursor-pointer px-4 py-3 rounded text-surface-700 hover:bg-surface-100 dark:text-surface-0 dark:hover:bg-surface-800 duration-150 transition-colors p-ripple">
                                <i class="pi pi-bookmark mr-2"></i>
                                <span class="font-medium">Drafts</span>
                                <span class="inline-flex items-center justify-center ml-auto text-primary-contrast rounded-full" style="min-width: 1.5rem; height: 1.5rem">
                                    <Badge :value="stashedOrders.length" class="p-badge-success"  />
                                </span>
                            </a>
                            <OverlayPanel ref="stashed_orders_op" class="w-5 lg:w-3" style="max-height:60vh;overflow-y: auto;">
                                <h4 class="m-2" style="color:#c2c2c2">{{ t('stashed_orders') }}</h4>
                                <StashedOrder :order="order" v-for="(order,index) in stashedOrders" :key="index" @back_to_checkout="BackStashedOrderToCheckout(index)" />
                            </OverlayPanel>
                        </li>
                    </ul>
                </div>
            </div>
        </template>
    </Drawer>
    <Dialog v-model:visible="version_dialog_visible" header="Nutrix" :style="{ width: '45rem' }">
        <p class="text-justify">
            Nutrix is an open-source restaurant management system
            designed to make managing your restaurant easy and efficient.
            It's built with modern web technologies and provides a simple
            and intuitive interface to manage your menu, orders, customers,
            and more. Nutrix is completely free and open source under the GPL-2 license, meaning
            you have complete control over the system and can modify it
            to suit your needs. With Nutrix, you can focus on what matters
            most - providing great food and service to your customers.
        </p>
        <p>
            For more support & collaboration visit &nbsp;<a style="font-size:large;" href="https://nutrixpos.com" target="_blank"><i class="pi pi-external-link mr-2"></i>https://nutrixpos.com </a>
        </p>
        <p>
            version / commit hash : {{ app_version }}
        </p>
    </Dialog>
</template>

<script setup lang="ts">
  import InputIcon from 'primevue/inputicon';
  import IconField from 'primevue/iconfield';
  import Toolbar from 'primevue/toolbar';
  import Dialog from 'primevue/dialog';
  import Panel from 'primevue/panel';
  import InputText from 'primevue/inputtext';
  import Chip from 'primevue/chip';
  import DataTable from 'primevue/datatable';
  import Column from 'primevue/column';
  import Button from 'primevue/button';
  import { useToast } from "primevue/usetoast";
  import axios from 'axios'
  import OrderItemView from '@/components/OrderItemView.vue'
  import {OrderItem} from '@/classes/OrderItem'
  import Order from '@/classes/Order'
  import Divider from 'primevue/divider';
  import Slider from 'primevue/slider';
  import Badge from 'primevue/badge'
  import NotificationView from '@/components/NotificationView.vue';
  import OverlayPanel from 'primevue/overlaypanel';
  import { Notification} from '@/classes/Notification';
  import { ref,watch,computed,getCurrentInstance, nextTick, useTemplateRef, version  } from "vue";
  import StashedOrder from '@/components/StashedOrder.vue'
  import InlineMessage from 'primevue/inlinemessage'
  import MainSearchResultView from '@/components/MainSearchResultView.vue';
  import OrderView from '@/components/OrderView.vue';
  import ProgressSpinner from "primevue/progressspinner";
  import StepList from 'primevue/steplist';
  import Step from 'primevue/step';
  import Stepper from 'primevue/stepper';
  import StepPanels from 'primevue/steppanels';
  import StepPanel from 'primevue/steppanel';
  import Card from 'primevue/card';
  import Popover from 'primevue/popover';
  import Textarea from 'primevue/textarea';
  import PickCustomer from '@/components/PickCustomer.vue';
  import AddCustomer from '@/components/AddCustomer.vue';
  import { useI18n } from 'vue-i18n'
  import { ToggleButton,Drawer,Avatar,ButtonGroup, Select } from 'primevue';
  import { globalStore } from '@/stores';


const app_version = ref("")
const version_dialog_visible = ref(false)
const payment_sources = ref<any[]>([])

const zitadel_enabled = ref(true)


const { t } = useI18n() 
const { proxy } = getCurrentInstance();
const store = globalStore()

const payment_source = ref()
const current_order_tip = ref(0)
const is_print_receipt_client = ref(true)
const is_print_receipt_kitchen = ref(true)
const is_auto_start_order = ref(false)
const is_auto_finish_order = ref(true)
const is_serve_inside = ref(true)
const is_delivery = ref(false)
const is_take_away = ref(false)
const is_collecting_money = ref(true)
const is_pay_later = ref(false)
const order_comment = ref("")
const new_order_delivery_customer = ref([])
const pick_customer_dialog = ref(false)
const add_customer_dialog = ref(false)
const new_custom_data_value = ref("")
const new_custom_data_key = ref("")
const custom_data : any = ref([])
const drawer_visible = ref(false)

const toast = useToast();
const itemToEditIndex = ref(0)
const edit_item_dialog = ref(false)
const is_order_valid = ref(true)
const chat_text = ref("")
const chats = ref<any[]>([])
const has_new_message = ref(false)
const chat_container = useTemplateRef("chat_container")
const mainSearchText = ref("")
const order_details_dialog = ref(false)
const order_additional_details_dialog = ref(false)
const order_to_show = ref<Order>()
const isSearchingProduct = ref(false)

import MealCard from '@/components/MealCard.vue';


const comment = ref("")
const subtotal = ref(0)
const discount = ref(0)
const discount_percent = ref(0)
const total = ref(0)
const namewithcomment = ref("")
const idwithcomment = ref("")
const visible = ref(false)
const selectedCategory = ref();

const stashedOrders = ref<Order[]>([])
const inProgressOrders = ref<Order[]>([])

const notifications_op = ref();
const stashed_orders_op = ref();
const chats_op = ref();
const current_orders_op = ref()
const mainsearch_op = ref()

const mainSearchResult = ref<any[]>([])

const discount_op = ref();

const order_details_steps : any = ref([
    {"number": 1, "label": "Main details"},
    {"number": 3, "label": "Confirmation"},
])

const toggleDarkMode = () => {
    store.toggleDarkMode()
    document.documentElement.classList.toggle('my-app-dark');
}

watch(is_delivery, (new_val) => {

    if (new_val){
        is_serve_inside.value = false
        is_take_away.value = false
        order_details_steps.value = [
            {"number": 1, "label": "Main details"},
            {"number": 2, "label": "Delivery"},
            {"number": 3, "label": "Confirmation"},
        ]
    }else {
        order_details_steps.value = [
            {"number": 1, "label": "Main details"},
            {"number": 3, "label": "Confirmation"},
        ]
    }

})

watch(is_serve_inside, (new_val) => {

    if (new_val){
        is_delivery.value = false
        is_take_away.value = false
    }

})

watch(is_take_away, (new_val) => {

    if (new_val){
        is_delivery.value = false
        is_serve_inside.value = false
    }

})


watch(is_collecting_money, (new_val) => {

    if (new_val)
        is_pay_later.value = false

})

watch(is_pay_later, (new_val) => {

if (new_val)
    is_collecting_money.value = false

})

const toggle_discount_popover = (event) => {
    discount_op.value.toggle(event);
}


const user : any = computed(() => {

    return proxy.$zitadel?.oidcAuth.userProfile

})

const finishOrderDisplayed = () => {
    if (order_to_show.value){
        order_to_show.value.state = "finished"
    }
}

const cancelOrderDisplayed = () => {
    if (order_to_show.value){
        order_to_show.value.state = "cancelled"
    }
    
    order_details_dialog.value=false

    getCurrentOrders()
}

const orderToShowAmountCollected = () => {

    if (order_to_show.value){
        order_to_show.value.is_paid = true
    }

    getCurrentOrders()
}

const mainSearchTextChanged = (event:any) => {

    if (mainSearchText.value != ""){
        mainsearch_op.value.show(event)

        axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders?filter[display_id]=${mainSearchText.value}`, {
            headers: {
                Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
            }
        })
        .then(response => {
            mainSearchResult.value = []
            for (var i=0;i<response.data.data.length;i++){
                mainSearchResult.value.push(response.data.data[i])
            }
        })
        .catch(error => {
            console.error(error);
        });

    }else{
        mainSearchResult.value = []
        mainsearch_op.value.hide()
    }


    console.log(mainSearchText.value)
}

const getCurrentOrders = () => {

    inProgressOrders.value = []

    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders?filter[is_paid]=false&filter[state]=!stashed`,{
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        inProgressOrders.value = inProgressOrders.value.concat(response.data.data)
    })


    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders?filter[is_paid]=true&filter[state]=in_progress&filter[state]=pending&filter[state]=!stashed`,{
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        inProgressOrders.value = inProgressOrders.value.concat(response.data.data)
    })

    
}

const claims : any = computed(() => {

    if (user.value) { 
        return Object.keys(user.value).map(key => ({
          key,
          value: user.value[key]
        }))
      }

      return []

})

const roles : any = computed(()=>{
    if (claims.value.length > 0){

        for (var i=0;i<claims.value.length;i++){
            if (claims.value[i].key == "urn:zitadel:iam:org:project:roles"){
                return Object.keys(claims.value[i].value).map(key => {
                    return key
                })
            }
        }
    }

    return []
})


const BackStashedOrderToCheckout = async (stashed_order_index:number) => {

    const order = stashedOrders.value[stashed_order_index]
    let tmp_subtotal = 0


    for (var index=0;index<order.items.length;index++){
        const tmp_order_item = new OrderItem()
        await tmp_order_item.FromItemData(order.items[index])
        await tmp_order_item.RefreshProductData()
        tmp_order_item.price = order.items[index].product.price
        tmp_subtotal+=order.items[index].price
        order.items[index] = tmp_order_item
    }


    subtotal.value = tmp_subtotal
    discount_percent.value =  isNaN(order.discount * 100 / tmp_subtotal) ? 0 : order.discount * 100 / tmp_subtotal
    orderItems.value = order.items


    // discount.value = order.discount == null || order.discount == undefined ? 0 : order.discount


    axios.delete(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders/${order.id}`,
    {
        headers:{
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        },
    },)
    .then(()=>{
        stashedOrders.value.splice(stashed_order_index,1)

        stashed_orders_op.value.toggle()
    })
    .catch(() => {
        toast.add({severity:'error', summary: 'Failed to remove stashed order from db', detail: "", life: 3000,group:'br'});
        stashed_orders_op.value.toggle()
    })


}

const notifications_toggle = (event: any) => {
    notifications_op.value.toggle(event);
}

const stashed_toggle = (event: any) => {
    stashed_orders_op.value.toggle(event);
}

const paylater_toggle = (event: any) => {
    current_orders_op.value.toggle(event);
}


const chats_toggle = async (event: any) => {
    has_new_message.value = false
    chats_op.value.toggle(event);
    await nextTick()

    if (chat_container.value != null)
        chat_container.value.scrollTop = chat_container.value?.scrollHeight

}

const notifications = ref<Notification[]>([])


const getStashedOrders = () => {
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders?filter[state]=stashed`,{
        headers:{
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    }).then(async (response) => {

        const dataCopy = JSON.parse(JSON.stringify(response.data.data))

        for (var i=0;i<dataCopy.length;i++){

            let order = new Order()
            order = JSON.parse(JSON.stringify(dataCopy[i]))
            order.items = []


            for (var j=0;j<dataCopy[i].items.length;j++){

                const item = new OrderItem()
                await item.FromItemData(dataCopy[i].items[j])

                order.items.push(item)
            }

            stashedOrders.value.push(order)
        }
    }).catch(() => {
        toast.add({severity:'error', summary: 'Failed to get stashed orders', detail: "", life: 3000,group:'br'});
    })
}


const stashOrder = () => {

    const order = new Order()
    order.items = orderItems.value
    order.discount = discount.value == null ? 0 : discount.value
    order.state = "stashed"


    axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders`,
        {
            data: order
        },
    {
        headers:{
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`,
            "Accept-Language": proxy.$i18n.locale
        },
    }).then(async (response) => {
        orderItems.value=[]
        discount.value = 0
        discount_percent.value = 0
        total.value =0
        subtotal.value = 0


        console.log(response)


        for (var index=0;index<response.data.data.items.length;index++){
            const temp_order_item = new OrderItem()
            await temp_order_item.FromItemData(response.data.data.items[index])
            await temp_order_item.RefreshProductData()
            temp_order_item.ValidateItem()
            response.data.data.items[index] = temp_order_item
        }

        
        stashedOrders.value.push(response.data.data)
        toast.add({severity:'success', summary: `Order ${order.display_id} stashed successfully !`, detail: "successfully stashed order !", life: 3000,group:'br'});
    }).catch( () => {
        toast.add({severity:'error', summary: 'Error Stashing Item', detail: "", life: 3000,group:'br'});
    })
}


const clearNotifications = () => {
    notifications.value = []
}

let socket : WebSocket


const SendChatMessage = (msg: string) => {
    if (socket.readyState === WebSocket.OPEN) {
        socket.send(`{"type":"chat_message","message":"${msg}","sender_name":"${user.value.name}","user_sub":"${user.value.sub}","to":"*","date": "${new Date().toLocaleString()}"}`)
    }else {
        console.log("WS closed")
    }
    chat_text.value = ""
}


const startWebsocket = () => {
    socket = new WebSocket(`ws://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/ws`);
    socket.onopen = () => {
        console.log("Opened ws connection");
        socket.send(`{"type":"subscribe","topic_name":"all"}`);
    }

    socket.onmessage = async (event) => {
        console.log("Received message: " + event.data);

        const data = JSON.parse(event.data);

        if (data.type == "topic_message") {
            if (data.topic_name == "order_finished"){

                toast.removeGroup('br')
                toast.add({ severity: 'success', summary: 'Order Finished', detail: `order ( ${data.order_id} ) finished and is ready to be served !`, life: 3000,group:'br' });

                const notification = new Notification();
                notification.description = `order (${data.order_id}) finished and is ready to be served !`
                notification.severity = "success"
                notification.topic_name = "Order Finished"
                notification.type = "topic_message"
                notification.date = data.date
                notifications.value.push(notification);


                getCurrentOrders()
            }else {
                const notification = new Notification();
                notification.description = data.message
                notification.severity = data.severity
                notification.topic_name = data.topic_name
                notification.type = data.type
                notification.date = data.date
                notifications.value.push(notification);

                toast.removeGroup('br')
                toast.add({ severity: data.severity, summary: data.topic_name, detail: data.message, life: 30000,group:'br' });
            }
        }

        if (data.type == "chat_message") {

            if (!chats_op.value?.visible) {
                has_new_message.value = true
            }

            chats.value.push({
                message:data.message,
                sender_name: data.sender_name,
                user_sub: data.user_sub,
                date: data.date,
            })


            if (chat_container.value != null){
                await nextTick()
                chat_container.value.scrollTop = chat_container.value?.scrollHeight + 100
            }
            
        }

    }
    socket.onerror = (event) => {
        console.log("Error occurred");
        console.log(event);
    }
    socket.onclose = () => {
        console.log("Connection closed");
        const retryConnection = async () => {
            if (socket.readyState !== WebSocket.OPEN) {
                await new Promise(r => setTimeout(r, 5000));
                console.log("Reconnecting to WebSocket...");
                startWebsocket()
            }
        }
        retryConnection();
    }
}

const loading = ref(true)
const { locale,setLocaleMessage } = useI18n({ useScope: 'global' })

const loadSettings = async () => {

    await axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/settings`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        },
    })
    .then(async (response)=>{

        payment_sources.value = response.data.data.payment_sources == null ? [] : response.data.data.payment_sources
        if (payment_sources.value.length > 0){
            payment_source.value = {"name": payment_sources.value[0].name}
        }

        await axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/languages/${response.data.data.language.code}`, {
            headers: {
                Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
            }
        })
        .then(response2 => {

            setLocaleMessage(response2.data.data.code, response2.data.data.pack);
            locale.value = response2.data.data.code
            store.setOrientation(response2.data.data.orientation)
            loading.value = false
        })
        .catch((err) => {
            console.log(err)
        });
        loading.value = false
    })
    .catch((err) => {
        console.log(err)
    });

}

const init = async () => {

    console.log(import.meta.env)

    app_version.value = import.meta.env.VITE_APP_APP_VERSION || ""

    if (import.meta.env.VITE_APP_ZITADEL_ENABLED === 'false'){
        zitadel_enabled.value = false
    }

    if (!zitadel_enabled.value || proxy.$zitadel?.oidcAuth.isAuthenticated){
        await loadSettings()
    }else {
        window.location.href = import.meta.env.VITE_APP_ZITADEL_ISSUER
        loading.value = false
    }


    startWebsocket()
    getStashedOrders()
    getCurrentOrders()

    console.log("user:")
    console.log(user.value)

    console.log("claims:")
    console.log(claims.value)
}

init()

// const notifications_severity_counter = ref<number[]>([])

const notifications_severity_counter = computed(() => {
    const counter = [0,0,0,0]
    notifications.value.forEach(notification => {
        switch (notification.severity) {
            case "success":
                counter[0]++;
                break;
            case "info":
                counter[1]++;
                break;
            case "warn":
                counter[2]++;
                break;
            case "error":
                counter[3]++;
                break;
        }
    })
    return counter
})

  
  
const searchtext = ref("")
const categories : any = ref([])

const orderItems = ref<OrderItem[]>([])


const addItem = async (item) => {

    const new_item = new OrderItem()
    new_item.product.name = item.name
    new_item.product.id = item.id
    await new_item.ReloadDefaults()
    new_item.ValidateAllMaterials()
    new_item.ValidateItem()


    orderItems.value.push(new_item)
}

const addWithComment = async () => {

    const new_item = new OrderItem()
    new_item.product.name = namewithcomment.value
    new_item.comment = comment.value
    new_item.product.id = idwithcomment.value
    await new_item.ReloadDefaults()
    new_item.ValidateItem()

    orderItems.value.push(new_item)
    visible.value=false
    comment.value = ""
    idwithcomment.value = ""
}


const getCategories = async () => {
    const response = await axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/categories`,{
        headers:{
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    categories.value = categories.value.concat(response.data.data)
    selectedCategory.value = categories.value[0]
}

getCategories();


const submitOrder = () => {

    let custom_data_map : any = {}
    custom_data.value.forEach((item : any) => {
        custom_data_map[item.key] = item.value
    })

    let order : any =  {
        items:orderItems.value,
        discount:discount.value,
        is_auto_start: is_auto_start_order.value,
        is_auto_finish: is_auto_finish_order.value,
        is_dine_in: is_serve_inside.value,
        is_take_away: is_take_away.value,
        is_delivery: is_delivery.value,
        is_paid: is_collecting_money.value,
        is_pay_later: is_pay_later.value,
        tips: current_order_tip.value,
        payment_source: payment_source.value.name,
        custom_data: custom_data_map,
        comment: order_comment.value,
        customer: new_order_delivery_customer.value.length > 0 ? new_order_delivery_customer.value[0] : null
    }

    if (orderItems.value.length > 0){

        axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders`,
            {
                meta: {
                    is_print_client_receipt: is_print_receipt_client.value,
                    is_print_kitchen_receipt: is_print_receipt_kitchen.value
                },
                data:order
            },
            {
                headers:{
                    Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`,
                    'Accept-Language': proxy.$i18n.locale,
                },
            },
        ).then(() => {
            toast.add({ severity: 'success', summary: 'Success', detail: 'Order in progress !', life: 3000,group:'br' });
            refreshAvailabilities()
            if (!is_collecting_money.value || is_auto_start_order.value){
                getCurrentOrders()
            }

            is_delivery.value = false
            is_take_away.value = false
            order_additional_details_dialog.value = false
        })
    
        orderItems.value = []
        
    }
};


const isUpdatingDiscount = ref(false)
const isUpdatingDiscountPercent = ref(false)


watch(searchtext, (newSearchText) => {

    isSearchingProduct.value = true
    
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/products?filter[search]=${newSearchText}`,
    {
        headers:{
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    }
    ).then((response) => {
        products.value = []
        products.value = response.data.data
        refreshAvailabilities();
        isSearchingProduct.value=false
    })

})

watch(subtotal, (new_subtotal) => {
  total.value = new_subtotal - discount.value
  if (total.value < 0)
    total.value = 0
})

watch(discount, (new_discount) => {
  if (!isUpdatingDiscountPercent.value){
    isUpdatingDiscount.value = true
    total.value = subtotal.value - new_discount

    if (new_discount != 0)
        discount_percent.value = new_discount*100 / subtotal.value
    if (total.value < 0)
    total.value = 0
    }else{
      isUpdatingDiscountPercent.value = false
  }
})

watch(discount_percent, (new_discount_percent) => {
 if (!isUpdatingDiscount.value){
    isUpdatingDiscountPercent.value = true
    discount.value = new_discount_percent * subtotal.value / 100
    total.value = subtotal.value - discount.value
    isUpdatingDiscount.value = false
  }else {
    isUpdatingDiscount.value = false
  }
})


watch(() => orderItems.value, 
  (currentValue) => {
    let x = 0
    let valid = true
    currentValue.forEach((item) => {

        x += item.price
        if (!item.isValid)
            valid=false

    })

    is_order_valid.value = valid
    subtotal.value = x
    discount.value = discount_percent.value * subtotal.value / 100
  },
  {deep: true}
);

  
const products = ref([
])



const refreshAvailabilities = () => {
    var product_ids = ""

    if (products.value != null){
        
        products.value.forEach((product,index) => {
            product_ids += index > 0 ? "," +product.id : product.id
        })

        axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/products/availability?ids=`+product_ids,{
            headers:{
                Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
            }
        })
        .then((response) => {
            products.value.forEach((product,index) => {
                products.value[index].availability = Math.round(response.data.data.filter((x) => x.recipe_id == product.id)[0].available * 100) / 100
            })
        })
    }
}


watch(selectedCategory, (category) => {
    if (category != null){
        products.value = []
        category.products.forEach((recipe) => {
            products.value.push({
                id: recipe.id,
                name:recipe.name,
                price:recipe.price,
                image_url: recipe.image_url,
                enable_inventory_consumption: recipe.enable_inventory_consumption
            })
        })
        refreshAvailabilities();
    }
})


  const navbar_links = ref([
    {
          label: {
              title:'cashier',
              plural:false
          },
          icon: 'pi pi-desktop',
          link: '/home',
          focused: true,
          authority: ['cashier','admin']
      },
      {
          label: {
              title:'kitchen',
              plural:false
          },
          focused: false,
          icon: 'fa fa-kitchen-set',
          link:'/kitchen',
          authority: ['chef','admin']
      },
      {
          label: {
              title:'admin',
              plural:false
          },
          focused: false,
          icon: 'pi pi-cog',
          link: '/admin',
          authority: ['admin']
      }
  ]);
  
  
  </script>
  
  <style>
  html,
  body {
    height: 100%;
    margin:0;
  }
  </style>
  