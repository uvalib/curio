import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

export default new Vuex.Store({
   state: {
      working: false,
      viewType: "none",
      iiifURL: "",
      rightsURL: "",
      pagePIDs: [],
      startPage: 0,
      PID: "",
      wslsData: {},
      falied: false
   },
   mutations: {
      setWorking(state, flag) {
         state.working = flag
         if (flag == true ) {
            state.failed = false
         }
      },
      setFailed(state) {
         state.failed = true
      },
      setPID(state, pid) {
         state.PID = pid
      },
      setViewData(state, resp) {
         let data = resp.data
         state.viewType = resp.type
         if ( resp.type == 'iiif') {
            state.iiifURL  = data.iiif
            state.rightsURL = data.rights
            state.pagePIDs.slice(0, state.pagePIDs.length-1)
            data.page_pids.split(",").forEach( p=>{
               state.pagePIDs.push(p)
            })
            state.startPage = data.page
         } else if (resp.type == 'wsls') {
            state.wslsData = data
         }
      }
   },
   actions: {
      async getPIDViewData(ctx, {pid, page}) {
         ctx.commit("setWorking", true)
         ctx.commit("setPID", pid)
         await axios.get(`/api/view/${pid}?page=${page}`).then(response => {
            ctx.commit('setViewData', response.data)
            ctx.commit("setWorking", false)
         }).catch(error => {
            ctx.commit('setFailed')
            ctx.commit("setWorking", false)
         })
      }
   },
   modules: {
   }
})
