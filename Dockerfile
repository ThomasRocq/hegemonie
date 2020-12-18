
FROM golang:1.15-buster AS builder

LABEL maintainers "Jean-Francois SMIGIELSKI <jf.smigielski@gmail.com>"

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOPATH=/gopath

WORKDIR /dist

RUN set -x \
&& apt-get update -y \
&& apt-get install -y --no-install-recommends \
  make \
  librocksdb-dev \
  librocksdb5.17

COPY Makefile LICENSE AUTHORS.md go.sum go.mod \
  /gopath/src/github.com/jfsmig/hegemonie/
COPY pkg /gopath/src/github.com/jfsmig/hegemonie/pkg
COPY api /gopath/src/github.com/jfsmig/hegemonie/api
COPY cmd /gopath/src/github.com/jfsmig/hegemonie/cmd

# Build & Install the code
RUN set -x \
&& cd /gopath/src/github.com/jfsmig/hegemonie \
&& go mod download

RUN set -x \
&& cd /gopath/src/github.com/jfsmig/hegemonie \
&& make \
&& cp -v /gopath/bin/hege /gopath/bin/heged /dist

# Install the dependencies.
# Inspired by https://dev.to/ivan/go-build-a-minimal-docker-image-in-just-three-steps-514i
RUN set -x \
&& mkdir -p /dist/lib64 \
&& ldd ./hegemonie | tr -s '[:blank:]' '\n' | grep '^/' | \
   xargs -I % sh -c 'mkdir -p $(dirname ./%); cp % ./%;' \
&& cp /lib64/ld-linux-x86-64.so.2 /dist/lib64/

# Mangle the maps to build ther raw shape based on the seed definitions
# JFS: we do this here because it is very fast to execute and it benefits
#      from the rich shell environments.
COPY docs/maps        /data/maps
COPY docs/definitions /data/defs
COPY docs/lang        /data/lang
RUN set -ex \
&& D=/data/maps \
&& HEGE=/gopath/bin/hege \
&& ls $D | \
   grep '.seed.json$' | \
   while read F ; do echo $F $F ; done | \
   sed -r 's/^(\S+).seed.json /\1.final.json /' | \
   while read FINAL SEED ; do \
    echo $D $SEED $FINAL ; \
    $HEGE map tools init < $D/$SEED | $HEGE map tools normalize > $D/$FINAL ; \
  done



#------------------------------------------------------------------------------
# Create the minimal runtime image

FROM scratch as runtime
COPY --chown=0:0 --from=builder /dist /
# Expose each default port of each module. --> pkg/utils/constants.go
EXPOSE 8080/tcp
EXPOSE 8081/tcp
EXPOSE 8082/tcp
EXPOSE 8083/tcp
EXPOSE 8084/tcp
USER 0
WORKDIR /
ENTRYPOINT ["/heged"]



#------------------------------------------------------------------------------
# Create a bundle with minimal runtime plus demonsration data

FROM runtime as demo
COPY --chown=65534:0 --from=builder /data/defs/        /data/defs/
COPY --chown=65534:0 --from=builder /data/lang/        /data/lang/
COPY --chown=65534:0 --from=builder /data/maps/        /data/maps/
COPY --chown=65534:0                pkg/web/templates  /data/templates/
COPY --chown=65534:0                pkg/web/static     /data/static/
USER 65534
WORKDIR /data
ENTRYPOINT ["/heged"]

