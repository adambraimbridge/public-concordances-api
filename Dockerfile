FROM alpine:3.3
ADD *.go /public-concordances-api/
ADD people/*.go /public-concordances-api/people/
RUN apk add --update bash \
  && apk --update add git bzr gcc \
  && apk --update add go \
  && export GOPATH=/gopath \
  && REPO_PATH="github.com/Financial-Times/public-concordances-api" \
  && mkdir -p $GOPATH/src/${REPO_PATH} \
  && cp -r public-concordances-api/* $GOPATH/src/${REPO_PATH} \
  && cd $GOPATH/src/${REPO_PATH} \
  && go get -t ./... \
  && cd $GOPATH/src/${REPO_PATH} \
  && go build  \
  && mv public-concordances-api /app \
  && apk del go git bzr \
  && rm -rf $GOPATH /var/cache/apk/*
CMD exec /app --neo-url=$NEO_URL --port=$APP_PORT --graphiteTCPAddress=$GRAPHITE_ADDRESS --graphitePrefix=$GRAPHITE_PREFIX --logMetrics=$LOG_METRICS --cache-duration=$CACHE_DURATION