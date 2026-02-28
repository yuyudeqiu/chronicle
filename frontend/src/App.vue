<script setup>
import { ref, computed, onMounted } from 'vue'

// Import our components
import TaskBoard from './components/TaskBoard.vue'
import TaskFormModal from './components/TaskFormModal.vue'
import TaskDetailModal from './components/TaskDetailModal.vue'
import ProgressModal from './components/ProgressModal.vue'

// -- GLOBAL STATE --
const tasks = ref([])
const activeTask = ref(null) // Used for detail, edit, and progress

// -- MODAL VISIBILITY STATE --
const isFormModalOpen = ref(false)
const isDetailModalOpen = ref(false)
const isProgressModalOpen = ref(false)

const isEditMode = ref(false) // Whether FormModal is Create or Edit

// -- COMPUTED COLUMNS --
const todos = computed(() => tasks.value.filter(t => t.status === 'todo'))
const inProgress = computed(() => tasks.value.filter(t => t.status === 'in-progress'))
const done = computed(() => tasks.value.filter(t => t.status === 'done'))

// -- TOAST STATE --
const toastMsg = ref('')
const showToastMsg = ref(false)
let toastTimeout = null

function showToast(msg) {
  toastMsg.value = msg
  showToastMsg.value = true
  clearTimeout(toastTimeout)
  toastTimeout = setTimeout(() => {
    showToastMsg.value = false
  }, 3000)
}

// -- DATA FETCHING --
async function loadTasks() {
  try {
    const [activeRes, doneRes] = await Promise.all([
      fetch('/api/v1/tasks?status=todo,in-progress').then(r => r.json()),
      fetch('/api/v1/tasks?status=done').then(r => r.json())
    ])
    
    let allTasks = []
    if (activeRes.code === 0 && activeRes.data) {
      allTasks = [...allTasks, ...activeRes.data]
    }
    if (doneRes.code === 0 && doneRes.data) {
      allTasks = [...allTasks, ...doneRes.data]
    }
    tasks.value = allTasks
  } catch (err) {
    showToast('Failed to load tasks')
    console.error(err)
  }
}

async function loadTaskDetail(id) {
  try {
    const res = await fetch(`/api/v1/tasks/${id}`).then(r => r.json())
    if (res.code === 0) {
      activeTask.value = res.data
    } else {
      showToast('Failed to load task details')
    }
  } catch (err) {
    showToast('Network error loading details')
  }
}

onMounted(() => {
  loadTasks()
})

// -- HANDLERS --
function handleNewTaskClick() {
  activeTask.value = null
  isEditMode.value = false
  isFormModalOpen.value = true
}

async function handleTaskClick(task) {
  await loadTaskDetail(task.id)
  isDetailModalOpen.value = true
}

async function handleFormSubmit(payload) {
  try {
    let res
    if (isEditMode.value) {
      // Edit
      res = await fetch(`/api/v1/tasks/${payload.id}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      }).then(r => r.json())
    } else {
      // Create
      res = await fetch('/api/v1/tasks', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      }).then(r => r.json())
    }

    if (res.code === 0) {
      isFormModalOpen.value = false
      showToast(isEditMode.value ? 'Task updated successfully' : 'Task created successfully')
      await loadTasks()
      
      // Refresh detail modal if it was an edit
      if (isEditMode.value && activeTask.value) {
         await loadTaskDetail(activeTask.value.id)
         isDetailModalOpen.value = true
      }
    } else {
      showToast('Failed to save task: ' + res.msg)
    }
  } catch (err) {
    showToast('Network error saving task')
  }
}

function handleOpenEdit() {
  isDetailModalOpen.value = false
  isEditMode.value = true
  isFormModalOpen.value = true
}

async function handleDeleteTask() {
  if (!activeTask.value) return
  if (!confirm("Are you sure you want to delete this task? All of its worklogs will also be permanently deleted.")) return

  try {
    const res = await fetch(`/api/v1/tasks/${activeTask.value.id}`, { method: 'DELETE' }).then(r => r.json())
    if (res.code === 0) {
      isDetailModalOpen.value = false
      showToast('Task deleted successfully')
      loadTasks()
    } else {
      showToast('Failed to delete task: ' + res.msg)
    }
  } catch (err) {
    showToast('Network error deleting task')
  }
}

function handleOpenProgress() {
  isDetailModalOpen.value = false
  isProgressModalOpen.value = true
}

async function handleProgressSubmit(payload) {
  try {
    const res = await fetch(`/api/v1/tasks/${payload.taskId}/progress`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        log_text: payload.log_text,
        mark_as_done: payload.mark_as_done,
        new_status: payload.new_status,
        deadline: payload.deadline
      })
    }).then(r => r.json())

    if (res.code === 0) {
      isProgressModalOpen.value = false
      showToast(payload.mark_as_done ? 'Task marked as done! 🎉' : 'Progress logged')
      loadTasks()
    } else {
       showToast('Failed to log progress: ' + res.msg)
    }
  } catch(err) {
     showToast('Network error logging progress')
  }
}

async function handleDeleteWorklog(logId) {
  if (!confirm("Are you sure you want to delete this worklog?")) return
  try {
    const res = await fetch(`/api/v1/worklogs/${logId}`, { method: 'DELETE' }).then(r => r.json())
    if (res.code === 0) {
       showToast('Worklog deleted successfully')
       if (activeTask.value) {
         // Refresh detail view silently
         await loadTaskDetail(activeTask.value.id)
       }
    } else {
       showToast('Failed to delete worklog: ' + res.msg)
    }
  } catch (err) {
    showToast('Network error deleting worklog')
  }
}
</script>

<template>
  <nav class="sticky top-0 z-40 w-full backdrop-blur flex-none transition-colors duration-500 lg:z-50 border-b border-dark-border bg-dark-bg/80">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-lg bg-gradient-to-tr from-indigo-500 to-purple-500 flex items-center justify-center text-white font-bold text-xl shadow-lg shadow-indigo-500/20">
            C
          </div>
          <span class="font-bold text-xl tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-slate-200 to-slate-400">Chronicle</span>
        </div>
        <button @click="handleNewTaskClick" class="flex items-center gap-2 bg-indigo-500 hover:bg-indigo-600 text-white px-4 py-2 rounded-lg text-sm font-medium transition-all shadow-lg shadow-indigo-500/20 hover:shadow-indigo-500/40 active:scale-95">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
          </svg>
          New Task
        </button>
      </div>
    </div>
  </nav>

  <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <TaskBoard 
      :todos="todos" 
      :inProgress="inProgress" 
      :done="done" 
      @task-click="handleTaskClick" 
    />
  </main>

  <!-- Modals -->
  <TaskFormModal 
    :is-open="isFormModalOpen"
    :task="isEditMode ? activeTask : null"
    @close="isFormModalOpen = false"
    @submit="handleFormSubmit"
  />

  <TaskDetailModal 
    :is-open="isDetailModalOpen"
    :task="activeTask"
    @close="isDetailModalOpen = false"
    @edit="handleOpenEdit"
    @delete="handleDeleteTask"
    @progress="handleOpenProgress"
    @delete-worklog="handleDeleteWorklog"
  />

  <ProgressModal 
    :is-open="isProgressModalOpen"
    :task="activeTask"
    @close="isProgressModalOpen = false"
    @submit="handleProgressSubmit"
  />

  <!-- Toast Notification -->
  <div class="fixed bottom-4 right-4 z-50 transform transition-all duration-300" :class="[showToastMsg ? 'translate-y-0 opacity-100' : 'translate-y-20 opacity-0']">
    <div class="bg-emerald-500 text-white px-4 py-3 rounded-lg shadow-lg flex items-center gap-3">
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
      </svg>
      <span class="text-sm font-medium">{{ toastMsg }}</span>
    </div>
  </div>
</template>
