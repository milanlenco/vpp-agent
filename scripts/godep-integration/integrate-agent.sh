#!/bin/bash

set -e

BLUE='\033[0;34m'

# checkout agent code
mkdir -p $GOPATH/src/github.com/ligato
cd $GOPATH/src/github.com/ligato
git clone https://github.com/ligato/vpp-agent
cd vpp-agent
git checkout $AGENT_COMMIT

/integration/continue/integrate-cont.sh
