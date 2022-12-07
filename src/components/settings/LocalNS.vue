<template>
    <div class="grid grid-col-1 grid-flow-row space-y-6 w-[50%]">
        <span class="text-gray-700 block">Current Local NS</span>
        <template v-for="(ns, index) in currentLocalNS" :key="index">
            <div class="flex justify-start">
                <input type="text" v-model="currentLocalNS[index]"
                    class="mb-1 w-1/2 rounded-md border-gray-700 shadow-sm block">
                <button @click="removeNS(index)"
                    class="block px-5 text-base font-semibold rounded-md text-slate-200">X</button>
            </div>
        </template>
        <div class="flex ">
            <button @click="addNS"
                class="block bg-blue-500 hover:bg-blue-600 px-5 active:bg-blue-700 py-2 text-base font-semibold text-white rounded-md">Add</button>
            <button @click="testNS"
                class="block bg-blue-500 hover:bg-blue-600 px-5 active:bg-blue-700 py-2 ml-4 text-base font-semibold text-white rounded-md">Test</button>
            <button @click="saveNS"
                class="block bg-blue-500 hover:bg-blue-600 px-5 active:bg-blue-700 py-2 ml-4 text-base font-semibold text-white rounded-md">Save</button>
        </div>
    </div>
    <table class="border-0 border-collapse w-1/2 mt-6">
        <tr v-for="r in testResults" :key="(r.host + r.server)">
            <td>{{r.host}}</td>
            <td>{{r.server}}</td>
            <td>{{r.return}}</td>
        </tr>
    </table>
    <p v-if="allTestPassed" class="text-teal-600">All tests passed.</p>
    <p v-if="localNSUpdated" class="text-teal-600">Local NS IP changed.</p>
</template>

<script>
import axios from 'axios';
export default {
    data() {
        return {
            currentLocalNS: [],
            shadowCurrentLocalNS: [],
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
        async testNS() {
            this.localNSUpdated = false;
            this.allTestPassed = false;
            this.testResults = [];
            this.shadowCurrentLocalNS = this.currentLocalNS.map((ns) => {
                if (!ns.endsWith(':53')) {
                    return ns + ':53';
                }
                return ns
            });

            let tests = [];
            this.testDomains.map((h) => {
                this.shadowCurrentLocalNS.forEach((ns) => {
                    this.testResults.push({ host: h.host, server: ns, return: '?', success: false });
                    tests.push(this.runTest({ server: ns, host: h.host, shouldAnswer: h.shouldAnswer }));
                })
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
                { server: t.server, host: t.host }
            ).then((r) => {
                this.testResults.map((s) => {
                    if (s.host === t.host && s.server === t.server) {
                        s.return = this.ok(t.shouldAnswer, r.data.data.error, r.data.data.answer);
                        s.success = s.return === 'ok';
                    }
                });
            });
        },
        saveNS() {
            this.testNS().then(() => {
                if (this.allTestPassed) {
                    axios.post('/api/settings/local_ns',
                        { servers: this.shadowCurrentLocalNS }
                    ).then(() => {
                        this.localNSUpdated = true;
                        this.loadNS();
                    });
                }
            });
        },
        loadNS() {
            axios.get('/api/settings/local_ns').then((r) => {
                this.currentLocalNS = r.data.data.servers.map((ns) => ns.replace(':53', ''));
            })
        },
        addNS() {
            this.currentLocalNS.push('');
        },
        removeNS(i) {
            this.currentLocalNS.splice(i, 1);
        }
    },
    mounted() {
        this.loadNS();
    },
    unmounted() {
    }
}
</script>