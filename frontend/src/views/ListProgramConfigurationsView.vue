<template>
  <v-container>
    <v-row justify="space-between"
      >
      <v-col cols="3">
    <v-btn color="blue" @click="$router.push({ name: 'Home' })"
        >Indietro</v-btn>
      </v-col>
      <v-col cols="3">
      <v-btn color="green" @click="buttonClickedHandler('')">Nuovo programma</v-btn>
    </v-col>
  
    </v-row>
    
    
  </v-container>
  <ListProgramsPage @buttonClicked="buttonClickedHandler" />
</template>

<script setup>
import { useRouter } from "vue-router";
import ListProgramsPage from "@/components/ListProgramsPage.vue";
const router=useRouter()
function buttonClickedHandler(name) {
  if (name==""){
    router.push({name:"NewOvenProgram"})
  } else {
    router.push({name:"EditOvenProgram", params:{"programName": name}})
  }
}
</script>
