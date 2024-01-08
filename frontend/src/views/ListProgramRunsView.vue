<template>
  <v-container>
    <v-row
      ><v-btn color="blue" @click="$router.push({ name: 'Home' })"
        >Indietro</v-btn
      ></v-row
    >
  </v-container>
  <ListProgramPage @buttonClicked="buttonClickedHandler" />
</template>

<script setup>
import { defineEmits } from "vue";
import { onBeforeMount } from "vue";
import ListProgramPage from "@/components/ListProgramPage.vue";
const emits = defineEmits({ programSelected: String });
onBeforeMount(() =>
  fetch("http://localhost:3333/api/processes/is-working").then((a) => {
    if (a.ok) {
      a.json().then((t) => {
        if (t["is-working"]) {
          this.$router.push({
            name: "OvenRun",
            programName: t["program-name"],
          });
        }
      });
    }
  })
);
function buttonClickedHandler(name) {
  emits("programSelected", name);
}
</script>
