<template>
  <!-- <div>
    {{ rows[5] }}
  </div> -->
  <q-table
    flat
    bordered
    :title="matrixTitle"
    :rows="rows"
    :columns="(columns as any)"
    row-key="name"
    hide-bottom
    :separator="'cell'"
    v-model:pagination="pagination"
  >
    <template v-slot:body="props">
      <q-tr :props="props">
        <q-td v-for="col in props.cols" :key="col.name" :props="props">
          {{ props.row[col.name] }}
        </q-td>
      </q-tr>
    </template>
  </q-table>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';

const props = defineProps<{
  matrix: number[][];
  matrixTitle?: string;
  columnsLabels: string[];
  rowsLabels?: string[];
}>();


const rowsLabelsComputed = computed(
  () => props.rowsLabels || props.columnsLabels
);

const columns = ref<Array<Record<string, unknown>>>([]);
const rows = ref<Array<Record<string, unknown>>>([]);
const pagination = ref({ rowsPerPage: 0 }); // Show all rows

onMounted(() => {
  const columnWidth = 100 / (props.matrix[0].length + 1);
  const style = `width: ${columnWidth}%;`;

  columns.value.push({
    name: 'name',
    label: '',
    style,
  });

  props.columnsLabels.forEach((label) => {
    columns.value.push({
      name: label,
      label: label,
      style,
    });
  });

  if (props.columnsLabels.length === props.matrix[0].length) {
    console.log(props.matrix)
    props.matrix.forEach((row, i) => {
      const newRow: Record<string, unknown> = {
        name: rowsLabelsComputed.value[i],
      };

      row.forEach((value, j) => {
        newRow[props.columnsLabels[j]] = value;
      });
      console.log(row)
      rows.value.push(newRow);
    });
  }
});
</script>

<style lang="scss" scoped></style>
