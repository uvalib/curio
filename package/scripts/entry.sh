# run application

# run from here, since application expects web template in web/
cd bin; ./curio \
  -apollo $APOLLO_URL \
  -iiif $CURIO_IIIF_MAN_URL \
  -rights $RIGHTS_WRAPPER_URL \
  -host $CURIO_HOST \
  --archivematicaBucket $ARCHIVEMATICA_CURIO_BUCKET \

#
# end of file
#
