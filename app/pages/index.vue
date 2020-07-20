<template lang="pug">
    v-layout(column,justify-center,align-center)
        v-flex(xs12 sm8 md6)
            v-card
                v-card-title.headline
                    | Dashboard
                v-card-text
                    p
                        | Vuetify is a progressive Material Design component
                        | framework for Vue.js. It was designed to empower
                        | developers to create amazing applications.
                    p
                        | Thank you for developing with Vuetify and I look forward
                        | to bringing more exciting features in the future.
                    .text-xs-right
                        em
                            small — John Leider
                    hr.my-3
                    a(href="https://nuxtjs.org/" target="_blank" rel="noopener noreferrer")
                        | Nuxt Documentation
                    br
                    a(href="https://github.com/nuxt/nuxt.js" target="_blank" rel="noopener noreferrer")
                        | Nuxt GitHub
</template>

<script>
import Logo from '~/components/Logo.vue'
import VuetifyLogo from '~/components/VuetifyLogo.vue'
import NavDrawer from '~/components/NavDrawer.vue'

export default {
    components: {
        Logo,
        VuetifyLogo,
        NavDrawer
    }
}
</script>
