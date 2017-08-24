#!/bin/bash

set -e

BLUE='\033[0;34m'
EOL='\033[0m'

function checkout_agent {
    echo -e "${BLUE}[1/7] Checking out vpp-agent${EOL}"
    mkdir -p ${GOPATH}/src/github.com/ligato
    cd ${GOPATH}/src/github.com/ligato
    git clone https://github.com/ligato/vpp-agent
    cd vpp-agent
    git checkout ${AGENT_COMMIT}
    cd $HOME
}

function import_godep {
    echo -e "${BLUE}[2/7] Importing godep dependencies into glide.yaml${EOL}"
    cd ${GOPATH}/src/${TARGET}
    rm -rf glide.yaml glide.lock
    glide init --non-interactive 
    cd $HOME
}

function resolve_deps {
    echo -e "${BLUE}[3/7] Resolving dependencies${EOL}"
    python /integration/scripts/resolve-deps.py -a ${GOPATH}/src/github.com/ligato/vpp-agent -t ${GOPATH}/src/${TARGET}
}

function download_deps {
    echo -e "${BLUE}[4/7] Downloading dependencies using glide${EOL}"
    cd ${GOPATH}/src/${TARGET}
    rm -rf vendor
    glide --yaml glide-resolved.yaml install --strip-vendor
    cd $HOME
}

function install_deps {
    echo -e "${BLUE}[5/7] Installing dependencies under the GOPATH${EOL}"
    rm -rf ${GOPATH}/src/github.com/ligato # remove checked-out agent to simplify the copy-merge
    rsync -aP ${GOPATH}/src/${TARGET}/vendor/ ${GOPATH}/src/
    # TODO: will need to use goget instead, shit
}

function godep_save {
    echo -e "${BLUE}[6/7] Running godep to save the dependencies${EOL}"
    cd ${GOPATH}/src/${TARGET}
    rm -rf vendor
    godep update -goversion
    godep save # TODO: specify what to build
    cd $HOME
}

function cleanup {
    echo -e "${BLUE}[7/7] Cleaning up${EOL}"
    cd ${GOPATH}/src/${TARGET}
    rm -rf glide.yaml glide-resolved.yaml
    cd $HOME
}


echo -e "${BLUE}Integrating vpp-agent into: ${TARGET}${EOL}"

checkout_agent
import_godep
resolve_deps
download_deps
install_deps
#godep_save
#cleanup

/bin/bash
