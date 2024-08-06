<template>
    <div class="w-full">
        <div class="grid m-2">
            <div class="col-12 flex pt-3">
                <div class="gird w-full">
                    <div class="col-12">
                        <h3>Sales</h3>
                    </div>
                    <div class="col-12 flex justify-content-center align-items-center w-full">
                        <div class="flex flex-column w-full">
                            <div class="gird">
                                <div class="col-4">
                                    <div class="card">
                                        <Chart style="min-height: 25vh;" type="line" :data="chartData" :options="chartOptions" />
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
                                                <DataTable v-model:expandedRows="expandedSalesLogOrderItemComponents" :value="slotProps.data.Items">
                                                    <Column expander style="width: 5rem" />
                                                    <Column sortable field="itemname" header="Name"></Column>
                                                    <Column sortable field="cost" header="Cost"></Column>
                                                    <Column sortable field="sale_price" header="Sale"></Column>
                                                    <Column sortable field="profit" header="Profit">
                                                        <template #body="slotProps">
                                                            <div :style="`${ (slotProps.data.sale_price - slotProps.data.cost) > 0 ? 'color:green' : 'color:red' }`">{{ slotProps.data.sale_price - slotProps.data.cost }}</div>
                                                        </template>
                                                    </Column>
                                                    <template #expansion="slotProps">
                                                        <DataTable :value="slotProps.data.Components">
                                                            <Column sortable field="componentname" header="Component Name"></Column>
                                                            <Column sortable field="cost" header="Cost"></Column>
                                                        </DataTable>        
                                                    </template>
                                                </DataTable>
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

const sales_log = ref([])
const expandedSalesLogRows = ref([])
const expandedSalesLogOrderItems = ref([])
const expandedSalesLogOrderItemComponents = ref([])


const chartData = ref();
const chartOptions = ref();


const chartLabels = ref([])
const chartSales = ref([])
        
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
                borderColor: documentStyle.getPropertyValue('--cyan-500')
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
            per_day_log[`${year}-${month}-${day}`].orders.push(log)
            per_day_log[`${year}-${month}-${day}`].cost += log.cost
            per_day_log[`${year}-${month}-${day}`].sales += log.sale_price
        })


        for (var day in per_day_log){
            sales_log.value.push(per_day_log[day])
            chartSales.value.push(per_day_log[day].sales)
        }


        chartData.value = setChartData();
        chartOptions.value = setChartOptions();
    })
    .catch(error => {
        console.log(error)
    })
}

loadSales()
</script>