<template>
  <v-container>
    <v-container>
      <h2>Controller</h2>
      <v-text-field
        label="kp rampa"
        v-model="configData.config.controller['kp-ramp']"
        type="number"
      >
      </v-text-field>
      <v-text-field
        label="ki rampa"
        v-model="configData.config.controller['ki-ramp']"
        type="number"
      >
      </v-text-field>
      <v-text-field
        label="kd rampa"
        v-model="configData.config.controller['kd-ramp']"
        type="number"
      >
      </v-text-field>
      <v-text-field
        label="kp stazionamento"
        v-model="configData.config.controller['kp-maintain']"
        type="number"
      >
      </v-text-field>
      <v-text-field
        label="ki stazionamento"
        v-model="configData.config.controller['ki-maintain']"
        type="number"
      >
      </v-text-field>
      <v-text-field
        label="kd stazionamento"
        v-model="configData.config.controller['kd-maintain']"
        type="number"
      >
      </v-text-field>
      <v-text-field
        label="Tempo di controllo secondi"
        v-model="configData.config.controller['step-time']"
        type="number"
      >
      </v-text-field>
      <v-text-field
        label="Intervallo per salvataggio dati secondi"
        v-model="configData.config.controller['save-time']"
        type="number"
      >
      </v-text-field>
      <v-text-field
        label="Cartella salvataggio run interna"
        v-model="configData.config.controller['saved-run-folder']"
      >
      </v-text-field>
      <v-text-field
        label="Usb device path"
        v-model="configData.config.controller['usb-path']"
      >
      </v-text-field>
      <v-text-field
        label="Usb nome cartella"
        v-model="configData.config.controller['usb-save-folder-name']"
      >
      </v-text-field>
    </v-container>
    <v-container>
      <h2>Server</h2>
      <v-text-field
        label="Temperatura rampa di test"
        v-model="configData.config.server['test-ramp-temperature']"
        type="number"
      >
      </v-text-field>
      <v-text-field
        label="Tempo minuti rampa di test"
        v-model="configData.config.server['test-ramp-time-minutes']"
        type="number"
      >
      </v-text-field>
    </v-container>
    <v-container>
      <h2>Forno</h2>
      <v-text-field
        label="Lunghezza"
        v-model="configData.config.oven['length']"
        type="number"
      >
      </v-text-field>
      <v-text-field
        label="Larghezza"
        v-model="configData.config.oven['width']"
        type="number"
      >
      </v-text-field>
      <v-text-field
        label="Altezza"
        v-model="configData.config.oven['height']"
        type="number"
      >
      </v-text-field>
      <v-text-field
        label="Capacita termica"
        v-model="configData.config.oven['thermal-capacity']"
        type="number"
      >
      </v-text-field>
      <v-text-field
        label="Peso"
        v-model="configData.config.oven['weight']"
        type="number"
      >
      </v-text-field>
      <v-text-field
        label="Potenza massima in Watt"
        v-model="configData.config.oven['max-power']"
        type="number"
      >
      </v-text-field>
    </v-container>

    <v-container>
      <v-btn color="blue" @click="saveData">Save</v-btn>
    </v-container>
  </v-container>
</template>

<script setup>
import { useRouter } from "vue-router";
import { onMounted, reactive } from "vue";
import { useAppStore } from "@/store/app";
const store = useAppStore();
const router = useRouter();
const configData = reactive({
  config: {
    server: {
      "distribution-directory": "../../frontend/dist",
      port: "3333",
      "oven-program-folder": "./programs",
      "test-ramp-temperature": 1000,
      "test-ramp-time-minutes": 1,
    },
    oven: {
      length: "0.11",
      height: "0.08",
      width: "0.23",
      "insulation-widths": [0.075, 0.05],
      "thermal-conductivities": [0.3, 0.14],
      "thermal-capacity": "900",
      weight: "24",
      "max-power": "10000",
    },
    controller: {
      "kp-ramp": "0.1",
      "ki-ramp": "0.05",
      "kd-ramp": "0.001",
      "kp-maintain": "0.01",
      "ki-maintain": "0.0001",
      "kd-maintain": "0.0001",
      "step-time": "1",
      "save-time": "60",
      "saved-run-folder": "./runs",
      "usb-path": "/dev/sda1",
      "usb-save-folder-name": "ovenruns",
    },
  },
});
function saveData() {
  fetch("http://localhost:3333/api/configuration/oven-config", {
    method: "POST",
    body: JSON.stringify(configData.config),
  }).then((response) => {
    if (response.ok) {
      router.back();
    } else {
      response.json().then((data) => store.setAPIError(data["Error"]));
    }
  });
}
onMounted(() => {
  fetch("http://localhost:3333/api/configuration/oven-config").then(
    (response) => {
      if (response.ok) {
        response.json().then((data) => {
          console.log(data);
          configData.config = data;
        });
      } else {
        response.json().then((data) => store.setAPIError(data["Error"]));
      }
    }
  );
});
</script>
