<template>
    <v-container>
      <v-row justify="space-between"
        >
        <v-col cols="3">
      <v-btn color="blue" @click="$router.push({ name: 'ListProgramConfigurations' })"
          >Indietro</v-btn>
        </v-col>
      </v-row>
      
      
    </v-container>
    <EditProgramPage :programName="props.programName" />
</template>

<script setup>
import {defineProps} from "vue"
import EditProgramPage from "@/components/EditProgramPage.vue";
const props =defineProps({programName:String})
console.log(props)
</script>
  