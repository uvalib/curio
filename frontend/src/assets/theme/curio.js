import { definePreset } from '@primeuix/themes'
import Aura from '@primeuix/themes/aura'
import ripple from '@primeuix/themes/aura/ripple'
import tooltip from '@primeuix/themes/aura/tooltip'
import colors from './colors.module.scss'

const Curio = definePreset(Aura, {
   root: {
      borderRadius: {
         none: '0',
         xs: '2px',
         sm: '3px',
         md: '4px',
         lg: '4px',
         xl: '8px'
      },
   },
   semantic: {
      primary: {
         50:  colors.brandBlue300,
         100: colors.brandBlue300,
         200: colors.brandBlue300,
         300: colors.brandBlue300,
         400: colors.brandBlue100,
         500: colors.brandBlue100,
         600: colors.brandBlue100,
         700: colors.brandBlue100,
         800: colors.brandBlue,
         900: colors.brandBlue,
         950: colors.brandBlue
      },
      focusRing: {
         width: '2px',
         style: 'solid',
         offset: '2px',
      },
      disabledOpacity: '0.3',
      colorScheme: {
         light: {
            primary: {
               color: '{primary.500}',
               contrastColor: '#ffffff',
               hoverColor: '{primary.100}',
               activeColor: '{primary.500}'
            },
            highlight: {
               background: '#ffffff',
               focusBackground: '#ffffff',
               color: colors.textBase,
               focusColor: '#ffffff'
            }
         },
      }
   },
   components: {
      button: {
         colorScheme: {
            light: {
               root: {
                  secondary: {
                     background: colors.grey200,
                     hoverBackground: colors.grey100,
                     hoverBorderColor: colors.grey,
                     borderColor: colors.grey100,
                     color: colors.textBase,
                  },
                  contrast:  {
                     focusRing: {
                        color: '#99c8ff',
                        shadow: 'none'
                     }
                  },
                  primary: {
                     focusRing: {
                        color: colors.blueAlt300,
                     }
                  },
               },
               text: {
                  primary: {
                     color: 'white',
                     hoverBackground: '#0003',
                     activeBackground: '#0003'
                  }
               }
            }
         }
      },
      toast: {
         colorScheme: {
            light: {
               success: {
                  background: '{green.200}',
                  borderColor: '{green.700}',
               },
               error: {
                  background: '{red.100}',
                  borderColor: '{ref.400}',
               }
            }
         }
      }
   },
   directives: {
      tooltip,
      ripple
   }
});

export default Curio;