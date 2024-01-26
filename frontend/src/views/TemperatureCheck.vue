<template>
  <v-container>
    <v-row
      ><v-btn color="blue" @click="$router.push({ name: 'ListTests' })"
        >Indietro</v-btn
      >
      <v-text-field  v-if="!isWorking"  label="Power % (0-1)" 
            v-model="power"
            type="number">
        </v-text-field>
      <v-btn color="red"  v-if="!isWorking" @click="powerOvenOneMinute">Power one min</v-btn></v-row
    >
    <v-row justify="center" style="height: 200px ">
      <Scatter :data="chartData" :options="chartOptions" ref="chart" width="500px" />
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
const isWorking = ref(false);
var temperatureData = [];
var i = 0;
const chart = ref(null);
const power = ref(0.0)
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
        if (temperatureData.length > 600) {
          temperatureData.shift();
        }
      });
    }
  });
};

function IsWorkingEnabler() {
  return fetch("http://localhost:3333/api/processes/is-working").then(
    async (a) => {
      if (a.ok) {
        await a.json().then((t) => {
          isWorking.value = t["is-working"];
        });
      }
    }
  );
};
const timerWorking = setInterval(IsWorkingEnabler, 1000);
const powerOvenOneMinute = ()=>{
  fetch("http://localhost:3333/api/processes/set-power-one-minute",{method:"POST", body:JSON.stringify({'power':power.value})}).then((a) => {
    if (!a.ok) {
      a.json().then(data=>{
        alert(data.error)
      })
    }
  })
}
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
onBeforeUnmount(() => {clearInterval(timer);clearInterval(timerWorking)});
//
</script>
