<template>
   <Toast position="top-center" />
   <div class="viewer">
      <WaitSpinner v-if="curio.working" :overlay="true" message="Loading viewer..." />
      <template v-else>
         <template v-if="curio.viewType=='iiif'">
            <template v-if="curio.hasAdvisory">
               <div class="advisory-dimmer"></div>
               <div class="advisory">
                  <span class="icon"></span>
                  <h5>Content advisory</h5>
                  <p>{{  curio.advisory }}</p>
                  <button @click="curio.clearAdvisory()">Show content</button>
               </div>
            </template>
            <div class="extra-tools hidden">
               <Button icon="pi pi-download" rounded text aria-label="download image" @click="downloadImage"/>
               <Button  v-if="canClipboard" rounded text aria-label="copy iiif manifest" @click="iiifManifestClicked">
                  <template #icon>
                     <img src="/iiif.svg" style="width: 24px;"/>
                  </template>
               </Button>
            </div>
            <div v-if="curio.hasAdvisory" id="tify-viewer" style="height:100%;" inert></div>
            <div v-else id="tify-viewer" style="height:100%;"></div>
         </template>
         <div v-else-if="curio.viewType=='wsls'" class="wsls">
            <div class="overview">
               <h3>{{curio.wslsData.title}}</h3>
               <p>{{curio.wslsData.description}}</p>
            </div>
            <div v-if="curio.wslsData.has_video" class="video-container" >
               <video class="video-js vjs-default-skin vjs-big-play-centered vjs-fluid" controls preload="auto"
                  :poster="curio.wslsData.poster_url" data-setup='{"inactivityTimeout": 0}'
                  crossorigin="anonymous"
               >
                  <source :src="curio.wslsData.video_url" type='video/mp4'>
                  <track kind="subtitles"
                     :src="curio.wslsData.video_url.replace(/\.[^/.]+$/, '.vtt')"
                     label="English" srclang="en"
                  />
                  <p class="vjs-no-js">
                     To view this video please enable JavaScript, and consider upgrading to a web browser that
                     <a href="http://videojs.com/html5-video-support/" target="_blank">supports HTML5 video</a>
                  </p>
               </video>
               <p class="duration">Duration: {{curio.wslsData.duration}}</p>
            </div>
            <div v-if="curio.wslsData.has_script" class="anchorscript-container">
               <h4>Anchor Script</h4>
               <img :src="curio.wslsData.thumb_url"/>
               <div class="anchorscript-links">
                  <a :href="curio.wslsData.pdf_url" target="_blank">View anchor script PDF in new tab</a>
                  <a :href="curio.wslsData.transcript_url" target="_blank">View anchor script transcription in new tab</a>
               </div>
             </div>
         </div>
         <div v-else-if="curio.viewType==='archivematica'">
               <TreeViewer :treeData="curio.archivematicaData"/>
         </div>
         <div v-else class="not-found">
            <h2>Sorry, but the resource you requested could not be found.</h2>
         </div>
      </template>
   </div>
</template>

<script setup>
import { useCurioStore } from "@/stores/curio"
import { onMounted, ref, onBeforeUnmount } from "vue"
import { useRoute, useRouter } from "vue-router"
import  TreeViewer  from "@/components/TreeViewer.vue"
import { useClipboard, usePermission, useTimeoutPoll } from '@vueuse/core'
import { useToast } from "primevue/usetoast"
import Toast from 'primevue/toast'

import 'tify'
import 'tify/dist/tify.css'

const canClipboard = usePermission('clipboard-write')
const { copy, copied } = useClipboard()

const toast = useToast()
const curio = useCurioStore()
const route = useRoute()
const router = useRouter()

const tgtDomain = ref("")
const viewer = ref(null)

onMounted( async () => {
   let pid = route.params.pid
   let page = route.query.page
   let unitID = route.query.unit
   if (!page) page = "1"

   await curio.getPIDViewData(pid, page, unitID)

   // the domain param is the transport and host of the parent window.
   // it is used to post messages from the viewer iFrame to the parent so the URL can be
   // updated to reflect image settings and viewer size.
   tgtDomain.value = route.query.domain

   if ( curio.viewType == 'iiif' ) {
      let pages = null
      let zoom = null
      let rotation = null
      let pan = {}
      if (route.query.page) {
         pages = [parseInt(route.query.page,10)]
      }
      if (route.query.zoom) {
         zoom = parseFloat(route.query.zoom)
      }
      if (route.query.rotation) {
         rotation = parseInt(route.query.rotation, 10)
      }
      if (route.query.x) {
         pan.x = parseFloat(route.query.x)
      }
      if (route.query.y) {
         pan.y = parseFloat(route.query.y)
      }

      viewer.value = new Tify({
         manifestUrl: curio.iiifURL,
         optionsResetOnPageChange: [],
         pages: pages,
         zoom: zoom,
         pan: pan,
         rotation: rotation,
         viewer: {
            immediateRender: false,
         },
      })
      viewer.value.mount('#tify-viewer')
      useTimeoutPoll( changeParam, 1000, { immediate: true } )
      setTimeout( ()=> {
         let extra = document.getElementsByClassName("extra-tools")[0]
         let orig = document.getElementsByClassName("tify-scan-buttons")[0]
         orig.appendChild(extra)
         extra.classList.remove("hidden")
      }, 1000)
   }

   if ( tgtDomain.value) {
      setTimeout(dimensionsMessage, 500)
   }
})

onBeforeUnmount(()=>{
   if ( viewer.value ) {
      viewer.value.destroy()
   }
})

const dimensionsMessage = (() =>{
   const message = {
      dimensions: {
         height: document.documentElement.scrollHeight + 'px',
         width: document.body.scrollWidth + 'px',
      }
   }
   window.top.postMessage(message, tgtDomain.value)
})

const iiifManifestClicked = (() => {
   copy(curio.iiifURL)
   if (copied) {
      toast.add({severity:'success', summary:  "Copied", detail:  "IIIF URL copied to clipboard.", life: 5000})
   }
})

const changeParam = (() => {
   let opts = viewer.value.options
   let origQ = route.query
   let q = Object.assign({}, origQ)
   delete q.x
   delete q.y
   delete q.zoom
   delete q.rotation
   delete q.page
   if (opts.zoom ) {
      q.zoom = opts.zoom
   }
   if (opts.rotation ) {
      q.rotation = opts.rotation
   }
   if (opts.pan ) {
      q.x = opts.pan.x
      q.y = opts.pan.y
   }
   if ( opts.pages ) {
      q.page = opts.pages[0]
   }

   if (q.zoom != origQ.zoom || q.rotation != origQ.rotation || q.x != origQ.x || q.y != origQ.y || q.page != origQ.page ) {
      router.replace({query: q})

      if ( tgtDomain.value ) {
         let evt = {name: "curio"}
         if ( q.x ) evt.x = q.x
         if ( q.y ) evt.y = q.y
         if ( q.zoom ) evt.zoom = q.zoom
         if ( q.rotation ) evt.rotation = q.rotation
         if ( q.page ) evt.page = q.page
         window.top.postMessage(evt, tgtDomain.value)
      }
   }
})

const downloadImage = (() => {
   let page = 0
   let url = new URL(window.location.href)
   let pageStr = url.searchParams.get("page")
   if (pageStr && pageStr.length > 0) {
      page = parseInt(pageStr, 10)-1
   }
   if (page < 0) page = 0
   let tgtPID =  curio.pagePIDs[page]
   let dlURL = `${curio.rightsURL}/${tgtPID}`
   var link = document.createElement('a')
   link.href = dlURL+"?download=1"
   document.body.appendChild(link)
   link.click()
   document.body.removeChild(link)
})
</script>

<style lang="scss">
.tify-header {
   background: $uva-grey-200;
   // toolbar button focs styles
   button.tify-header-button, button.tify-page-select-button {
     &:focus {
         color: $uva-text-color-dark;
         outline: 2px solid  $uva-brand-blue-100;
         outline-offset: -1px;
     }
   }
}

// styles for buttons inside the view
.tify-scan {
   background-color: $uva-grey-A;
   // pagination
   button.tify-scan-page-button {
      &:focus {
         background-color: white;
         outline: 2px solid $uva-blue-alt-300;
         outline-offset: 0.2rem;
     }
   }

   // pan, zoom, etc
   button.tify-scan-button {
      border-radius: 50px;
      &:focus {
         background-color: black;
         outline: 2px solid  $uva-blue-alt-300;
         outline-offset: 0px;
     }
   }
}

@media only screen and (min-width: 768px) {
   .advisory {
      max-width: 390px;
      max-height: 340px;
      border: 1px solid #F3EC45;
      padding: 25px;
   }
}
@media only screen and (max-width: 768px) {
   .advisory {
      width: 90%;
   }
}
h3 {
   text-align: left;
}
.tify-info-section.-title {
   text-align: left;
}
div.tify-info-metadata {
   text-align: left;
   h4 {
      font-weight: bold;
      margin-bottom: 5px;
      font-size: 0.95em;
   }
   .tify-info-content {
      margin-left: 15px;
      font-size: 0.95em;
   }
}
div.tify-info-section.-logo {
   border-top: 1px solid #dedede;
   padding-top:15px;
   img {
      margin: 0 auto;
   }
}
.viewer {
   height: 100%;
   position: relative;
   .advisory-dimmer {
      position: fixed;
      top: 55px;
      left: 0;
      width: 100%;
      height: 100%;
      background-color: rgba(10,10,10,0.9);
      z-index: 9999;
      -webkit-backdrop-filter: blur(10px);
      backdrop-filter: blur(10px);
   }
   .advisory {
      position: absolute;
      top: 50%; left: 50%;
      transform: translate(-50%,-50%);
      opacity: 1;
      background: #2b2b2b;
      border-radius: 10px;
      z-index: 100000;
      .icon {
         display: block;
         width: 60px;
         height: 60px;
         background-image: url(/src/assets/eye-slash.svg);
         background-repeat: no-repeat;
         background-position: center center;
         margin: 10px auto 20px auto;
      }
      h5 {
         font-family: "franklin-gothic-urw-medium", arial, sans-serif;
         -webkit-font-smoothing: antialiased;
         -moz-osx-font-smoothing: grayscale;
         font-size: 20px;
         margin: 17px 0;
         padding: 0;
         color: white;
      }
      p {
         font-family: "franklin-gothic-urw-medium", arial, sans-serif;
         padding: 0;
         font-size: 17px;
         padding: 0;
         margin: 0 0 17px 0;
         color: white;
      }
      button {
         margin: 5px 0 20px 0;
         border-radius:  5px;
         background-color: #BFE7F7;
         border: 2px solid #007bac;
         padding: 0.5rem 1rem;
         font-size: 17px;
         font-family: "franklin-gothic-urw-medium", arial, sans-serif;
         cursor: pointer ;
         &:hover {
            background-color: #91d8f2;
         }
      }
   }
}

.extra-tools.hidden {
   display: none;
}

.extra-tools {
   display: flex;
   flex-direction: column;
   align-items: flex-start;
   justify-content: center;
   gap: 8px;
   i {
      padding: 8px;
      font-size: 1.25em;
      color: white;
      cursor: pointer;
      border-radius: 50px;
      &:hover {
         -webkit-backdrop-filter: blur(2px);
         backdrop-filter: blur(2px);
         background: rgba(0, 0, 0, .2);
      }
   }
   .iiif {
      padding: 3px;
   }
   button.p-button-icon-only {
      .p-button-icon {
         font-size: 1.2rem;
      }
      &:focus, &:focus-visible {
         background: black;
      }
   }
   .iiif, .download {
      cursor: pointer;
      border-radius: 50px;
         &:focus {
            background-color: black;
            outline: 2px solid $uva-blue-alt-300;
            outline-offset: 0px;
      }
      img {
         padding: 6px 4px 0 4px;
         width: 32px;
         &:hover {
            -webkit-backdrop-filter: blur(2px);
            backdrop-filter: blur(2px);
            background: rgba(0, 0, 0, .2);
         }
      }
   }
}
.-controls {
   .-view:nth-child(2) {
      display: none !important;
   }
   .-view:first-of-type {
     button:last-of-type {
         display: none !important;
      }
   }
}

.not-found {
   display: inline-block;
   padding: 20px 50px;
   margin: 4% auto 0 auto;
   h2 {
      font-size: 1.5em;
      color: $uva-text-color-base;
   }
}
.wsls {
   max-width: 640px;
   margin: 0 auto;
   video {
      width: 100%;
   }
   .overview {
      text-align: left;
      h3 {
         margin:  5px;
      }
   }
   .duration {
      text-align: right;
      font-size:0.8em;
      margin: 5px 0 0 0;
   }
   .anchorscript-container {
      margin: 0 auto; text-align: left;
      color: $uva-text-color-base;
      h4 {
        margin: 0 0 10px 0;
        border-bottom: 1px solid $uva-text-color-base;
        padding-bottom: 2px;
      }
      img {
        float: left;
      }
      .anchorscript-links {
         float: left; padding-left: 10px;
         a {
            display: block;
            font-size: 1.1em;
            text-decoration: none;
            color: $uva-blue-alt-A;
            &:hover {
               text-decoration: underline;
            }
         }
      }
   }
}
</style>
