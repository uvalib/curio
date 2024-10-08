import { defineStore } from 'pinia'
import { useFetch } from '@vueuse/core'

export const useCurioStore = defineStore('curio', {
	state: () => ({
      working: false,
      viewType: "none",
      iiifURL: "",
      rightsURL: "",
      pagePIDs: [],
      startPage: 0,
      PID: "",
      wslsData: {},
      archivematicaData: {},
      failed: false,
      advisoryCleared: false,
      advisory: "",
   }),
   getters: {
      hasAdvisory: state => {
         return state.advisory != "" && state.advisoryCleared == false
      }
   },
   actions: {
      setViewData(resp) {
         let data = resp.data
         this.viewType = resp.type
         if ( resp.type == 'iiif') {
            this.iiifURL  = data.iiif
            this.rightsURL = data.rights
            this.pagePIDs.slice(0, this.pagePIDs.length-1)
            data.page_pids.split(",").forEach( p=>{
               this.pagePIDs.push(p)
            })
            this.startPage = data.page
         } else if (resp.type == 'wsls') {
            this.wslsData = data
         } else if (resp.type == 'archivematica') {
            this.archivematicaData = [data].flat()
         }
      },

      clearAdvisory() {
         this.advisoryCleared = true
      },

      async getPIDViewData( pid, page, unit ) {
         this.working =  true
         this.failed = false
         this.advisory = ""
         this.advisoryCleared = false
         this.pid = pid
         let url = `/api/view/${pid}?page=${page}`
         if (unit ) {
            url += `&unit=${unit}`
         }
         const { error, data } = await useFetch(url)
         if ( error.value ) {
            this.failed = true
            this.working = false
         } else {
            const resp = JSON.parse(data.value)
            if (resp.type == "iiif") {
               const { data } = await useFetch( resp.data.iiif)
               const iiifResp = JSON.parse(data.value)
               if ( iiifResp && iiifResp.metadata ) {
                  iiifResp.metadata.forEach( m => {
                     if (m.label == 'Content Advisory') {
                        this.advisory = m.value
                     }
                  })
               }
            }
            this.setViewData(resp)
            this.working = false
         }
      }
   },
})
