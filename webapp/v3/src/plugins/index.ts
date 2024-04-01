// Plugins
import router from '../router'

import pinia from '../stores'



// Types
import type { App } from 'vue'

export function registerPlugins (app: App) {
  app
    .use(router)
    .use(pinia)
}
