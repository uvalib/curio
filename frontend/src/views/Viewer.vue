<template>
   <div class="viewer">
      <WaitSpinner v-if="curio.working" :overlay="true" message="Loading viewer..." />
      <template v-else>
         <template v-if="curio.viewType=='iiif'">
            <div id="tify-viewer" style="height:100%;"></div>
            <div class="extra-tools">
               <span class="image-download" @click="downloadImage">
                  <i class="fas fa-download"></i>
                  <span class="dl-text">Download image</span>
               </span>
            </div>
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

import 'tify'
import 'tify/dist/tify.css'

const curio = useCurioStore()
const route = useRoute()
const router = useRouter()

const tgtDomain = ref("")
const intervalID = ref(-1)
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
      intervalID.value = setInterval( changeParam, 1000)
   }

   if ( tgtDomain.value) {
      setTimeout(dimensionsMessage, 500)
   }
})

onBeforeUnmount(()=>{
   if ( intervalID.value > -1) {
      clearInterval(intervalID.value)
      intervalID.value = -1
   }
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
         if (q.y ) evt.y = q.y
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
.viewer {
   height: 100%;
}
.extra-tools {
   z-index: 1000;
   position: absolute;
   left: 12px;
   top: 12px;
   font-size: 1.1em;
   color: #222;
   cursor: pointer;
   .dl-text {
      margin-left: 5px;
      font-weight: 500;
      &:hover {
         text-decoration: underline;
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
      color: var(--uvalib-text);
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
      margin: 0 auto; text-align: left; color: var(--uvalib-text);
      h4 {
        margin: 0 0 10px 0;
        border-bottom: 1px solid var(--uvalib-text);
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
            color: var(--color-link);
            &:hover {
               text-decoration: underline;
            }
         }
      }
   }
}
</style>
