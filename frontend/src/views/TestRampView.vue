<template>
  <v-container>
    <v-row
      ><v-btn
        v-if="isWorking"
        color="blue"
        @click="$router.push({ name: 'ListTests' })"
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
      <v-card class="mx-auto" width="400" title="Temperatura aspettata">
        <v-card-text class="py-0">
          <v-row align="center" no-gutters>
            <span class="text-h3 font-weight-bold">{{ tempExpected }}</span>
            <span class="text-h3 font-weight-bold">&deg;C </span>
          </v-row>
        </v-card-text>
      </v-card>
    </v-row>
    <v-row
      ><v-btn :disabled="isWorking" color="blue" @click="TryStartTest"
        >Test</v-btn
      ></v-row
    >
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
const tempExpected = ref(0);
const isWorking = ref(false);
var temperatureData = [];
var temperatureExpected = [];
var i = 0;
const chart = ref(null);
const getTemp = () => {
  fetch("http://localhost:3333/api/processes/get-temperatures-process").then(
    (a) => {
      if (a.ok) {
        a.json().then((t) => {
          temp.value = t.Temperature;
          tempExpected.value = t.ExpectedTemperature;
          temperatureData.push({
            x: i,
            y: t.Temperature,
          });
          temperatureExpected.push({
            x: i,
            y: t.ExpectedTemperature,
          });
          if (chart.value !== null) {
            chart.value.chart.update();
          }
          i += 10;
        });
      }
    }
  );
};
var timer = 0;
const timerWorking = setInterval(
  () =>
    fetch("http://localhost:3333/api/processes/is-working").then((a) => {
      if (a.ok) {
        a.json().then((t) => {
          isWorking.value = t["IsWorking"];
          if (!isWorking.value) {
            if (timer !== 0) {
              clearInterval(timer);
              timer = 0;
            }
          } else {
            if (timer === 0) {
              timer = setInterval(getTemp, 10000);
            }
          }
        });
      }
    }),
  1000
);

function TryStartTest() {
  if (!isWorking.value) {
    fetch("http://localhost:3333/api/processes/test-ramp", {
      method: "POST",
    }).then((a) => {
      if (a.ok) {
        timer = setInterval(getTemp, 10000);
      }
    });
  }
}

const chartOptions = reactive({
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    x: [
      {
        type: "linear",
      },
    ],
    y: [{ type: "linear", min: -1, max: 1000 }],
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
      {
        label: "Desiderata",
        borderColor: "blue",
        backgroundColor: "blue",
        showLine: true,
        data: temperatureExpected,
      },
    ],
  };
});
onBeforeUnmount(() => {
  clearInterval(timer);
  clearInterval(timerWorking);
});
//
</script>
