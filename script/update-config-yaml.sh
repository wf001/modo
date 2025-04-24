#!/bin/bash

LLVM_ARCH=$(llvm-config --host-target)
VERSION="0.0.1"
COMMIT_HASH=$(git rev-parse --short=7 HEAD)

cat <<EOF > config.yaml
version: "$VERSION"
arch: "$LLVM_ARCH"
commit: "$COMMIT_HASH"
EOF
