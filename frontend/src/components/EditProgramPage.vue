<template>
  <v-container>
    <h2>Programma: {{ ovenProgramContainer.ovenProgram["name"] }}</h2>

    <v-form @submit.prevent="saveData">
      <v-text-field
        label="Nome programma"
        v-model="ovenProgramContainer.ovenProgram['name']"
      ></v-text-field>
      <v-label for="color">Colore icona </v-label>
      <v-btn :color="ovenProgramContainer.ovenProgram['icon-color']" id="color">
        <v-menu activator="parent" offset-y>
          <v-color-picker
            v-model="ovenProgramContainer.ovenProgram['icon-color']"
            hide-canvas
            hide-inputs
            show-swatches
            class="mx-auto"
          >
          </v-color-picker>
        </v-menu>
      </v-btn>
      <v-text-field
        label="Temperatura chiusura aria"
        v-model="ovenProgramContainer.ovenProgram['air-closed-at-degrees']"
        type="number"
      >
      </v-text-field>
      <v-row>
        <h3>Punti</h3>
      </v-row>
      <v-row ga-3 class="pa-2 ma-2">
        <v-btn @click="addPoint(point)" color="blue">Aggiungi punto</v-btn>
      </v-row>

      <v-row
        ga-3
        v-for="point in ovenProgramContainer.ovenProgram['points']"
        :key="point"
      >
        <v-col cols="3">
          <v-text-field label="Segmento" v-model="point['segment-name']">
          </v-text-field>
        </v-col>
        <v-col cols="2">
          <v-text-field
            label="Temperatura"
            v-model="point['temperature']"
            type="number"
          >
          </v-text-field>
        </v-col>
        <v-col cols="2">
          <v-text-field
            label="Tempo minuti"
            v-model="point['time-minutes']"
            type="number"
          >
          </v-text-field>
        </v-col>
        <v-col cols="3">
          <v-checkbox
            label="Ricomincia se offline"
            v-model="point['restart-from-last-ascending-ramp']"
            type="checkbox"
          >
          </v-checkbox>
        </v-col>
        <v-col cols="2">
          <v-text-field
            label="Minuti no restart"
            v-model="point['time-after-no-restart-minutes']"
            type="number"
          >
          </v-text-field>
        </v-col>
        <v-col cols="3">
          <v-btn class="pa-2 ma-2" @click="removePoint(point)" color="red"
            >Elimina</v-btn
          >
        </v-col>
        <v-divider></v-divider>
      </v-row>
      <v-row justify="center" style="height: 200px">
        <Scatter :data="chartData" :options="chartOptions" ref="chart" />
      </v-row>
      <v-row no-gutters class="pa-2 ma-2">
        <v-btn type="submit" block color="success">Salva</v-btn>
      </v-row>
    </v-form>
    <v-btn @click="deleteProgram" color="red">ELIMINA</v-btn>
  </v-container>
</template>
<script setup>
import { useRouter } from "vue-router";
import { defineProps, ref, reactive, onMounted, watch } from "vue";
import { Scatter } from "vue-chartjs";
import { useAppStore } from "@/store/app";
import {
  Chart as ChartJS,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip,
  Legend,
} from "chart.js";

const store = useAppStore();
ChartJS.register(LinearScale, PointElement, LineElement, Tooltip, Legend);
var temperatureData = [{ x: 0, y: 0 }];
const props = defineProps({ programName: String });
const router = useRouter();

const ovenProgramContainer = reactive({
  ovenProgram: {
    name: "",
    "icon-color": "#00AAAA",
    points: [],
    "air-closed-at-degrees": "30",
  },
});
const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    x: { type: "linear", min: 0 },
    y: { type: "linear", min: 0, max: 1400 },
  },
};
const chartData = {
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
const chart = ref(null);
function recalc_chart(points) {
  temperatureData = [{ x: 0, y: 0 }];
  var t = 0;
  for (const v of points) {
    console.log(v);
    t += parseFloat(v["time-minutes"]);
    temperatureData.push({ x: parseFloat(t), y: parseFloat(v["temperature"]) });
  }
  chartData["datasets"][0]["data"] = temperatureData;
  if (chart.value !== null) {
    chart.value.chart.data = chartData;
    chart.value.chart.update();
  }
  console.log(temperatureData);
}

watch(ovenProgramContainer.ovenProgram.points, (newValue) => {
  recalc_chart(newValue);
});
const saveData = () => {
  if (ovenProgramContainer.ovenProgram.name !== "") {
    fetch("http://localhost:3333/api/configuration/programs", {
      method: "POST",
      body: JSON.stringify(ovenProgramContainer.ovenProgram),
    }).then((a) => {
      if (a.ok) {
        router.push({ name: "ListProgramConfigurations" });
      } else {
        a.json().then((data) => store.setAPIError(data["Error"]));
      }
    });
  }
};
const deleteProgram = () => {
  if (ovenProgramContainer.ovenProgram.name !== "") {
    fetch(
      "http://localhost:3333/api/configuration/programs/" +
        ovenProgramContainer.ovenProgram.name,
      { method: "DELETE" }
    ).then((a) => {
      if (a.ok) {
        router.push({ name: "ListProgramConfigurations" });
      } else {
        a.json().then((data) => store.setAPIError(data["Error"]));
      }
    });
  }
};

const removePoint = (point) => {
  ovenProgramContainer.ovenProgram["points"] = ovenProgramContainer.ovenProgram[
    "points"
  ].filter((a) => a !== point);
};

const addPoint = () => {
  ovenProgramContainer.ovenProgram["points"].push({
    temperature: "0.0",
    "time-minutes": "0.0",
  });
};

onMounted(() => {
  if (
    props.programName !== "" &&
    props.programName !== undefined &&
    props.programName !== null
  ) {
    fetch(
      "http://localhost:3333/api/configuration/programs/" + props.programName
    ).then((response) => {
      if (response.ok) {
        response.json().then((data) => {
          ovenProgramContainer.ovenProgram = data;
          recalc_chart(data.points);
        });
      } else {
        response.json().then((data) => store.setAPIError(data["Error"]));
        router.back();
      }
    });
  }
});
</script>
