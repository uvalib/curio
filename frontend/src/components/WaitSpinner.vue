<template>
   <div v-if="overlay" class="spinner-overlay">
      <div class="spinner">
         <h3 v-if="message">{{message}}</h3>
         <div class="spinner-animation">
            <div class="bounce1" :style="{backgroundColor: color}"></div>
            <div class="bounce2" :style="{backgroundColor: color}"></div>
            <div class="bounce3" :style="{backgroundColor: color}"></div>
         </div>
      </div>
   </div>
   <div v-else class="spinner embed">
      <h3 v-if="message">{{message}}</h3>
      <div class="spinner-animation">
         <div class="bounce1" :style="{backgroundColor: color}"></div>
         <div class="bounce2" :style="{backgroundColor: color}"></div>
         <div class="bounce3" :style="{backgroundColor: color}"></div>
      </div>
   </div>
</template>

<script>
export default {
   name: "WaitSpinner",
   props: {
      message: {
         type: String,
         default: ""
      },
      overlay: {
         type: Boolean,
         default: false
      },
      color: {
         type: String,
         default: "var(--uvalib-brand-orange)"
      }
   },
   computed: {
      backgroundColor() {
         if ( this.transparent) return "transparent"
         return "white"
      }
   }
}
</script>

<style lang="scss" scoped>
div.spinner-overlay {
   position: fixed;
   left: 0;
   top: 0;
   width: 100%;
   height: 100%;
   z-index: 1000;
   background:white;
}
div.spinner {
   background: white;
   margin: 12vw auto;
   text-align: center;
   padding: 10px 150px 25px 150px;
   display: inline-block;
   font-weight: bold;
   color: var(--uvalib-text);
}
div.spinner.embed {
   box-shadow: none;
   padding: 0;
   margin: 0;
   background: transparent;
}
@media only screen and (min-width: 768px) {
   div.spinner {
      padding: 40px 90px;
   }
}
@media only screen and (max-width: 768px) {
   div.spinner {
      width: 95%;
      padding: 40px 0;
      margin-top: 30%;
   }
}
div.spinner h1 {
  color: var(--uvalib-text);
  border: none;
}
.spinner-animation {
  margin: 0 auto;
  width: 80px;
  text-align: center;
}
.spinner-animation > div {
  width: 18px;
  height: 18px;
  border-radius: 100%;
  display: inline-block;
  -webkit-animation: sk-bouncedelay 1.4s infinite ease-in-out both;
  animation: sk-bouncedelay 1.4s infinite ease-in-out both;
  margin: 0 2px;
}
.spinner-animation .bounce1 {
  -webkit-animation-delay: -0.32s;
  animation-delay: -0.32s;
}
.spinner-animation .bounce2 {
  -webkit-animation-delay: -0.16s;
  animation-delay: -0.16s;
}
@-webkit-keyframes sk-bouncedelay {
  0%, 80%, 100% { -webkit-transform: scale(0) }
  40% { -webkit-transform: scale(1.0) }
}
@keyframes sk-bouncedelay {
  0%, 80%, 100% {
    -webkit-transform: scale(0);
    transform: scale(0);
  } 40% {
    -webkit-transform: scale(1.0);
    transform: scale(1.0);
  }
}
</style>
