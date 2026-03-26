import axios from 'axios'
import { getRuntimeConfig } from '../utils/runtimeConfig'

export const apiClient = axios.create({
  baseURL: getRuntimeConfig().apiBaseUrl,
  withCredentials: true,
  timeout: 10000,
})

export function buildApiUrl(path: string): string {
  return `${getRuntimeConfig().apiBaseUrl}${path}`
}
