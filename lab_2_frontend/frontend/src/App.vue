<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useAuthStore } from './stores/auth';

  const store = useAuthStore();

  onMounted(() => {
    console.log('main mounted')

    const url = new URL(window.location.href);
    console.log('url', url)
    const jwt = url.searchParams.get('jwt');
    console.log(url.searchParams)
    console.log('jwt', jwt)

    if (jwt) {
      console.log('in set jwt')
      // Store JWT in localStorage
      localStorage.setItem('token', jwt);
      url.searchParams.delete('jwt');
      window.history.replaceState({}, document.title, url.pathname + url.search);
      console.log('jwt saved', jwt)

      if (store.tryApplyJwt(jwt)){
        console.log('applied', store.userType)
      } else {
        console.log('failed', store.userType)
      }
      // window.location.reload()
    }
  });
</script>
