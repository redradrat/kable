<template lang="pug">
    v-navigation-drawer(v-model='drawer', :mini-variant='miniVariant', absolute, dark)
        v-list.py-0(dense='', nav='')
            v-list-item(two-line='', :class="miniVariant && 'px-0'")
                v-list-item-avatar
                    img(src='https://randomuser.me/api/portraits/men/81.jpg')
                v-list-item-content
                    v-list-item-title Application
                    v-list-item-subtitle Subtext
            v-divider
            v-list-item(v-for='item in items', :key='item.title', link='')
                v-list-item-icon
                    v-icon {{ item.icon }}
                v-list-item-content
                    v-list-item-title {{ item.title }}
        template(v-slot:append)
            v-list-item(@click.stop='miniVariant = !miniVariant; if (collapseIcon == "mdi-chevron-left") { collapseIcon = "mdi-chevron-right"} else { collapseIcon = "mdi-chevron-left"}')
                v-list-item-icon
                    v-icon {{ collapseIcon }}
                v-list-item-content
                    span.text-uppercase Collapse
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
    data() {
        return {
            drawer: true,
            collapseIcon: "mdi-chevron-left",
            items: [
                { title: 'Dashboard', icon: 'mdi-view-dashboard' },
                { title: 'Photos', icon: 'mdi-image' },
                { title: 'About', icon: 'mdi-help-box' }
            ],
            miniVariant: false
        }
    }
})
</script>
