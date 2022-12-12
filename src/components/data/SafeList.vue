<template>
    <div class="w-[35%]">
        Safelist Filter:
        <input type="text" placeholder="Filter here. eg. github.com" 
            @input="loadSafeList" v-model="filter" class="rounded-md mb-6 w-full border-slate-300">
    </div>
    <table class="border-0 min-w-[75%] w-auto">
        <tr class="border-b border-slate-600">
            <th class="w-1/2 text-left cursor-pointer hover:text-blue-600" @click="sortList('name')">Host</th>
            <th class="w-1/4 text-left">QType</th>
            <th class="w-1/4 text-left cursor-pointer hover:text-blue-600" @click="sortList('ttl')">TTL</th>
        </tr>
        <tr v-for="item in list" :key="(item.name+item.qtype)" class="border-b border-slate-300">
            <td>{{item.name}}</td>
            <td>{{item.qtype}}</td>
            <td>{{item.ttl == -999999 ? '∞' : item.ttl}}</td>
        </tr>
    </table>
</template>

<script>
import axios from 'axios';
import _ from 'lodash';

export default {
    data() {
        return {
            list: [],
            filter: '',
            ascendCol: 'name',
        }
    },
    methods: {
        loadSafeList: _.debounce(function() {
            axios.get('/api/data/safelist')
                .then((r) => {
                    this.list = r.data.data.list.sort((a, b) => {
                        return a.name > b.name ? 1 : -1;
                    }).filter((a) => {
                        if (this.filter != '') {
                            return a.name.includes(this.filter);
                        } else {
                            return true;
                        }
                    })
                });
        }, 200),
        sortList(col) {
            var ascend = (col === this.ascendCol);
            if (ascend) {
                this.ascendCol = '';
            } else {
                this.ascendCol = col;
            }
            if (col === 'name') {
                this.list.sort((a, b) => {
                    if (ascend) {
                        return a.name > b.name ? 1 : -1;
                    } else {
                        return a.name < b.name ? 1 : -1;
                    }
                })
            } else if (col === 'ttl') {
                this.list.sort((a, b) => {
                    if (ascend) {
                        return a.ttl > b.ttl ? 1 : -1;
                    } else {
                        return a.ttl < b.ttl ? 1 : -1;
                    }
                })
            }
        }
    },
    mounted() {
        this.loadSafeList();
    }
}
</script>