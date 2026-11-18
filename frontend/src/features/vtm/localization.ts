import { t } from '../../lib/i18n'
import { labelKeys } from './labelKeys'
export type MessageValue = string | number | Message | Message[]
export type Message = { key: string; params?: Record<string, MessageValue> }
export const msg = (key: string, params?: Message['params']): Message => ({ key, params })
export function formatMessage(value: MessageValue | undefined | null): string {
  if (value == null) return ''
  if (typeof value === 'string' || typeof value === 'number') return String(value)
  if (Array.isArray(value)) return value.map(formatMessage).join(' ')
  return t(
    value.key,
    Object.fromEntries(
      Object.entries(value.params ?? {}).map(([key, value]) => [key, formatMessage(value)]),
    ),
  )
}
export const sameMessage = (a: Message, b: Message) => JSON.stringify(a) === JSON.stringify(b)
export const joinMessages = (values: Message[]) => msg('vtm.messages', { messages: values })
export const l = (label: string) => (labelKeys[label] ? t(labelKeys[label]) : label)
export const labelMessage = (label: MessageValue | undefined): MessageValue =>
  typeof label === 'string' && labelKeys[label] ? msg(labelKeys[label]) : (label ?? '')
// Preserve separation between translated text and adjacent numeric JSX expressions.
export const jsxText = (key: string) => ` ${t(key)} `
export const message = (source: string, params?: Message['params']) =>
  msg(labelKeys[source], params)
export const copy = (source: string, params?: Message['params']) =>
  formatMessage(message(source, params))
