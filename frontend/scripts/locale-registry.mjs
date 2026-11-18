import { createServer } from 'vite'

// Load the same explicit registry as the application, including its key checks.
export async function loadLocaleRegistry() {
  const server = await createServer({
    logLevel: 'error',
    server: { middlewareMode: true, hmr: false, ws: false },
    appType: 'custom',
  })
  try {
    const registry = await server.ssrLoadModule('/src/locales/index.ts')
    const merger = await server.ssrLoadModule('/src/locales/merge.ts')
    return { ...registry, ...merger }
  } finally {
    await server.close()
  }
}
