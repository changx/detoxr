<template>
    <div class="fixed inset-0 top-20 w-60 right-auto left-auto overflow-y-auto pb-4">
        <nav id="nav" class="relative">
            <ul>
                <li class="mb-12" v-for="(menu, index) in menuItems" :key="index">
                    <h5 class="mb-8 font-semibold text-slate-900">{{ menu.title }}</h5>
                    <ul class="border-l border-slate-100 space-y-2 ">
                        <li v-for="(item, index) in menu.items" :key="index">
                            <router-link active-class="font-bold border-blue-500 text-blue-500"
                                class="pl-4 py-2 -ml-px block border-l border-transparent hover:border-blue-500 border-slate-100 text-slate-700 hover:text-blue-500"
                                :to="{ path: item.route }">{{ item.title }}</router-link>
                        </li>
                    </ul>
                </li>
            </ul>
        </nav>
    </div>
    <div class="pl-80 pt-10 pb-10">
        <router-view></router-view>
    </div>
</template>

<script>
export default {
  data() {
    return {
      menuItems: [
        {
          title: 'Data View',
          items: [
            { title: 'Statistics', route: '/data/stats' },
            { title: 'Safe List', route: '/data/safelist' },
            { title: 'Victims', route: '/data/victims' }
          ]
        },
        {
          title: 'Settings',
          items: [
            { title: 'Honeypot', route: '/settings/honeypot' },
            { title: 'Nameserver', route: '/settings/local_ns' },
            { title: 'DNS-over-HTTPS', route: '/settings/doh_service' }
          ]
        },
      ]
    }
  }
}
</script>