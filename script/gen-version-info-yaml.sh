#!/bin/bash

VERSION=$(grep -s 'version:' ./cmd/modo/version.yaml | cut -d' ' -f2 )

LLVM_ARCH=$(llvm-config --host-target)
COMMIT_HASH=$(git rev-parse --short=7 HEAD)

cat <<EOF > cmd/modo/version-info.yaml
version: $VERSION
arch: $LLVM_ARCH
commit: $COMMIT_HASH
EOF
