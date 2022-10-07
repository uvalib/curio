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
      failed: false
   }),
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

      async getPIDViewData( pid, page, unit ) {
         this.working =  true
         this.failed = false
         this.pid = pid
         let url = `/api/view/${pid}?page=${page}`
         if (unit ) {
            url += `&unit=${unit}`
         }
         await axios.get(url).then(response => {
            this.setViewData(response.data)
            this.working = false
         }).catch( () => {
            this.failed = true
            this.working = false
         })
      }
   },
   modules: {
   }
})
