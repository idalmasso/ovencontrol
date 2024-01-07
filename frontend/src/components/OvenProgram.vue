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
import { ref, reactive } from "vue";
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
import { onMounted } from "vue";
ChartJS.register(LinearScale, PointElement, LineElement, Tooltip, Legend);
const temp = ref(0);
const tempExpected = ref(0);
const isWorking = ref(false);
var temperatureData = [];
var temperatureExpected = [];
const chart = ref(null);
const chartData = {
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
const updatePlotData = () => {
  chartData["datasets"][0]["data"] = temperatureData;
  chartData["datasets"][1]["data"] = temperatureExpected;
  if (chart.value !== null) {
    chart.value.chart.data = chartData;
    chart.value.chart.update();
  }
};
const getTemp = () => {
  fetch("http://localhost:3333/api/processes/get-temperatures-process").then(
    (a) => {
      if (a.ok) {
        a.json().then((t) => {
          temp.value = t["oven-temperature"];
          tempExpected.value = t["expected-temperature"];
          temperatureData.push({
            x: t["time-seconds"],
            y: t["oven-temperature"],
          });
          temperatureExpected.push({
            x: t["time-seconds"],
            y: t["expected-temperature"],
          });
          updatePlotData();
        });
      }
    }
  );
};
function IsWorkingEnabler() {
  return fetch("http://localhost:3333/api/processes/is-working").then(
    async (a) => {
      if (a.ok) {
        await a.json().then((t) => {
          isWorking.value = t["is-working"];
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
    }
  );
}
var timer = 0;
const timerWorking = setInterval(IsWorkingEnabler, 1000);

function TryStartTest() {
  if (!isWorking.value) {
    fetch("http://localhost:3333/api/processes/test-ramp", {
      method: "POST",
    }).then((a) => {
      if (a.ok) {
        temperatureData = [];
        temperatureExpected = [];
        getTemp();
        timer = setInterval(getTemp, 10000);
      }
    });
  }
}

const chartOptions = reactive({
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    x: { type: "linear", min: 0 },
    y: { type: "linear", min: 0, max: 1000 },
  },
});

onMounted(() => {
  IsWorkingEnabler().then(() => {
    if (isWorking.value) {
      fetch(
        "http://localhost:3333/api/processes/get-actual-process-data?step=10"
      ).then((a) => {
        if (a.ok) {
          a.json().then((data) => {
            temperatureExpected = data.map((row) => {
              return {
                x: row["seconds-from-start"],
                y: row["desired-temperature"],
              };
            });
            temperatureData = data.map((row) => {
              return {
                x: row["seconds-from-start"],
                y: row["oven-temperature"],
              };
            });
            temp.value = data[data.length - 1]["oven-temperature"];
            tempExpected.value = data[data.length - 1]["desired-temperature"];
            updatePlotData();
          });
        }
      });
    }
  });
});
onBeforeUnmount(() => {
  clearInterval(timer);
  clearInterval(timerWorking);
});
//
</script>
