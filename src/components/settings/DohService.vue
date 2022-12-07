<template>
    <div class="grid grid-cols-1 gap-6 space-y-4">
        <label class="block">
            <span class="text-gray-700 block">Current DoH service</span>
            <input type="text" disabled v-model="currentDohServiceUrl" class="disabled bg-slate-50 mt-1 w-1/3 rounded-md border-gray-700 shadow-sm">
        </label>    
        <label class="block">
            <span class="block text-gray-700">DoH Service</span>
            <input type="text" v-model="pendingDohServiceUrl" class="mt-1 w-1/3 rounded-md border-gray-700 shadow-sm">
            <button @click="testDohService"
                class="bg-blue-500 hover:bg-blue-600 px-5 active:bg-blue-700 py-2 ml-4 leading-5 text-base font-semibold text-white rounded-md">Test</button>
            <button @click="saveDohService"
                class="bg-blue-500 hover:bg-blue-600 px-5 active:bg-blue-700 py-2 ml-4 leading-5 text-base font-semibold text-white rounded-md">Save</button>
        </label>
        <p v-if="dohServiceUpdated" class="text-teal-600">URL of DNS-over-HTTPS has been changed.</p>
    </div>
</template>

<script>
import axios from 'axios';

export default {
    data() {
        return {
            currentDohServiceUrl: '',
            pendingDohServiceUrl: '',
            shadowPendingDohServiceUrl: '',
            dohServiceUpdated: false,
        }
    },
    methods: {
        setShadowPendingDohServiceUrl() {
            this.shadowPendingDohServiceUrl = this.pendingDohServiceUrl;
            this.allTestPassed = false;
            this.dohServiceUpdated = false;
        },
        loadCurrentDohServiceUrl() {
            axios.get('/api/settings/doh_service')
            .then((r) => {
                this.currentDohServiceUrl = r.data.data.server;
            })
        },
        testDohService() { 
            this.setShadowPendingDohServiceUrl();
        },
        saveDohService() {
            this.testDohService();
            axios.post('/api/settings/doh_service',
                { server: this.shadowPendingDohServiceUrl}
            ).then(() => {
                this.dohServiceUpdated = true;
                this.loadCurrentDohServiceUrl();
            });
        }
    },
    mounted() {
        this.loadCurrentDohServiceUrl();
    }
}
</script>