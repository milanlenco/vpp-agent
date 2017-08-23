#!/bin/bash


AGENT_COMMIT=`git rev-parse HEAD`
echo "Repo agent commit number: "$AGENT_COMMIT

while [ "$1" != "" ]; do
    case $1 in
        -a | --agent )          shift
                                AGENT_COMMIT=$1
                                ;;
        -t | --target )         shift
                                TARGET=$1
                                ;;
        * )                     echo "invalid parameter "$1
                                exit 1
    esac
    shift
done

if [ -z "$TARGET" ] || [ -z "$AGENT_COMMIT" ]; then
    echo "USAGE: integrate.sh --target <target-package> [--agent <agent-commit>]"
    exit 1
fi


echo "Integrate agent commit number: $AGENT_COMMIT, into the package: $TARGET"

#docker build -t godep_vpp_agent_integration --build-arg AGENT_COMMIT=$AGENT_COMMIT --build-arg TARGET=$TARGET --no-cache .
docker run -it --rm --name ${TARGET##*/}_vpp_agent_integration \
    -v ${GOPATH}/src/$TARGET:/integration/go/src/$TARGET \
    -v ${PWD}/continue:/integration/continue \
    -v ${PWD}/.resolve-deps.cache:/integration/continue/.resolve-deps.cache godep_vpp_agent_integration
