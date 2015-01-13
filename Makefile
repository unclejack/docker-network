PKG:=github.com/docker/docker-network
PKG_PATH:=/go/src/$(PKG)
IMGNAME:=dockernetwork
TEST_PKGS:=${PKG}/drivers/simplebridge ${PKG}/namespace
RUN_CMD:=docker run --privileged --rm -e TEST_PKGS="$(TEST_PKGS)" -v "$(shell pwd)":$(PKG_PATH) -w $(PKG_PATH) $(IMGNAME)
RUN_CMD:=docker run --privileged --rm -e TEST_PKGS="$(TEST_PKGS)" -v "$(shell pwd)":$(PKG_PATH) -w $(PKG_PATH) -it $(IMGNAME)

build: dockerbuild
	$(RUN_CMD) project/build.sh

test: dockerbuild
	$(RUN_CMD) project/test.sh

shell: dockerbuild
	$(RUN_CMD) bash

dockerbuild:
	docker build -t $(IMGNAME) .	
