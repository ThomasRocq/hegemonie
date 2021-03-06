# Copyright (C) 2020 Hegemonie's AUTHORS
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

BASE=github.com/jfsmig/hegemonie
GO=go
PROTOC=protoc
COV_OUT=coverage.txt

AUTO=
AUTO+= pkg/auth/backend/mem_auto.go
AUTO+= pkg/auth/proto/auth.pb.go
AUTO+= pkg/map/graph/map_auto.go
AUTO+= pkg/map/proto/map.pb.go
AUTO+= pkg/event/proto/event.pb.go
AUTO+= pkg/region/model/world_auto.go
AUTO+= pkg/region/proto/region.pb.go
AUTO+= pkg/healthcheck/healthcheck.pb.go

all: prepare
	$(GO) install $(BASE)/cmd/gen-set
	#$(GO) install $(BASE)/cmd/hege-init
	$(GO) install $(BASE)/cmd/heged
	$(GO) install $(BASE)/cmd/hege

prepare: $(AUTO)

pkg/auth/backend/mem_auto.go: pkg/auth/backend/mem.go cmd/gen-set/gen-set.go
	-rm $@
	$(GO) generate github.com/jfsmig/hegemonie/pkg/auth/backend

pkg/map/graph/map_auto.go: pkg/map/graph/map.go cmd/gen-set/gen-set.go
	-rm $@
	$(GO) generate github.com/jfsmig/hegemonie/pkg/map/graph

pkg/region/model/world_auto.go: pkg/region/model/types.go cmd/gen-set/gen-set.go
	-rm $@
	$(GO) generate github.com/jfsmig/hegemonie/pkg/region/model

pkg/map/proto/%.pb.go: api/map.proto
	$(PROTOC) -I api api/map.proto --go_out=plugins=grpc:pkg/map/proto

pkg/auth/proto/%.pb.go: api/auth.proto
	$(PROTOC) -I api api/auth.proto --go_out=plugins=grpc:pkg/auth/proto

pkg/region/proto/%.pb.go: api/region.proto
	$(PROTOC) -I api api/region.proto  --go_out=plugins=grpc:pkg/region/proto

pkg/event/proto/%.pb.go: api/event.proto
	$(PROTOC) -I api api/event.proto  --go_out=plugins=grpc:pkg/event/proto

pkg/healthcheck/%.pb.go: api/healthcheck.proto
	$(PROTOC) -I api api/healthcheck.proto  --go_out=plugins=grpc:pkg/healthcheck

clean:
	-rm $(AUTO)

.PHONY: all prepare clean test bench fmt try gen-set

fmt:
	go list ./... | grep -v vendor | while read D ; do go fmt $$D ; done

test: all
	go list ./... | while read D ; do go test -race -coverprofile=profile.out -covermode=atomic $$D ; if [ -f profile.out ] ; then cat profile.out >> $(COV_OUT) ; fi ; done

benchmark: all
	go list ./... | while read D ; do go test -race -coverprofile=profile.out -covermode=atomic -bench=$$D $$D ; if [ -f profile.out ] ; then cat profile.out >> $(COV_OUT) ; fi ;  done

img_tag:
	 ( export L='(C) Quentin Minten / CC BY-NC-SA 3.0' ; \
		for F in website/www/static/img0/quentin-minten*/*.jpg ; do \
			BN=$(basename $$F) ; \
			convert img0/$$BN -gravity south -stroke '#000C' -strokewidth 2 -annotate 0 "$L" -stroke  none -fill yellow -annotate 0 "$L" website/www/static/img/$$BN ; \
		done )

docker: Dockerfile
	for T in runtime demo ; do sudo docker build --target=$$T --tag=jfsmig/hegemonie-$$T . ; done
	for T in runtime demo ; do sudo docker push jfsmig/hegemonie-$$T:latest ; done

try: all
	sudo docker-compose up

