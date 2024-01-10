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
import { useRouter } from "vue-router";
import { onMounted } from "vue";
import ListProgramsPage from "@/components/ListProgramsPage.vue";
const router = useRouter();
onMounted(() =>
  fetch("http://localhost:3333/api/processes/is-working").then((a) => {
    if (a.ok) {
      a.json().then((t) => {
        if (t["is-working"]) {
          router.push({
            name: "OvenRun",
            params: { programName: t["program-name"] },
          });
        }
      });
    }
  })
);
function buttonClickedHandler(name) {
  console.log(name);
  router.push({ name: "OvenRun", params: { programName: name } });
}
</script>
