<template>
  <v-container>
    <v-row
      ><v-btn color="blue" @click="$router.push({ name: 'Home' })"
        >Indietro</v-btn
      ></v-row
    >
    <v-row class="mx-auto" align="center" justify="center">
      Sei davvero sicuro di voler spegnere? Questo interromperà ogni cottura in
      corso
    </v-row>
    <v-row class="mx-auto" align="center" justify="center"
      ><v-col align="center" justify="center"
        ><v-btn color="green" @click="okClicked">Sì</v-btn></v-col
      ><v-col align="center" justify="center"
        ><v-btn color="red" @click="$router.push({ name: 'Home' })"
          >No</v-btn
        ></v-col
      ></v-row
    >
  </v-container>
</template>

<script setup>
import { useAppStore } from "@/store/app";
const store = useAppStore();

const okClicked = () =>
  fetch("http://localhost:3333/api/power-off", { method: "POST" }).then((a) => {
    if (!a.ok) {
      a.json().then((t) => {
        if (t["Error"] !== null) {
          store.setAPIError(t["Error"]);
        }
      });
    }
  });
</script>
