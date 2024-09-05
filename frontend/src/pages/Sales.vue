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
                            <DataTable v-model:expandedRows="expandedSalesLogRows" :value="sales_log" stripedRows tableStyle="min-width: 50rem" class="w-full pr-5">
                                    <Column expander style="width: 5rem" />
                                    <Column sortable field="date" header="Date"></Column>
                                    <Column sortable field="cost" header="Cost"></Column>
                                    <Column sortable field="sales" header="Sales"></Column>
                                    <Column sortable field="profit" header="Profit">
                                        <template #body="slotProps">
                                            <div :style="`${ (slotProps.data.sales - slotProps.data.cost) > 0 ? 'color:green' : 'color:red' }`">{{ slotProps.data.sales - slotProps.data.cost }}</div>
                                        </template>
                                    </Column>
                                    <template #expansion="slotProps">
                                        <DataTable v-model:expandedRows="expandedSalesLogOrderItems" :value="slotProps.data.orders">
                                            <Column expander style="width: 5rem" />
                                            <Column sortable field="date" header="Date"></Column>
                                            <Column sortable field="cost" header="Cost"></Column>
                                            <Column sortable field="sale_price" header="Sales"></Column>
                                            <Column sortable field="profit" header="Profit">
                                                <template #body="slotProps">
                                                    <div :style="`${ (slotProps.data.sale_price - slotProps.data.cost) > 0 ? 'color:green' : 'color:red' }`">{{ slotProps.data.sale_price - slotProps.data.cost }}</div>
                                                </template>
                                            </Column>
                                            <template #expansion="slotProps">
                                                <SalesLogTableItems :items="slotProps.data.Items" />
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
import {ref} from 'vue'
import axios from 'axios'
import SalesLogTableItems from '@/components/SalesLogTableItems.vue'

const sales_log = ref([])
const expandedSalesLogRows = ref([])
const expandedSalesLogOrderItems = ref([])


const chartData = ref();
const chartOptions = ref();


const chartLabels = ref([])
const chartSales = ref([])
const chartCost = ref([])




const productPiechartData = ref();
const productPiechartOptions = ref();
const productPieChartLabels = ref([])
const productPieChartSales = ref([])

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


const loadSales = () => {
    axios.get('http://localhost:8000/api/sales_logs')
    .then(response => {
        var per_day_log = {}
        var product_sale_count = {}


        response.data.forEach((log) => {

            const date = new Date(log.date);

            // Get the date
            const year = date.getFullYear();
            const month = String(date.getMonth() + 1).padStart(2, '0');
            const day = String(date.getDate()).padStart(2, '0');


            if (!per_day_log[`${year}-${month}-${day}`]){
                per_day_log[`${year}-${month}-${day}`] = {
                    cost: 0.0,
                    sales: 0.0,
                    orders: [],
                    date: `${year}-${month}-${day}` 
                }

                chartLabels.value.push(`${year}-${month}-${day}`)
            }


            log.Items.forEach((component) => {
                if (!product_sale_count[component.ItemName]){
                    product_sale_count[component.ItemName] = 0
                }
                
                product_sale_count[component.ItemName]++;
            })


            per_day_log[`${year}-${month}-${day}`].orders.push(log)
            per_day_log[`${year}-${month}-${day}`].cost += log.cost
            per_day_log[`${year}-${month}-${day}`].sales += log.sale_price
        })


        for (var day in per_day_log){
            sales_log.value.push(per_day_log[day])
            chartSales.value.push(per_day_log[day].sales)
            chartCost.value.push(per_day_log[day].cost)
        }

        for (var product in product_sale_count){
            productPieChartLabels.value.push(product)
            productPieChartSales.value.push(product_sale_count[product])
        }


        chartData.value = setChartData();
        chartOptions.value = setChartOptions();

        setTimeout(() => {            
            productPiechartData.value = setProductPieChartData();
            productPiechartOptions.value = setProductPieChartOptions();
        }, 2000);
    })
    .catch(error => {
        console.log(error)
    })
}

loadSales()
</script>