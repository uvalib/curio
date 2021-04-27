<template>
   <div class="vue">
      <div class="header" role="banner">
         <div class="library-link">
            <UvaLibraryLogo v-if="isKiosk" />
            <a v-else target="_blank" href="https://library.virginia.edu">
               <UvaLibraryLogo />
            </a>
         </div>
         <div class="site-link">
            <router-link to="/">
               Curio
            </router-link>
         </div>
      </div>
      <div class="body">
         <h2>Welcome to Curio, the UVA Library digital object viewer</h2>
         <div class="content">
            <p>
               To use this viewer, make a request like the following:
            </p>
            <a :href="`${basePath}view/uva-lib:2218948`" target="_blank">
               {{basePath}}view/uva-lib:2218948
            </a>
            <a :href="`${basePath}view/uva-lib:770581?page=10`" target="_blank">
               {{basePath}}view/uva-lib:770581?page=10
            </a>
            <p class="pad">
               If you need an HTML snippet that can be used to embed the object in another web page, use:
            </p>
            <a :href="`${basePath}oembed?url=${basePath}view/uva-lib:2218948&format=json`" target="_blank">
               {{basePath}}oembed?url={{basePath}}oembed?url={{basePath}}view/uva-lib:2218948&format=json
            </a>
            <a :href="`${basePath}oembed?url=${basePath}view/uva-lib:770581?page=10&format=xml`" target="_blank">
               {{basePath}}oembed?url={{basePath}}view/uva-lib:770581?page=10&format=xml
            </a>
         </div>
      </div>
   </div>
</template>

<script>
import UvaLibraryLogo from "@/components/UvaLibraryLogo"
export default {
   name: "Home",
   components: {
      UvaLibraryLogo
   },
   computed: {
      basePath() {
         let url = window.location.href
         return url
      }
   }
};
</script>

<style scoped lang="scss">
div.header {
   background-color: var(--uvalib-brand-blue);
   color: white;
   padding: 1vw 20px;
   text-align: left;
   position: relative;
   box-sizing: border-box;
   display: flex;
   flex-direction: row;
   flex-wrap: nowrap;
   justify-content: space-between;
   align-content: stretch;
   align-items: center;
}
#app div.header div.site-link a {
  color: white;
  text-decoration: none;
}

div.library-link {
  height: 45px;
  width: 220px;
  order: 0;
  flex: 0 1 auto;
  align-self: flex-start;
}
div.site-link {
  height: 35px;
  width: 70px;
  order: 0;
  flex: 0 1 auto;
  align-self: center;
  font-size: 1.85em;
}
.body {
   padding: 25px;
   h2 {
      font-size: 1.15em;
      color: var(--uvalib-text);
   }
   .content {
      width: 50%;
      margin: 50px auto;
      text-align: left;
      p.pad {
         margin-top: 30px;
      }
      a {
         margin: 10px 20px;
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
</style>
