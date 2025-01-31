import { RouteRecordRaw } from 'vue-router';

const MainLayout = () => import('layouts/MainLayout.vue');
const ErrorPage = () => import('pages/ErrorNotFound.vue');

import AuthPageVue from 'src/pages/auth/AuthPage.vue';
import Lab3 from 'src/pages/lab3/lab-3.vue';
import lab3BVue from 'src/pages/lab3b/lab-3-b.vue';
import lab1A from 'src/pages/lab1a/lab-1-a.vue';
import lab3C from 'src/pages/lab3c/lab-3-c.vue';
import ChooseLab from 'src/pages/chooseLab/ChooseLab.vue';
import lab2 from 'src/pages/lab2/lab-2.vue';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: MainLayout,
    children: [
      { path: '/', component: AuthPageVue },
      {
        path: '/choose-lab',
        name: 'choose-lab',
        component: ChooseLab,
      },
      {
        path: '/lab3-a',
        name: 'lab3-a',
        component: Lab3,
      },
      {
        path: '/lab3-b',
        name: 'lab3-b',
        component: lab3BVue,
      },
      {
        path: '/lab1-a',
        name: 'lab1-a',
        component: lab1A,
      },
      {
        path: '/lab2',
        name: 'lab2',
        component: lab2,
      },
      {
        path: '/lab3-c',
        name: 'lab3-c',
        component: lab3C,
      },
    ],
  },
  {
    path: '/:catchAll(.*)*',
    component: ErrorPage,
  },
];

export default routes;
