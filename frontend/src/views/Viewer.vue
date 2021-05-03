<template>
   <div class="viewer">
      <WaitSpinner v-if="working" :overlay="true" message="Loading viewer..." />
      <template v-else>
         <template v-if="viewType=='iiif'">
            <div id="tify-viewer"></div>
            <div class="extra-tools">
               <span class="image-download" @click="downloadImage">
                  <i class="fas fa-download"></i>
                  <span class="dl-text">Download image</span>
               </span>
            </div>
         </template>
         <div v-else-if="viewType=='wsls'" class="wsls">
            <div class="overview">
               <h3>{{wslsData.title}}</h3>
               <p>{{wslsData.description}}</p>
            </div>
            <div v-if="wslsData.has_video" class="video-container" >
               <video class="video-js vjs-default-skin vjs-big-play-centered vjs-fluid" controls preload="auto"
                  :poster="wslsData.poster_url" data-setup='{"inactivityTimeout": 0}'
               >
                  <source :src="wslsData.video_url" type='video/mp4'>
                  <p class="vjs-no-js">
                     To view this video please enable JavaScript, and consider upgrading to a web browser that
                     <a href="http://videojs.com/html5-video-support/" target="_blank">supports HTML5 video</a>
                  </p>
               </video>
               <p class="duration">Duration: {{wslsData.duration}}</p>
            </div>
            <div v-if="wslsData.has_script" class="anchorscript-container">
               <h4>Anchor Script</h4>
               <img :src="wslsData.thumb_url"/>
               <div class="anchorscript-links">
                  <a :href="wslsData.pdf_url" target="_blank">View anchor script PDF in new tab</a>
                  <a :href="wslsData.transcript_url" target="_blank">View anchor script transcription in new tab</a>
               </div>
             </div>
         </div>
         <div v-else class="not-found">
            <h2>Sorry, but the resource you requested could not be found.</h2>
         </div>
      </template>
   </div>
</template>

<script>
import { mapState } from "vuex"
export default {
   name: "Viewer",
   computed: {
      ...mapState({
         working : state => state.working,
         iiifURL: state => state.iiifURL,
         pagePIDs: state => state.pagePIDs,
         rightsURL: state => state.rightsURL,
         startPage: state => state.startPage,
         viewType: state => state.viewType,
         wslsData: state => state.wslsData,
      })
   },
   async created() {
      let pid = this.$route.params.pid
      let page = this.$route.query.page
      let unitID = this.$route.query.unit
      if (!page) page = "1"
      await this.$store.dispatch("getPIDViewData", {pid: pid, page: page, unit: unitID})
      window.tifyOptions = {
         container: '#tify-viewer',
         immediateRender: false,
         manifest: this.iiifURL,
         stylesheet: '/tify_mods.css',
         title: null,
      }
      await import ('tify/dist/tify.css')
      await import ('tify/dist/tify.js')
      this.$nextTick( ()=> {
         if (this.startPage > 1) {
            let testQ = Object.assign({}, this.$route.query)
            delete testQ.page
            delete testQ.tify
            let tify = {pages: [this.startPage]}
            testQ.tify = JSON.stringify(tify)
            this.$router.replace({query: testQ})
         }
      })
   },
   methods: {
      downloadImage() {
         let page = 0
         let url = new URL(window.location.href)
         let tifyParamsStr = url.searchParams.get("tify")
         if (tifyParamsStr && tifyParamsStr.length > 0) {
            let tifyParams = JSON.parse(tifyParamsStr)
            if (tifyParams.pages) {
               page = tifyParams.pages[0]-1
            }
         }
         let tgtPID =  this.pagePIDs[page]
         let dlURL = `${this.rightsURL}/${tgtPID}`
         var link = document.createElement('a')
         link.href = dlURL+"?download=1"
         document.body.appendChild(link)
         link.click()
         document.body.removeChild(link)
      }
   }
}
</script>

<style lang="scss">
::v-deep .tify-header_column.-controls.-visible  {
   display: none !important;
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
.wsls {
   max-width: 640px;
   min-width: 410px;
   margin: 0 auto;
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
@media only screen and (max-width: 600px) {
   .dl-text {
      display: none;
   }
}
</style>
