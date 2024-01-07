<template>
  <v-container>
    <v-row
      ><v-btn color="blue" @click="$router.push({ name: 'ListTests' })"
        >Indietro</v-btn
      ></v-row
    >
    <v-row style="height: 400px">
      <Scatter :data="chartData" :options="chartOptions" ref="chart" />
    </v-row>
    <v-row text-center align-content="center" justify="center">
      <v-card class="mx-auto" width="400" title="Temperatura forno">
        <v-card-text class="py-0">
          <v-row align="center" no-gutters>
            <span class="text-h3 font-weight-bold">{{ temp }}</span>
            <span class="text-h3 font-weight-bold">&deg;C </span>
          </v-row>
        </v-card-text>
      </v-card>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, computed, reactive } from "vue";
import { onBeforeUnmount } from "vue";
import { Scatter } from "vue-chartjs";
import {
  Chart as ChartJS,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip,
  Legend,
} from "chart.js";
ChartJS.register(LinearScale, PointElement, LineElement, Tooltip, Legend);
const temp = ref(0);
var temperatureData = [];
var i = 0;
const chart = ref(null);
const getTemp = () => {
  fetch("http://localhost:3333/api/processes/get-temperature").then((a) => {
    if (a.ok) {
      a.json().then((t) => {
        temp.value = t["oven-temperature"];
        temperatureData.push({
          x: i,
          y: t["oven-temperature"],
        });
        if (chart.value !== null) {
          chart.value.chart.update();
        }
        i++;
        if (temperatureData.length > 10) {
          temperatureData.shift();
        }
      });
    }
  });
};
const timer = setInterval(getTemp, 1000);
const chartOptions = reactive({
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    x: [
      {
        type: "linear",
      },
    ],
    y: [{ type: "linear" }],
  },
});

const chartData = computed(() => {
  return {
    datasets: [
      {
        label: "Temperatura",
        borderColor: "red",
        backgroundColor: "red",
        showLine: true,
        data: temperatureData,
      },
    ],
  };
});
onBeforeUnmount(() => clearInterval(timer));
//
</script>
