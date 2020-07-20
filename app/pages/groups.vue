<template lang="pug">
            v-card.ma-2
                v-card-title
                    | Groups
                    v-spacer
                    v-text-field(v-model='search', append-icon='mdi-magnify', label='Search', single-line='', hide-details='')
                v-data-table(:headers='headers', :items='groups', :search='search')
                    template(v-slot:top='')
                            v-toolbar(flat,dark)
                                v-spacer
                                v-dialog(v-model='dialog', max-width='500px')
                                    template(v-slot:activator='{ on, attrs }')
                                        v-btn.mb-2(color='primary', dark='', v-bind='attrs', v-on='on') New Item
</template>

<script lang="ts">
import Vue from 'vue'
import NavDrawer from '~/components/NavDrawer.vue'

export default Vue.extend({
    components: {
        NavDrawer
         },
    data() {
        return {
            search: '',
            headers: [
                {
                    text: 'Name',
                    align: 'start',
                    sortable: false,
                    value: 'name'
                },
                { text: 'Description', value: 'desc' },
            ],
            groups: [
                {
                    name: 'Application Dev',
                    desc: 'The Application Dev stage project'
                },
                {
                    name: 'Application QA',
                    desc: 'The Application QA stage project'
                },
                {
                    name: 'Application Prod',
                    desc: 'The Application Prod stage project'
                },
                {
                    name: 'Application Prod',
                    desc: 'The Application Prod stage project'
                },
                {
                    name: 'Application Prod',
                    desc: 'The Application Prod stage project'
                },
                {
                    name: 'Application Prod',
                    desc: 'The Application Prod stage project'
                },
                {
                    name: 'Application Prod',
                    desc: 'The Application Prod stage project'
                },
                {
                    name: 'Application Staging',
                    desc: 'The Application Staging stage project'
                }
            ]
        }
    }
})
</script>
