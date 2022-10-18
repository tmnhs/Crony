//index.js
import Vue from 'vue'
import VueI18n from 'vue-i18n'
import elementEnLocale from 'element-ui/lib/locale/lang/en'
import elementZhLocale from 'element-ui/lib/locale/lang/zh-CN'
import elementLocal from 'element-ui/lib/locale'

import enLocale from './en'
import zhLocale from './zh'

Vue.use(VueI18n)

const messages = {
  en: {
    ...enLocale,
    ...elementEnLocale,
  },
  zh: {
    ...zhLocale,
    ...elementZhLocale,
  },
}

const i18n = new VueI18n({
  locale: localStorage.getItem('language') || 'zh',
  messages,
})

elementLocal.i18n((key, value) => i18n.t(key, value))


export default i18n
