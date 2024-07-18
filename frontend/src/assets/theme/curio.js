import { definePreset } from '@primevue/themes'
import Aura from '@primevue/themes/aura'
import ripple from '@primevue/themes/aura/ripple'
import tooltip from '@primevue/themes/aura/tooltip'
import './uva-colors.css'

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
         50: 'var(--uvalib-brand-blue-lightest)',
         100: 'var(--uvalib-brand-blue-lighter)',
         200: 'var(--uvalib-brand-blue-lighter)',
         300: 'var(--uvalib-brand-blue-lighter)',
         400: 'var(--uvalib-brand-blue-light)',
         500: 'var(--uvalib-brand-blue-light)',
         600: 'var(--uvalib-brand-blue-light)',
         700: 'var(--uvalib-brand-blue-light)',
         800: 'var(--uvalib-brand-blue)',
         900: 'var(--uvalib-brand-blue)',
         950: 'var(--uvalib-brand-blue)'
      },
      focusRing: {
         width: '2px',
         style: 'dotted',
         offset: '3px'
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
               color: 'var(--uvalib-text)',
               focusColor: '#ffffff'
            }
         },
      }
   },
   components: {
      button: {
         colorScheme: {
            light: {
               secondary: {
                  background: 'var(--uvalib-grey-lightest)',
                  hoverBackground: 'var(--uvalib-grey-light)',
                  hoverBorderColor: 'var(--uvalib-grey)',
                  borderColor: 'var(--uvalib-grey-light)',
                  color: 'var(--uvalib-text)',
               },
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