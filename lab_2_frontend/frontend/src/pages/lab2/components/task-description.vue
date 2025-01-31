<template>
  <div>
    <banner-component class="text-primary text">
      <template #title>
        <div>{{ variant?.variant }}</div>
      </template>
      <div>
        <span class="text-main"> Взять выборку №1: </span> Усреднять на уровне баллов
      </div>
      <div>
        <ol>
          <i>Вычисления:</i>
          <li>Рассчитать моду и медиану</li>
          <li>Рассчитать консолидированное мнение экспертов 4-мя способами(среднее арифметическое, среднее геометрическое, среднее гармоническое,среднее по Байесу)</li>
          <li>Рассчитать дисперсию и среднекв. отклонение на основе средних,
            полученных по каждому из 4-х способов</li>
          <li>Рассчитать отклонение других усредненных результатов в % от среднего
            по Байесу</li>
        </ol>
      </div>
      <div>
        <span class="text-main"> Взять выборки №1 и №2. </span>
        <div><span class="text-main">По обеим выборкам произвести вычисление средних весов:</span></div>
        <div><span class="text-main">При этом</span> усреднять на уровне баллов, рассчитывать консолидированное мнение экспертов по Байесу</div>
      </div>

      <br>
      <div><span class="text-main">По итогам:</span></div>
      <div>
        <ol>
          <li>Упорядочить представленные шесть критериев по убыванию по
            параметру «отличие весов, выставленных экспертами разных выборок».</li>
          <li>Логически объяснить полученные результаты.</li>
        </ol>
      </div>
    </banner-component>
  </div>
</template>

<script lang="ts" setup>
import BannerComponent from 'src/components/BannerComponent.vue';
import { Lab1a } from 'src/models/lab/lab1a';
import { Lab2 } from 'src/models/lab/lab2';

const props = defineProps<{
  variant: Lab2.Variant | null;
  // variantNumber?: number;
}>();
console.log('variant:', props.variant)
</script>

<style lang="scss" scoped>
.text {
  font-size: 18px;
  font-weight: 500;
}

.text-main {
  font-weight: 700;
}

table {
  border-collapse: collapse;
  border-spacing: 0;
  border: 2px solid var(--q-primary);
}

th {
  border: 1px solid var(--q-primary);
  padding: 8px;
  text-align: left;
  text-transform: capitalize;
}

td {
  border: 1px solid var(--q-primary);
  padding: 8px;
  text-align: left;
}

.alternative {
  display: flex;
  flex-direction: column;
}
</style>
