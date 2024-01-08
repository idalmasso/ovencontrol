<template>
  <v-container>
    <v-row
      ><v-btn color="blue" @click="$router.push({ name: 'Home' })"
        >Indietro</v-btn
      ></v-row
    >
  </v-container>
  <ListProgramsPage @buttonClicked="buttonClickedHandler" />
</template>

<script setup>
import ListProgramsPage from "@/components/ListProgramsPage.vue";

function buttonClickedHandler(name) {
  console.log(name)
}
</script>
