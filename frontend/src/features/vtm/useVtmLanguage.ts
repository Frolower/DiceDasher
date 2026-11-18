import { useSyncExternalStore } from 'react'
import { getLanguage, subscribeLanguage } from '../../lib/i18n'
export const useVtmLanguage = () =>
  useSyncExternalStore(subscribeLanguage, getLanguage, getLanguage)
