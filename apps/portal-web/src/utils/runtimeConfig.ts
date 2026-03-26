export interface RuntimeConfig {
  apiBaseUrl: string
  appTitle: string
}

declare global {
  interface Window {
    __PORTAL_CONFIG__?: RuntimeConfig
  }
}

export function getRuntimeConfig(): RuntimeConfig {
  return {
    apiBaseUrl: window.__PORTAL_CONFIG__?.apiBaseUrl ?? '/api',
    appTitle: window.__PORTAL_CONFIG__?.appTitle ?? 'Portal',
  }
}
