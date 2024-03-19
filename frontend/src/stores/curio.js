import { defineStore } from 'pinia'
import axios from 'axios'

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
         await axios.get(url).then( async response => {
            if (response.data.type == "iiif") {
               let manUrl = response.data.data.iiif
               let manResp = await axios.get(manUrl)
               let respMeta = manResp.data.metadata
               if ( respMeta ) {
                  respMeta.forEach( m => {
                     if (m.label == 'Content Advisory') {
                        this.advisory = m.value
                     }
                  })
               }
               this.setViewData(response.data)
               this.working = false
            } else {
               this.setViewData(response.data)
               this.working = false
            }
         }).catch( () => {
            this.failed = true
            this.working = false
         })
      }
   },
})
