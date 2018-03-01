# set blank options variables
DBHOST_OPT=""
DBNAME_OPT=""
DBUSER_OPT=""
DBPASSWD_OPT=""
IIIFURL_OPT=""
DOVHOST_OPT=""

# database host
if [ -n "$DBHOST" ]; then
   DBHOST_OPT="--dbhost $DBHOST"
fi

# database name
if [ -n "$DBNAME" ]; then
   DBNAME_OPT="--dbname $DBNAME"
fi

# database user
if [ -n "$DBUSER" ]; then
   DBUSER_OPT="--dbuser $DBUSER"
fi

# database password
if [ -n "$DBPASSWD" ]; then
   DBPASSWD_OPT="--dbpass $DBPASSWD"
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
./digital-object-viewer $DBHOST_OPT $DBNAME_OPT $DBUSER_OPT $DBPASSWD_OPT $IIIFURL_OPT $DOVHOST_OPT

#
# end of file
#
