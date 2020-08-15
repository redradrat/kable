<template lang="pug">
    v-app(dark)
        AppNav
        v-content
            v-container
                nuxt
        v-footer(:absolute='!fixed', app)
            span © {{ new Date().getFullYear() }}
</template>

<script lang="ts">
import Vue from 'vue'
import AppNav from '~/components/AppNav.vue'

export default Vue.extend({
    components: {
        AppNav
        },
      data () {
    return {
      fixed: false,
      title: 'kable'
    }
  }
})
</script>
