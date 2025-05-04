#!/bin/bash

VERSION="0.0.4"

LLVM_ARCH=$(llvm-config --host-target)
COMMIT_HASH=$(git rev-parse --short=7 HEAD)

cat <<EOF > cmd/modo/version.yaml
version: $VERSION
arch: $LLVM_ARCH
commit: $COMMIT_HASH
EOF
