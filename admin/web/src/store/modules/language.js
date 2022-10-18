export default{
  state: {
    language: localStorage.getItem('language') || 'zh',
  },

  mutations: {
    SET_LANGUAGE: (state, language) => {
      state.language = language
      localStorage.setItem('language', language)
    },
  },

  actions: {
    setLanguage({ commit }, language) {
      commit('SET_LANGUAGE', language)
    },
  },

}
