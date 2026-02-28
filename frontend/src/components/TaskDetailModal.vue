<script setup>
import { computed } from 'vue'

const props = defineProps({
  isOpen: Boolean,
  task: Object
})

const emit = defineEmits(['close', 'edit', 'delete', 'progress', 'delete-worklog', 'start'])

const formattedDeadline = computed(() => {
  if (!props.task || !props.task.deadline) return ''
  return new Date(props.task.deadline).toLocaleString()
})

const formattedCompletedAt = computed(() => {
  if (!props.task || !props.task.actual_completed_at) return ''
  return new Date(props.task.actual_completed_at).toLocaleString()
})

const parsedLinks = computed(() => {
  if (!props.task || !props.task.links) return []
  return props.task.links.split('\n').filter(l => l.trim())
})

const statusClass = computed(() => {
  if (!props.task) return ''
  if (props.task.status === 'done') {
    return 'text-[10px] uppercase tracking-wider font-semibold ml-2 text-emerald-400 bg-emerald-500/10 px-2 py-1 rounded border border-emerald-500/20'
  } else if (props.task.status === 'in-progress') {
    return 'text-[10px] uppercase tracking-wider font-semibold ml-2 text-indigo-400 bg-indigo-500/10 px-2 py-1 rounded border border-indigo-500/20'
  }
  return 'text-[10px] uppercase tracking-wider font-semibold ml-2 text-slate-400 bg-slate-500/10 px-2 py-1 rounded border border-slate-500/20'
})

function formatTime(isoString) {
  return new Date(isoString).toLocaleString()
}
</script>

<template>
  <div v-show="isOpen && task" class="fixed inset-0 z-50">
    <div class="fixed inset-0 bg-black/60 backdrop-blur-sm transition-opacity" @click="emit('close')"></div>
    <div class="fixed inset-0 z-10 overflow-y-auto">
      <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
        <div class="relative transform overflow-hidden rounded-2xl bg-dark-card border border-dark-border text-left shadow-2xl transition-all sm:my-8 sm:w-full sm:max-w-2xl">
          <div class="px-6 py-5 border-b border-dark-border flex justify-between items-center bg-gradient-to-r from-dark-card to-dark-bg">
            <h3 class="text-lg font-semibold text-white truncate pr-4">{{ task?.title }}</h3>
            <div class="flex gap-3">
              <button @click="emit('edit')" class="text-xs bg-slate-500/10 hover:bg-slate-500 text-slate-400 hover:text-white border border-slate-500/20 px-3 py-1.5 rounded-lg transition-colors">Edit</button>
              <button @click="emit('delete')" class="text-xs bg-red-500/10 hover:bg-red-500 text-red-500 hover:text-white border border-red-500/20 px-3 py-1.5 rounded-lg transition-colors">Delete</button>
              <button v-if="task?.status === 'todo'" @click="emit('start')" class="text-xs bg-emerald-500 hover:bg-emerald-600 text-white px-3 py-1.5 rounded-lg transition-colors shadow-lg shadow-emerald-500/20">Start Task</button>
              <button v-if="task?.status !== 'done'" @click="emit('progress')" class="text-xs bg-indigo-500 hover:bg-indigo-600 text-white px-3 py-1.5 rounded-lg transition-colors shadow-lg shadow-indigo-500/20">Update Progress</button>
              <button @click="emit('close')" class="flex-shrink-0 text-slate-400 hover:text-white transition-colors p-1">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
              </button>
            </div>
          </div>
          <div class="px-6 py-5 space-y-5 max-h-[70vh] overflow-y-auto">
            <div class="flex items-center">
              <span class="text-[10px] uppercase tracking-wider font-semibold text-indigo-400 bg-indigo-500/10 px-2 py-1 rounded border border-indigo-500/20">
                {{ task?.category }}
              </span>
              <span :class="statusClass">{{ task?.status }}</span>
            </div>
            
            <div class="flex flex-col gap-2 mt-1 mb-2">
              <div v-if="task?.deadline" class="flex items-center gap-2 text-xs text-slate-400">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
                <span>Deadline: <span class="text-slate-300 font-medium">{{ formattedDeadline }}</span></span>
              </div>
              <div v-if="task?.actual_completed_at" class="flex items-center gap-2 text-xs text-slate-400">
                <svg class="w-3.5 h-3.5 text-emerald-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                <span>Completed: <span class="text-emerald-400 font-medium">{{ formattedCompletedAt }}</span></span>
              </div>
            </div>

            <div v-if="task?.description">
              <h4 class="text-sm font-medium text-slate-400 mb-2">Description</h4>
              <div class="text-sm text-slate-300 bg-dark-bg p-4 rounded-xl border border-dark-border leading-relaxed">{{ task.description }}</div>
            </div>

            <div v-if="task?.targets">
              <h4 class="text-sm font-medium text-slate-400 mb-2">Targets</h4>
              <div class="text-sm text-slate-300 bg-dark-bg p-4 rounded-xl border border-dark-border leading-relaxed whitespace-pre-wrap">{{ task.targets }}</div>
            </div>

            <div v-if="parsedLinks.length > 0">
              <h4 class="text-sm font-medium text-slate-400 mb-2">Links</h4>
              <div class="text-sm text-slate-300 bg-dark-bg p-4 rounded-xl border border-dark-border leading-relaxed flex flex-col gap-1">
                <a v-for="(link, i) in parsedLinks" :key="i" :href="link" target="_blank" class="text-indigo-400 hover:text-indigo-300 underline break-all">{{ link }}</a>
              </div>
            </div>

            <div>
              <h4 class="text-sm font-medium text-slate-400 mb-3 flex items-center gap-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                Worklogs
              </h4>
              <div class="space-y-3">
                <div v-if="!task?.logs || task.logs.length === 0" class="text-sm text-slate-500 italic">No worklogs yet.</div>
                <div v-for="log in task?.logs" :key="log.id" class="bg-dark-bg p-3 rounded-lg border border-dark-border text-sm group relative">
                  <div class="flex justify-between items-start mb-1">
                    <div class="text-xs text-slate-500">{{ formatTime(log.created_at) }}</div>
                    <button @click="emit('delete-worklog', log.id)" title="Delete Worklog" class="text-xs text-red-500 hover:text-red-400 opacity-0 group-hover:opacity-100 transition-opacity p-1">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                    </button>
                  </div>
                  <div class="text-slate-200 whitespace-pre-wrap pr-6">{{ log.log_text }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
