<template>
    <div class="grid grid-cols-1 gap-6 space-y-4">
        <label class="block">
            <span class="text-gray-700 block">Current Local NS</span>
            <input type="text" disabled v-model="currentLocalNS" class="disabled bg-slate-50 mt-1 w-1/3 rounded-md border-gray-700 shadow-sm">
        </label>
        <label class="block">
            <span class="text-gray-700 block">Local NS IP</span>
            <div>
                <input type="text" v-model="pendingLocalNS" class="mt-1 rounded-md border-gray-700 shadow-sm w-1/3">
                <button @click="testPendingLocalNS"
                    class="bg-blue-500 hover:bg-blue-600 px-5 active:bg-blue-700 py-2 ml-4 leading-5 text-base font-semibold text-white rounded-md">Test</button>
                <button @click="savePendingLocalNS"
                    class="bg-blue-500 hover:bg-blue-600 px-5 active:bg-blue-700 py-2 ml-4 leading-5 text-base font-semibold text-white rounded-md">Save</button>
            </div>
        </label>
        <p v-for="r in testResults" :key="r.host">
            {{r.host}} ... {{r.return}}
        </p>
        <p v-if="allTestPassed" class="text-teal-600">All tests passed.</p>
        <p v-if="localNSUpdated" class="text-teal-600">Local NS IP changed.</p>
    </div>
</template>

<script>
import axios from 'axios';
export default {
    data() {
        return {
            currentLocalNS: '',
            pendingLocalNS: '',
            shadowPendingLocalNS: '',
            allTestPassed: false,
            localNSUpdated: false,
            testResults: [],
            testDomains: [
                {host: 'baidu.com', shouldAnswer: true},
                {host: 'qq.com', shouldAnswer: true},
                {host: 'bilibili.com', shouldAnswer: true}
            ],
        }
    },
    methods: {
        setShadowPendingLocalNS() {
            if (!this.pendingLocalNS.endsWith(':53')) {
                this.shadowPendingLocalNS = this.pendingLocalNS+ ':53';
            } else {
                this.shadowPendingLocalNS = this.pendingLocalNS;
            }
            this.allTestPassed = false;
            this.localNSUpdated = false;
        },
        async testPendingLocalNS() {
            this.setShadowPendingLocalNS();
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
            return axios.post('/api/settings/local_ns/try',
                { server: this.shadowPendingLocalNS, host: t.host }
            ).then((r) => {
                this.testResults.map((s) => {
                    if (s.host === r.data.data.host) {
                        s.return = this.ok(t.shouldAnswer, r.data.data.error, r.data.data.answer);
                        s.success = s.return === 'ok';
                    }
                });
            });
        },
        savePendingLocalNS() {
            this.testPendingLocalNS().then(() => {
                if (this.allTestPassed) {
                    axios.post('/api/settings/local_ns',
                        { server: this.shadowPendingLocalNS}
                    ).then(() => {
                        this.localNSUpdated = true;
                        this.loadCurrentLocalNS();
                    });
                }
            });
        },
        loadCurrentLocalNS() {
            axios.get('/api/settings/local_ns').then((r) => {
                this.currentLocalNS = r.data.data.server.replace(':53', '');
            })
        }
    },
    mounted() {
        this.loadCurrentLocalNS();
    },
    unmounted() {
    }
}
</script>