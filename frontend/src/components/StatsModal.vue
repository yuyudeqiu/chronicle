<script setup>
import { ref, onMounted } from 'vue'

const summary = ref(null)
const error = ref(null)

async function loadSummary() {
  try {
    const res = await fetch('/api/v1/stats/summary').then(r => r.json())
    if (res.code === 0) {
      summary.value = res.data
    } else {
      error.value = res.msg
    }
  } catch (err) {
    error.value = 'Failed to load summary stats'
  }
}

onMounted(() => {
  loadSummary()
})
</script>

<template>
  <div class="fixed inset-0 z-50 bg-dark-bg/95 flex items-center justify-center p-4 backdrop-blur-md overflow-y-auto" @click.self="$emit('close')">
    <div class="max-w-4xl w-full">
      <div class="flex items-center justify-between mb-8">
        <h1 class="text-3xl font-bold text-white tracking-tight flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-emerald-400 to-emerald-600 flex items-center justify-center text-white shadow-lg shadow-emerald-500/20">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path></svg>
          </div>
          General Statistics
        </h1>
        <button @click="$emit('close')" class="p-2 rounded-full hover:bg-white/10 text-slate-400 hover:text-white transition-colors">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </button>
      </div>

      <div v-if="error" class="bg-red-500/10 border border-red-500/20 text-red-400 p-4 rounded-xl text-sm mb-6">
        {{ error }}
      </div>

      <div v-else-if="summary" class="space-y-6">
        <!-- Top Stats Row -->
        <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
          <div class="bg-dark-card rounded-2xl p-6 border border-dark-border shadow-lg relative overflow-hidden group">
            <div class="absolute -right-4 -top-4 w-24 h-24 bg-blue-500/10 rounded-full blur-2xl group-hover:bg-blue-500/20 transition-colors"></div>
            <p class="text-sm font-medium text-slate-400 mb-2">Total Tasks</p>
            <p class="text-4xl font-bold text-white">{{ summary.total_tasks }}</p>
          </div>
          
          <div class="bg-dark-card rounded-2xl p-6 border border-dark-border shadow-lg relative overflow-hidden group">
            <div class="absolute -right-4 -top-4 w-24 h-24 bg-emerald-500/10 rounded-full blur-2xl group-hover:bg-emerald-500/20 transition-colors"></div>
            <p class="text-sm font-medium text-slate-400 mb-2">Completed</p>
            <p class="text-4xl font-bold text-emerald-400">{{ summary.completed_tasks }}</p>
          </div>

          <div class="bg-dark-card rounded-2xl p-6 border border-dark-border shadow-lg relative overflow-hidden group">
            <div class="absolute -right-4 -top-4 w-24 h-24 bg-indigo-500/10 rounded-full blur-2xl group-hover:bg-indigo-500/20 transition-colors"></div>
            <p class="text-sm font-medium text-slate-400 mb-2">In Progress</p>
            <p class="text-4xl font-bold text-indigo-400">{{ summary.in_progress_tasks }}</p>
          </div>

          <div class="bg-dark-card rounded-2xl p-6 border border-dark-border shadow-lg relative overflow-hidden group">
            <div class="absolute -right-4 -top-4 w-24 h-24 bg-purple-500/10 rounded-full blur-2xl group-hover:bg-purple-500/20 transition-colors"></div>
            <p class="text-sm font-medium text-slate-400 mb-2">Completion Rate</p>
            <p class="text-4xl font-bold text-purple-400">{{ Math.round(summary.completion_rate * 100) }}%</p>
          </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <!-- Categories -->
          <div class="bg-dark-card rounded-3xl p-6 border border-dark-border shadow-xl col-span-1">
            <h3 class="text-lg font-bold text-white mb-6">By Category</h3>
            <div class="space-y-4">
              <div v-for="(count, category) in summary.by_category" :key="category" class="flex items-center justify-between">
                <span class="text-sm font-medium text-slate-300">{{ category || 'Uncategorized' }}</span>
                <span class="text-sm text-slate-400 bg-dark-bg px-3 py-1 rounded-lg border border-dark-border">{{ count }}</span>
              </div>
              <div v-if="!summary.by_category || Object.keys(summary.by_category).length === 0" class="text-sm text-slate-500 italic py-4 text-center">
                No categories found.
              </div>
            </div>
          </div>

          <!-- Weekly Activity Array -->
          <div class="bg-dark-card rounded-3xl p-6 border border-dark-border shadow-xl col-span-2">
            <h3 class="text-lg font-bold text-white mb-6">Last 7 Days</h3>
            <div class="flex items-end justify-between h-48 gap-2 mt-4">
              <div v-for="day in summary.weekly_stats" :key="day.date" class="relative flex flex-col items-center justify-end h-full gap-2 group w-full">
                <!-- Data Tooltip -->
                <div class="absolute -top-12 bg-dark-bg border border-dark-border px-3 py-2 rounded-xl text-xs text-white opacity-0 group-hover:opacity-100 transition-opacity z-10 whitespace-nowrap shadow-xl pointer-events-none">
                  {{ day.date.slice(5) }}<br/>
                  <span class="text-emerald-400">Done: {{ day.completed }}</span> | 
                  <span class="text-indigo-400">New: {{ day.created }}</span>
                </div>

                <!-- Bars Container -->
                <div class="flex gap-1 items-end w-full max-w-[40px] h-full justify-center">
                  <!-- Created Bar -->
                  <div class="w-1/2 bg-indigo-500/50 group-hover:bg-indigo-400 rounded-t-sm transition-all" 
                       :style="{ height: `${Math.max(4, (day.created / Math.max(1, ...summary.weekly_stats.map(s => Math.max(s.created, s.completed)))) * 100)}%` }">
                  </div>
                  <!-- Completed Bar -->
                  <div class="w-1/2 bg-emerald-500/50 group-hover:bg-emerald-400 rounded-t-sm transition-all"
                       :style="{ height: `${Math.max(4, (day.completed / Math.max(1, ...summary.weekly_stats.map(s => Math.max(s.created, s.completed)))) * 100)}%` }">
                  </div>
                </div>

                <!-- Date Label -->
                <span class="text-[10px] text-slate-500 uppercase font-medium">{{ day.date.slice(5) }}</span>
              </div>
            </div>
            <!-- Legend -->
            <div class="flex gap-4 mt-6 justify-center">
              <div class="flex items-center gap-2"><div class="w-3 h-3 rounded-sm bg-indigo-500/50"></div><span class="text-xs text-slate-400">Created</span></div>
              <div class="flex items-center gap-2"><div class="w-3 h-3 rounded-sm bg-emerald-500/50"></div><span class="text-xs text-slate-400">Completed</span></div>
            </div>
          </div>
        </div>
      </div>
      
      <div v-else class="flex justify-center p-12">
        <div class="w-8 h-8 rounded-full border-2 border-indigo-500/20 border-t-indigo-500 animate-spin"></div>
      </div>
    </div>
  </div>
</template>
