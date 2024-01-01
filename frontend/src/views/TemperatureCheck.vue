<template>
  <v-container>
    <v-row
      ><v-btn color="blue" @click="$router.push({ name: 'Tests' })"
        >Indietro</v-btn
      ></v-row
    >
    <v-row text-center align-content="center" justify="center">
      <v-card class="mx-auto" width="400" title="Temperatura forno">
        <v-card-text class="py-0">
          <v-row align="center" no-gutters>
            <v-col cols="2">
              <span class="text-h3 font-weight-bold">{{ temp }}</span>
            </v-col>
            <v-col>
              <span class="text-h3 font-weight-bold">&deg;C </span>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref } from "vue";
import { onBeforeUnmount } from "vue";
const temp = ref(0);

const getTemp = () => {
  fetch("http://localhost:3333/api/processes/get-temperature").then((a) => {
    if (a.ok) {
      a.json().then((t) => {
        temp.value = t.Temperature;
      });
    }
  });
};
const timer = setInterval(getTemp, 1000);

onBeforeUnmount(() => clearInterval(timer));
//
</script>
