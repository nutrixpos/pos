<template>
    <div class="w-full">
        <div class="grid mx-2">
            <div class="col-12 flex">
                <div class="gird w-full">
                    <div class="col-12">
                        <h3>{{$t('sales')}}</h3>
                    </div>
                    <div class="col-12 flex justify-content-center align-items-center w-full">
                        <div class="flex flex-column w-full">
                            <div class="grid">
                                <div class="col-8">
                                    <Card style="max-width:100%;">
                                        <template #content>
                                            <Line v-if="chartData != undefined" style="min-height: 20rem;max-width:100%;" :data="chartData" :chartOptions="chartOptions" />
                                        </template>
                                    </Card>
                                </div>
                                <div class="col-4 flex justify-content-start align-items-start">
                                    <Card style="max-width:100%;">
                                        <template #content>
                                            <Pie style="max-width: 100%;" v-if="productPiechartData != undefined" :data="productPiechartData" :chartOptions="productPiechartOptions" />
                                        </template>
                                    </Card>
                                </div>
                            </div>
                            <DataTable @page="updatSalesTableRowsPerPage" :lazy="true" :totalRecords="salesTableTotalRecords" :loading="isSalesTableLoading" v-model:expandedRows="expandedSalesLogRows" paginatorPosition="both"  paginator :rows="salesTableRowsPerPage" :rowsPerPageOptions="[7, 14, 30, 90]" :value="sales_log" stripedRows tableStyle="min-width: 50rem;max-height:50vh;" class="w-full pr-2">
                                    <Column expander style="width: 5rem" />
                                    <Column sortable field="date" :header="$t('date')">
                                        <template #body="slotProps">
                                            <i class="pi pi-exclamation-circle" v-tooltip.top="'has refunds'" v-if="slotProps.data.refunds?.length > 0" style="margin-right:0.5rem;color:red"></i>                                            
                                            {{ slotProps.data.date }}
                                        </template>
                                    </Column>
                                    <Column sortable field="costs" :header="$t('cost')">
                                        <template #body="slotProps">
                                            <div class="flex gap-2 align-items-center">
                                                <div>{{ slotProps.data.costs }} </div>
                                                <Badge v-if="slotProps.data.inventory_refunds > 0 " :value="`-${slotProps.data.inventory_refunds}`" severity="success" class="mr-2" />
                                            </div>
                                        </template>
                                    </Column>
                                    <Column sortable field="total_sales" :header="$t('sales')"></Column>
                                    <Column sortable field="refunds_value" :header="$t('refunds')">
                                        <template #body="slotProps">
                                            <Badge :value="`${slotProps.data.inventory_refunds > 0 ? '-' : ''}${slotProps.data.refunds_value}`" severity="secondary" style="margin-right:0.5rem" />
                                        </template>
                                    </Column>
                                    <Column sortable field="profit" :header="$t('profit')">
                                        <template #body="slotProps">
                                            <div :style="`${ ( slotProps.data.total_sales - slotProps.data.costs - slotProps.data.refunds_value + ( slotProps.data.inventory_refunds || 0 )) > 0 ? 'color:green' : 'color:red' }`">
                                                {{ slotProps.data.total_sales - slotProps.data.costs  - slotProps.data.refunds_value + ( slotProps.data.inventory_refunds || 0 ) }}
                                                <Badge :value="`+${slotProps.data.tips} ${$t('tips')}`" severity="secondary" style="margin-right:0.5rem" />
                                            </div>
                                        </template>
                                    </Column>
                                    <template #expansion="slotProps">
                                        <DataTable v-model:expandedRows="expandedSalesLogOrderItems" :value="slotProps.data.orders">
                                            <Column expander style="width: 5rem" />
                                            <Column sortable field="order.display_id" header="Id">
                                                <template #body="slotProps">
                                                    <i class="pi pi-exclamation-circle" v-tooltip.top="'has refunds'" v-if="orders_refunds[slotProps.data.id]?.refunds.length > 0" style="margin-right:0.5rem;color:red"></i>                                            
                                                    {{ slotProps.data.order.display_id }}
                                                </template>
                                            </Column>
                                            <Column sortable field="order.submitted_at" header="Submitted At"></Column>
                                            <Column sortable field="order.cost" header="Cost">
                                                <template #body="slotProps">
                                                    <div class="flex gap-2 align-items-center">
                                                        <div>{{ slotProps.data.order.cost }} </div>
                                                        <Badge v-if="orders_refunds[slotProps.data.order.id]?.inventory_refunds > 0 " :value="`-${orders_refunds[slotProps.data.order.id]?.inventory_refunds}`" severity="success" class="mr-2" />
                                                    </div>
                                                </template>
                                            </Column>
                                            <Column sortable field="order.sale_price" header="Sales"></Column>
                                            <Column sortable field="order.refunds" header="Refunds">
                                                <template #body="slotProps">
                                                    {{ orders_refunds[slotProps.data.id]?.total_refunds }}
                                                </template>
                                            </Column>
                                            <Column sortable field="profit" header="Profit">
                                                <template #body="slotProps">
                                                    <div :style="`${ (slotProps.data.order.sale_price - slotProps.data.order.cost - ( orders_refunds[slotProps.data.order.id]?.total_refunds || 0 ) + (orders_refunds[slotProps.data.order.id]?.inventory_refunds || 0 )) > 0 ? 'color:green' : 'color:red' }`">
                                                        {{ slotProps.data.order.sale_price - slotProps.data.order.cost - (orders_refunds[slotProps.data.order.id]?.total_refunds || 0) + (orders_refunds[slotProps.data.order.id]?.inventory_refunds || 0 ) }} 
                                                        <Badge :value="`+${slotProps.data.order.tips} ${$t('tips')}`" severity="secondary" style="margin-right:0.5rem" />
                                                    </div>
                                                </template>
                                            </Column>
                                            <template #expansion="slotProps">
                                                <SalesLogTableItems :items="slotProps.data.costs" :order_refunds="orders_refunds[slotProps.data.id] || []" />
                                            </template>
                                        </DataTable>
                                    </template>
                            </DataTable>
                        </div>                        
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import DataTable from "primevue/datatable";
import Column from 'primevue/column'
import {Line,Pie} from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale,PointElement,LineElement,ArcElement,Filler} from 'chart.js'
import Card from "primevue/card";
import {getCurrentInstance, ref} from 'vue'
import axios from 'axios'
import SalesLogTableItems from '@/components/SalesLogTableItems.vue'
import { $dt } from '@primevue/themes';
import {Badge} from 'primevue';


ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale,PointElement,LineElement,ArcElement, Filler)


const {proxy} = getCurrentInstance()

const sales_log = ref([])
const orders_refunds = ref({})
const expandedSalesLogRows = ref([])
const expandedSalesLogOrderItems = ref([])
const salesTableRowsPerPage = ref(7)
const salesTableFirstIndex = ref(0)
const salesTableTotalRecords = ref(0)


const chartData = ref();
const chartOptions = ref();


const chartLabels = ref([])
const chartSales = ref([])
const chartCost = ref([])
const chartRefunds= ref([])
const chartInventoryReturns= ref([])
const chartProfit= ref([])

const isSalesTableLoading = ref(true)




const productPieChartColors = ref([])
const productPiechartData = ref();
const productPiechartOptions = ref();
const productPieChartLabels = ref([])
const productPieChartSales = ref([])


const updatSalesTableRowsPerPage = (event) => {

    expandedSalesLogRows.value = []
    expandedSalesLogOrderItems.value = []

    const { first, rows } = event;
    loadSales(first,rows)
}


const setProductPieChartData = () => {
    return {
        labels: productPieChartLabels.value,
        datasets: [
            {
                data: productPieChartSales.value,
                backgroundColor: productPieChartColors.value,
                hoverBackgroundColor: productPieChartColors.value
            }
        ]
    };
};

const setProductPieChartOptions = () => {
    const documentStyle = getComputedStyle(document.documentElement);
    const textColor = documentStyle.getPropertyValue('--p-blue-500');

    return {
        plugins: {
            legend: {
                labels: {
                    usePointStyle: true,
                    color: textColor
                }
            }
        }
    };
};



        
const setChartData = () => {
    return {
        labels: chartLabels.value,
        datasets: [
            {
                label: 'Sales',
                data: chartSales.value,
                fill: false,
                tension: 0.4,
                borderWidth:2,
                borderColor: '#4069ce',
            },
            {
                label: 'Cost',
                data: chartCost.value,
                fill: false,
                tension: 0.4,
                borderWidth:2,
                borderColor: '#8E44AD',
            },
            {
                label: 'Refunds',
                data: chartRefunds.value,
                fill: false,
                tension: 0.4,
                borderWidth:2,
                borderColor: '#C70039',
            },
            {
                label: 'Returns',
                borderWidth:2,
                data: chartInventoryReturns.value,
                fill: true,
                tension: 0.5,
                borderColor: '#FFBF00',
            },
            {
                label: 'Profit',
                data: chartProfit.value,
                fill: true,
                tension: 0.5,
                borderWidth: 4,
                borderColor: '#40ce52',
            },
        ]
    };
};
const setChartOptions = () => {
    const documentStyle = getComputedStyle(document.documentElement);
    const textColorSecondary = documentStyle.getPropertyValue('--text-color-secondary');
    const textColor = documentStyle.getPropertyValue('--text-color');
    const surfaceBorder = documentStyle.getPropertyValue('--surface-border');

    return {
        maintainAspectRatio: false,
        aspectRatio: 0.6,
        plugins: {
            legend: {
                labels: {
                    color: textColor
                }
            }
        },
        scales: {
            x: {
                ticks: {
                    color: textColorSecondary
                },
                grid: {
                    color: surfaceBorder
                }
            },
            y: {
                ticks: {
                    color: textColorSecondary
                },
                grid: {
                    color: surfaceBorder
                }
            }
        }
    };
}


const loadSales = (first=salesTableFirstIndex.value,rows=salesTableRowsPerPage.value) => {


    let page_number = Math.floor(first/rows) + 1


    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/logs/salesperday?page[number]=${page_number}&page[size]=${rows}`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {

        let temp_sales_log = []
        chartLabels.value = []
        chartCost.value = []
        chartSales.value = []
        chartRefunds.value = []
        chartInventoryReturns.value = []
        chartProfit.value = []
        productPieChartLabels.value = []
        productPieChartSales.value = []
        isSalesTableLoading.value = true
        orders_refunds.value = {}

        salesTableTotalRecords.value = response.data.meta.total_records


        for (let i=0;i<response.data.data.length;i++){
            chartLabels.value.push(response.data.data[i].date)
            temp_sales_log.push(response.data.data[i])
            chartSales.value.push(response.data.data[i].total_sales)
            chartCost.value.push(response.data.data[i].costs)
            chartRefunds.value.push(response.data.data[i].refunds_value)

            let day_inventory_refunds = 0
            let tips = 0
            
            for (let j=0;j<response.data.data[i].orders.length;j++){

                tips += response.data.data[i].orders[j].order.tips

                for (let k=0;k<response.data.data[i].orders[j].order.items.length;k++){

                    if (!productPieChartLabels.value.includes(response.data.data[i].orders[j].order.items[k].product.name)){
                        productPieChartLabels.value.push(response.data.data[i].orders[j].order.items[k].product.name)
                        productPieChartSales.value.push(response.data.data[i].orders[j].order.items[k].quantity)
                        productPieChartColors.value.push('#'+(Math.floor(Math.random()*16777215).toString(16)).padStart(6, '0'));
                    }else {

                        let index = productPieChartLabels.value.indexOf(response.data.data[i].orders[j].order.items[k].product.name)
                        productPieChartSales.value[index] += response.data.data[i].orders[j].order.items[k].quantity
                    }   
                }

                
            }

            for (let j=0;j<response.data.data[i].refunds.length;j++){

                let refund = response.data.data[i].refunds[j]

                if (orders_refunds.value[`${response.data.data[i].refunds[j].order_id}`]){
                    orders_refunds.value[`${response.data.data[i].refunds[j].order_id}`].refunds.push(refund)
                    orders_refunds.value[`${response.data.data[i].refunds[j].order_id}`].total_refunds += refund.amount
                    orders_refunds.value[`${response.data.data[i].refunds[j].order_id}`].inventory_refunds = 0

                }
                else{
                    orders_refunds.value[`${response.data.data[i].refunds[j].order_id}`] = {
                        refunds : [refund],
                        inventory_refunds: 0,
                        total_refunds : refund.amount
                    }
                }

                if (refund.destination == "inventory"){
                    // inventory_refund = the cost of material consumed not the customer sale price
                    day_inventory_refunds += refund.item_cost
                    orders_refunds.value[`${response.data.data[i].refunds[j].order_id}`].inventory_refunds += refund.item_cost
                }

                if (refund.destination == "custom"){
                    refund.material_refunds.forEach((material_refund) => {
                        orders_refunds.value[`${response.data.data[i].refunds[j].order_id}`].inventory_refunds += material_refund.inventory_return_qty * material_refund.cost_per_unit
                        day_inventory_refunds += material_refund.inventory_return_qty * material_refund.cost_per_unit
                    })
                }
                

            }


            temp_sales_log[i].inventory_refunds = day_inventory_refunds
            temp_sales_log[i].tips = tips
            chartInventoryReturns.value.push(temp_sales_log[i].inventory_refunds)
            chartProfit.value.push((temp_sales_log[i].total_sales - temp_sales_log[i].costs - temp_sales_log[i].refunds_value + ( temp_sales_log[i].inventory_refunds || 0 )))
        }

        sales_log.value = temp_sales_log


        chartData.value = setChartData();
        chartOptions.value = setChartOptions();


        productPiechartData.value = setProductPieChartData();
        productPiechartOptions.value = setProductPieChartOptions();



        isSalesTableLoading.value = false

    })
    .catch(error => {
        console.log(error)
    })

}

loadSales()
</script>