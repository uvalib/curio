# set blank options variables
TRACKSYSURL_OPT=""
IIIFURL_OPT=""
DOVHOST_OPT=""

# TRACKSYS URL
if [ -n "$TRACKSYS_URL" ]; then
   TRACKSYSURL_OPT="--tracksys $TRACKSYS_URL"
fi

# IIIF URL
if [ -n "$IIIF_URL" ]; then
   IIIFURL_OPT="--iiif $IIIF_URL"
fi

# DOV HOST
if [ -n "$DOVHOST" ]; then
   DOVHOST_OPT="--dovhost $DOVHOST"
fi

cd bin
./digital-object-viewer $TRACKSYSURL_OPT $IIIFURL_OPT $DOVHOST_OPT

#
# end of file
#
