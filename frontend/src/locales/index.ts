import commonEn from './common/en.json'
import commonRu from './common/ru.json'
import authEn from './auth/en.json'
import authRu from './auth/ru.json'
import profileEn from './profile/en.json'
import profileRu from './profile/ru.json'
import themesEn from './themes/en.json'
import themesRu from './themes/ru.json'
import demoEn from './demo/en.json'
import demoRu from './demo/ru.json'
import charactersEn from './characters/en.json'
import charactersRu from './characters/ru.json'
import vtmEn from './systems/vtm/en.json'
import vtmRu from './systems/vtm/ru.json'
import { mergeDictionaries } from './merge'

// Register each feature once; all modules share the application language.
export const localeModules = {
  common: { en: commonEn, ru: commonRu },
  auth: { en: authEn, ru: authRu },
  profile: { en: profileEn, ru: profileRu },
  themes: { en: themesEn, ru: themesRu },
  demo: { en: demoEn, ru: demoRu },
  characters: { en: charactersEn, ru: charactersRu },
  'systems/vtm': { en: vtmEn, ru: vtmRu },
}

export const en = mergeDictionaries(localeModules, 'en')
export const ru = mergeDictionaries(localeModules, 'ru')
