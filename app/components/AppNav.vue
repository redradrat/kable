<template lang="pug">
        v-app-bar(height="50px",fixed, app,color="primary")
            .d-inline-flex.flex-grow-0.justify-left
                v-btn(outlined,color="white",to="/",active-class="no-active",nuxt)
                    h2.text-uppercase.white--text(style='font-family:"Arco Typography"') Kapabilities
            .d-inline-flex.flex-grow-1.justify-center
                v-btn.ma-2(text,v-for="item in menuItems", :key="item.label",nuxt,:to='item.location') {{ item.label }}
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
    data: () => ({
        menuItems: [
            { label: 'Dashboard', location: '/' },
            { label: 'Groups', location: '/groups' },
            { label: 'Repository', location: '/repo' }
        ],
        gradient: ['red', 'orange', 'yellow']
    })
})
</script>

<style lang="sass">
.v-btn--active.no-active::before
  opacity: 0 !important
</style>
