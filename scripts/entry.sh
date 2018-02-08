# set blank options variables
DBHOST_OPT=""
DBNAME_OPT=""
DBUSER_OPT=""
DBPASSWD_OPT=""
IIFURL_OPT=""

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
   IIFURL_OPT="--iiif $IIIF_URL"
fi

cd bin
./digital-object-viewer $DBHOST_OPT $DBNAME_OPT $DBUSER_OPT $DBPASSWD_OPT $IIFURL_OPT

#
# end of file
#
