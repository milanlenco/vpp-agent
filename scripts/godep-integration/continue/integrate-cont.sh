#!/bin/bash

set -e

BLUE='\033[0;34m'

function import_godep {
	echo -e "${BLUE}[1/X] Importing godep dependencies"
	cd ${GOPATH}/src/${TARGET}
	rm -rf glide.yaml glide.lock
	glide init --non-interactive 
    cd $HOME
}

function resolve_deps {
	echo -e "${BLUE}[2/X] Resolving dependencies"
    python /integration/continue/resolve-deps.py -a ${GOPATH}/src/github.com/ligato/vpp-agent -t ${GOPATH}/src/${TARGET}
}

function install_deps {
	echo -e "${BLUE}[3/X] Installing dependencies using glide"
    cd ${GOPATH}/src/${TARGET}
    rm -rf vendor
    glide --yaml glide-resolved.yaml install --strip-vendor
    mv vendor/* ${GOPATH}/src/
    cd $HOME
}

function godep_save {
    cd ${GOPATH}/src/${TARGET}
    rm -rf vendor
    godep save
    cd $HOME
}

function cleanup {
    cd ${GOPATH}/src/${TARGET}
    rm -rf glide.yaml glide-resolved.yaml
    cd $HOME
}

echo -e "${BLUE}Integrating agent into: "$TARGET

import_godep
resolve_deps
install_deps
godep_save
#cleanup

/bin/bash
