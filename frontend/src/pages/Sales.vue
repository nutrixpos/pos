<template>
    <div class="w-full">
        <div class="grid mx-2">
            <div class="col-12 flex">
                <div class="gird w-full">
                    <div class="col-12">
                        <h3>Sales</h3>
                    </div>
                    <div class="col-12 flex justify-content-center align-items-center w-full">
                        <div class="flex flex-column w-full">
                            <div class="grid">
                                <div class="col-8">
                                    <div class="card">
                                        <Chart style="min-height: 20rem;" type="line" :data="chartData" :options="chartOptions" />
                                    </div>
                                </div>
                                <div class="col-4 flex justify-content-center align-items-center">
                                    <div class="card">
                                        <Chart type="pie" class="w-20rem" :data="productPiechartData" :options="productPiechartOptions" />
                                    </div>
                                </div>
                            </div>
                            <DataTable @page="updatSalesTableRowsPerPage" :lazy="true" :totalRecords="salesTableTotalRecords" :loading="isSalesTableLoading" v-model:expandedRows="expandedSalesLogRows" paginatorPosition="both"  paginator :rows="salesTableRowsPerPage" :rowsPerPageOptions="[7, 14, 30, 90]" :value="sales_log" stripedRows tableStyle="min-width: 50rem;max-height:50vh;" class="w-full pr-2">
                                    <Column expander style="width: 5rem" />
                                    <Column sortable field="date" header="Date"></Column>
                                    <Column sortable field="costs" header="Costs"></Column>
                                    <Column sortable field="total_sales" header="Sales"></Column>
                                    <Column sortable field="profit" header="Profit">
                                        <template #body="slotProps">
                                            <div :style="`${ (slotProps.data.total_sales - slotProps.data.costs) > 0 ? 'color:green' : 'color:red' }`">{{ slotProps.data.total_sales - slotProps.data.costs }}</div>
                                        </template>
                                    </Column>
                                    <template #expansion="slotProps">
                                        <DataTable v-model:expandedRows="expandedSalesLogOrderItems" :value="slotProps.data.orders">
                                            <Column expander style="width: 5rem" />
                                            <Column sortable field="order.display_id" header="Id"></Column>
                                            <Column sortable field="order.submitted_at" header="Submitted At"></Column>
                                            <Column sortable field="order.cost" header="Cost"></Column>
                                            <Column sortable field="order.sale_price" header="Sales"></Column>
                                            <Column sortable field="profit" header="Profit">
                                                <template #body="slotProps">
                                                    <div :style="`${ (slotProps.data.order.sale_price - slotProps.data.order.cost) > 0 ? 'color:green' : 'color:red' }`">{{ slotProps.data.order.sale_price - slotProps.data.order.cost }}</div>
                                                </template>
                                            </Column>
                                            <template #expansion="slotProps">
                                                <SalesLogTableItems :items="slotProps.data.costs" />
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

<script setup>
import DataTable from "primevue/datatable";
import Column from 'primevue/column'
import Chart from 'primevue/chart';
import {getCurrentInstance, ref} from 'vue'
import axios from 'axios'
import SalesLogTableItems from '@/components/SalesLogTableItems.vue'

const {proxy} = getCurrentInstance()

const sales_log = ref([])
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

const isSalesTableLoading = ref(true)




const productPiechartData = ref();
const productPiechartOptions = ref();
const productPieChartLabels = ref([])
const productPieChartSales = ref([])


const updatSalesTableRowsPerPage = (event) => {
    const { first, rows } = event;
    loadSales(first,rows)
}


const setProductPieChartData = () => {
    const documentStyle = getComputedStyle(document.body);

    return {
        labels: productPieChartLabels,
        datasets: [
            {
                data: productPieChartSales,
                backgroundColor: [documentStyle.getPropertyValue('--cyan-500'), documentStyle.getPropertyValue('--orange-500'), documentStyle.getPropertyValue('--gray-500')],
                hoverBackgroundColor: [documentStyle.getPropertyValue('--cyan-400'), documentStyle.getPropertyValue('--orange-400'), documentStyle.getPropertyValue('--gray-400')]
            }
        ]
    };
};

const setProductPieChartOptions = () => {
    const documentStyle = getComputedStyle(document.documentElement);
    const textColor = documentStyle.getPropertyValue('--text-color');

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
    const documentStyle = getComputedStyle(document.documentElement);

    return {
        labels: chartLabels.value,
        datasets: [
            {
                label: 'Sales',
                data: chartSales.value,
                fill: false,
                tension: 0.4,
                borderColor: documentStyle.getPropertyValue('--cyan-500'),
            },
            {
                label: 'Cost',
                data: chartCost.value,
                fill: false,
                tension: 0.4,
                borderColor: documentStyle.getPropertyValue('--orange-300'),
            },
        ]
    };
};
const setChartOptions = () => {
    const documentStyle = getComputedStyle(document.documentElement);
    const textColor = documentStyle.getPropertyValue('--text-color');
    const textColorSecondary = documentStyle.getPropertyValue('--text-color-secondary');
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

    axios.get(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/salesperday?first_index=${first}&rows=${rows}`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel.oidcAuth.accessToken}`
        }
    })
    .then(response => {

        let temp_sales_log = []
        chartLabels.value = []
        chartCost.value = []
        chartSales.value = []
        productPieChartLabels.value = []
        productPieChartSales.value = []
        isSalesTableLoading.value = true

        salesTableTotalRecords.value = response.data.total_records


        for (let i=0;i<response.data.sales.length;i++){
            chartLabels.value.push(response.data.sales[i].date)
            temp_sales_log.push(response.data.sales[i])
            chartSales.value.push(response.data.sales[i].total_sales)
            chartCost.value.push(response.data.sales[i].costs)

            for (let j=0;j<response.data.sales[i].orders.length;j++){



                for (let k=0;k<response.data.sales[i].orders[j].order.items.length;k++){

                    if (!productPieChartLabels.value.includes(response.data.sales[i].orders[j].order.items[k].product.name)){
                        productPieChartLabels.value.push(response.data.sales[i].orders[j].order.items[k].product.name)
                        productPieChartSales.value.push(response.data.sales[i].orders[j].order.items[k].quantity)
                    }else {

                        let index = productPieChartLabels.value.indexOf(response.data.sales[i].orders[j].order.items[k].product.name)
                        productPieChartSales.value[index] += response.data.sales[i].orders[j].order.items[k].quantity
                    }   
                }
                
            }
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