<template>
  <v-container>
    <v-row
      ><v-btn color="blue" @click="$router.push({ name: 'Tests' })"
        >Indietro</v-btn
      ></v-row
    >
    <v-container class="fill-height">
      <v-responsive class="align-center text-center fill-height">
        Read: {{ temp }}
      </v-responsive>
    </v-container>
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
