<script setup>
import TaskCard from './TaskCard.vue'

defineProps({
  todos: { type: Array, required: true },
  inProgress: { type: Array, required: true },
  done: { type: Array, required: true }
})

const emit = defineEmits(['task-click', 'start'])
</script>

<template>
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 items-start">
    
    <!-- Todo Column -->
    <div class="flex flex-col gap-4">
      <div class="flex items-center gap-2 mb-2">
        <div class="w-2 h-2 rounded-full bg-slate-400"></div>
        <h2 class="font-semibold text-slate-300">To Do</h2>
        <span class="text-xs bg-dark-border px-2 py-0.5 rounded-full text-slate-400">
          {{ todos.length }}
        </span>
      </div>
      <div class="flex flex-col gap-3 min-h-[200px]">
        <TaskCard 
          v-for="task in todos" 
          :key="task.id" 
          :task="task" 
          @click="emit('task-click', $event)" 
          @start="emit('start', $event)"
        />
        <div v-if="todos.length === 0" class="text-sm text-slate-500 italic text-center py-4">No tasks to do</div>
      </div>
    </div>

    <!-- In Progress Column -->
    <div class="flex flex-col gap-4">
      <div class="flex items-center gap-2 mb-2">
        <div class="w-2 h-2 rounded-full bg-indigo-400 animate-pulse"></div>
        <h2 class="font-semibold text-slate-300">In Progress</h2>
        <span class="text-xs bg-dark-border px-2 py-0.5 rounded-full text-slate-400">
          {{ inProgress.length }}
        </span>
      </div>
      <div class="flex flex-col gap-3 min-h-[200px]">
        <TaskCard 
          v-for="task in inProgress" 
          :key="task.id" 
          :task="task" 
          @click="emit('task-click', $event)" 
        />
        <div v-if="inProgress.length === 0" class="text-sm text-slate-500 italic text-center py-4">No tasks in progress</div>
      </div>
    </div>

    <!-- Done Column -->
    <div class="flex flex-col gap-4">
      <div class="flex items-center gap-2 mb-2">
        <div class="w-2 h-2 rounded-full bg-emerald-400"></div>
        <h2 class="font-semibold text-slate-300">History / Done</h2>
        <span class="text-xs bg-dark-border px-2 py-0.5 rounded-full text-slate-400">
          {{ done.length }}
        </span>
      </div>
      <div class="flex flex-col gap-3 min-h-[200px] max-h-[70vh] overflow-y-auto pr-2">
        <TaskCard 
          v-for="task in done" 
          :key="task.id" 
          :task="task" 
          @click="emit('task-click', $event)" 
        />
        <div v-if="done.length === 0" class="text-sm text-slate-500 italic text-center py-4">No completed tasks yet</div>
      </div>
    </div>

    <!-- Quick Actions Column -->
    <div class="flex flex-col gap-4 relative">
      <div class="flex items-center gap-2 mb-2">
        <div class="w-2 h-2 rounded-full bg-emerald-400"></div>
        <h2 class="font-semibold text-slate-300">Quick Actions</h2>
      </div>

      <div class="glass-card rounded-xl p-5 border border-dark-border hover:border-indigo-500/30 transition-colors group cursor-pointer" onclick="window.open('/api/v1/exports/daily-markdown', '_blank')">
        <div class="flex items-center gap-3 mb-2">
          <div class="p-2 bg-indigo-500/10 rounded-lg text-indigo-400 group-hover:bg-indigo-500 group-hover:text-white transition-colors">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
          </div>
          <h3 class="font-medium text-slate-200">Export Daily Markdown</h3>
        </div>
        <p class="text-xs text-slate-400">Generate Obsidian compatible markdown format of completed tasks for today.</p>
      </div>

      <div class="glass-card rounded-xl p-5 border border-dark-border hover:border-emerald-500/30 transition-colors group cursor-pointer" onclick="window.open('/summary.html', '_blank')">
        <div class="flex items-center gap-3 mb-2">
          <div class="p-2 bg-emerald-500/10 rounded-lg text-emerald-400 group-hover:bg-emerald-500 group-hover:text-white transition-colors">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
          </div>
          <h3 class="font-medium text-slate-200">View Daily Summary</h3>
        </div>
        <p class="text-xs text-slate-400">View a visual summary page of all tasks and progress you made today.</p>
      </div>

      <div class="glass-card rounded-xl p-5 border border-dark-border hover:border-emerald-500/30 transition-colors group cursor-pointer" onclick="window.open('/api/v1/reports/daily-summary', '_blank')">
        <div class="flex items-center gap-3 mb-2">
            <div class="p-2 bg-emerald-500/10 rounded-lg text-emerald-400 group-hover:bg-emerald-500 group-hover:text-white transition-colors">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                </svg>
            </div>
            <h3 class="font-medium text-slate-200">Daily JSON Summary</h3>
        </div>
        <p class="text-xs text-slate-400">View the raw JSON report built for AI Agents to process.</p>
      </div>

    </div>
  </div>
</template>
