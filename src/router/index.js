import { createRouter, createWebHistory } from 'vue-router';

const routes = [

    {
        path: '/',
        name: 'DashBoard',
        component: () => import('@/components/DashBoard.vue'),
        children: [
            {
                path: '/settings/honeypot',
                name: 'HoneyPot',
                component: () => import('@/components/settings/HoneyPot.vue')
            },
            {
                path: '/settings/local_ns',
                name: 'LocalNS',
                component: () => import('@/components/settings/LocalNS.vue')
            },
            {
                path: '/settings/doh_service',
                name: 'DohService',
                component: () => import('@/components/settings/DohService.vue')
            },
            {
                path: '/data/stats',
                name: 'Statistics',
                component: () => import('@/components/data/ServerStatistics.vue'),
            },
            {
                path: '/data/safelist',
                name: 'Safelist',
                component: () => import('@/components/data/SafeList.vue'),
            },
            {
                path: '/data/victims',
                name: 'VictimList',
                component: () => import('@/components/data/VictimList.vue')
            },
        ]
    },
    {
        path: '/about',
        name: 'AboutPage',
        component: () => import('@/views/AboutProject.vue'),
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router;