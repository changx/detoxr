<template>
    <div class="grid grid-cols-1 gap-6 space-y-4">
        <label class="block">
            <span class="text-gray-700 block">Current Honeypot</span>
            <input type="text" disabled v-model="currentHoneypot" class="disabled bg-slate-50 mt-1 w-1/3 rounded-md border-gray-700 shadow-sm">
        </label>
        <label class="block">
            <span class="text-gray-700 block">Honeypot IP/domain</span>
            <div>
                <input type="text" v-model="pendingHoneypot" class="mt-1 rounded-md border-gray-700 shadow-sm w-1/3">
                <button @click="testPendingHoneypot"
                    class="bg-blue-500 hover:bg-blue-600 px-5 active:bg-blue-700 py-2 ml-4 leading-5 text-base font-semibold text-white rounded-md">Test</button>
                <button @click="savePendingHoneypot"
                    class="bg-blue-500 hover:bg-blue-600 px-5 active:bg-blue-700 py-2 ml-4 leading-5 text-base font-semibold text-white rounded-md">Save</button>
            </div>
            <div class="text-slate-400">
                Fill an IP address: NOT allocated to China ISPs, NOT serving as forward name server.
            </div>
        </label>
        <p v-for="r in testResults" :key="r.host">
            {{r.host}} ... {{r.return}}
        </p>
        <p v-if="allTestPassed" class="text-teal-600">All tests passed.</p>
        <p v-if="honeypotUpdated" class="text-teal-600">Honeypot IP changed.</p>
    </div>
</template>

<script>
import axios from 'axios';
export default {
    data() {
        return {
            currentHoneypot: '',
            pendingHoneypot: '202.12.27.33',
            shadowPendingHoentypot: '202.12.27.33:53',
            allTestPassed: false,
            honeypotUpdated: false,
            testResults: [],
            testDomains: [
                {host: 'twitter.com', shouldAnswer: true},
                {host: 'baidu.com', shouldAnswer: false},
                {host: 'facebook.com', shouldAnswer: true},
                {host: 'qq.com', shouldAnswer: false}
            ],
        }
    },
    methods: {
        setShadowPendingHoneypot() {
            if (!this.pendingHoneypot.endsWith(':53')) {
                this.shadowPendingHoentypot = this.pendingHoneypot + ':53';
            } else {
                this.shadowPendingHoentypot = this.pendingHoneypot;
            }
            this.allTestPassed = false;
            this.honeypotUpdated = false;
        },
        async testPendingHoneypot() {
            this.setShadowPendingHoneypot();
            this.testResults = [];

            let tests = [];
            this.testDomains.map((h) => {
                this.testResults.push({ host: h.host, return: '?', success: false });
                tests.push(this.runTest(h));
            });

            return Promise.all(tests).then(() => {
                let allPassed = true;
                this.testResults.map((s) => {
                    if (!s.success) {
                        allPassed = false;
                    }
                });
                this.allTestPassed = allPassed;
            })
        },
        ok(shouldAnswer, error, answer) {
            if (shouldAnswer && !error && answer) {
                return 'ok';
            } else if (!shouldAnswer && (error || !answer)) {
                return 'ok';
            }
            return 'failed';
        },
        async runTest(t) {
            return axios.post('/api/settings/honeypot/try',
                { server: this.shadowPendingHoentypot, host: t.host }
            ).then((r) => {
                this.testResults.map((s) => {
                    if (s.host === r.data.data.host) {
                        s.return = this.ok(t.shouldAnswer, r.data.data.error, r.data.data.answer);
                        s.success = s.return === 'ok';
                    }
                });
            });
        },
        savePendingHoneypot() {
            this.testPendingHoneypot().then(() => {
                if (this.allTestPassed) {
                    axios.post('/api/settings/honeypot',
                        { server: this.shadowPendingHoentypot }
                    ).then(() => {
                        this.honeypotUpdated = true;
                        this.loadCurrentHoneypot();
                    });
                }
            });
        },
        loadCurrentHoneypot() {
            axios.get('/api/settings/honeypot').then((r) => {
                this.currentHoneypot = r.data.data.server.replace(':53', '');
            })
        }
    },
    mounted() {
        this.loadCurrentHoneypot();
    },
    unmounted() {
    }
}
</script>