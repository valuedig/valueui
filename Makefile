# Copyright 2019 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
#


BUILDCOLOR="\033[34;1m"
BINCOLOR="\033[37;1m"
ENDCOLOR="\033[0m"

PROTOC_CMD = protoc
PROTOC_ARGS = --proto_path=./api/ --go_opt=paths=source_relative --go_out=./pkg/uiapi/ --go-grpc_out=./pkg/uiapi/ ./api/uiapi.proto

HTOML_TAG_FIX_CMD = htoml-tag-fix
HTOML_TAG_FIX_ARGS = pkg/uiapi/uiapi.pb.go


ifndef V
	QUIET_BUILD = @printf '%b %b\n' $(BUILDCOLOR)BUILD$(ENDCOLOR) $(BINCOLOR)$@$(ENDCOLOR) 1>&2;
	QUIET_INSTALL = @printf '%b %b\n' $(BUILDCOLOR)INSTALL$(ENDCOLOR) $(BINCOLOR)$@$(ENDCOLOR) 1>&2;
endif

# npm install -g sass
CSS_BUILD_CMD = sass
CSS_BUILD_ARGS = --no-source-map scss/main.scss:main.css

.PHONY: api

all: build_main
	@echo ""
	@echo "build complete"
	@echo ""

build_main:
	$(QUIET_BUILD)$(CSS_BUILD_CMD) $(CSS_BUILD_ARGS) $(CCLINK)

clean:
	@echo ""
	@echo "clean complete"
	@echo ""

api:
	$(QUIET_BUILD)$(PROTOC_CMD) $(PROTOC_ARGS) $(CCLINK)
	$(QUIET_BUILD)$(HTOML_TAG_FIX_CMD) $(HTOML_TAG_FIX_ARGS) $(CCLINK)

