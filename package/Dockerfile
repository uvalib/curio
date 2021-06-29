#
# build the target application
#
FROM golang:1.16.5-alpine3.13 as builder
RUN apk update && apk upgrade && apk add make yarn

WORKDIR /build
COPY go.mod go.sum Makefile ./
COPY viewsrv ./viewsrv
COPY frontend ./frontend
COPY templates ./templates
RUN make linux-full

#
# build the target container
#
FROM alpine:3.13

# update the packages
RUN apk update && apk upgrade && apk add bash tzdata ca-certificates curl && rm -fr /var/cache/apk/*

# Create the run user and group
RUN addgroup webservice && adduser webservice -G webservice -D

# set the timezone appropriatly
ENV TZ=UTC
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Specify home
ENV APP_HOME /curio
WORKDIR $APP_HOME

# Create necessary directories
RUN mkdir -p $APP_HOME $APP_HOME/scripts $APP_HOME/bin
RUN chown -R webservice $APP_HOME && chgrp -R webservice $APP_HOME

# port and run command
EXPOSE 8080
CMD scripts/entry.sh

# Move in necessary assets
COPY package/scripts/entry.sh $APP_HOME/scripts/entry.sh
COPY package/data/container_bash_profile /home/webservice/.profile
COPY --from=builder /build/bin/curio.linux $APP_HOME/bin/curio
COPY --from=builder /build/bin/public $APP_HOME/bin/public
COPY --from=builder /build/bin/templates $APP_HOME/bin/templates

# Ensure permissions are correct
RUN chown webservice:webservice /home/webservice/.profile $APP_HOME/scripts/entry.sh $APP_HOME/bin/curio && chmod 755 /home/webservice/.profile $APP_HOME/scripts/entry.sh $APP_HOME/bin/curio
RUN chown -R webservice:webservice $APP_HOME/bin/public $APP_HOME/bin/templates && chmod -R 755 $APP_HOME/bin/templates $APP_HOME/bin/public

# Add the build tag
ARG BUILD_TAG
RUN test -n "$BUILD_TAG" && touch $APP_HOME/buildtag.build-$BUILD_TAG || touch $APP_HOME/buildtag.build-0

# Specify the user
USER webservice

#
# end of file
#
