if [ -z "$DOCKER_HOST" ]; then
   echo "ERROR: no DOCKER_HOST defined"
   exit 1
fi

# set the definitions
INSTANCE=digobjview-ws
NAMESPACE=uvadave

docker run -ti -p 8300:8085 $NAMESPACE/$INSTANCE /bin/bash -l
