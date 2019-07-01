#!/usr/bin/env bash

# Check.
go env
go version

# Determine the arch/os combos we're building for.
XC_ARCH=${XC_ARCH:-"amd64"}
XC_OS=${XC_OS:-linux darwin windows}

# Build!
for OS in ${XC_OS[@]}; do
  for ARCH in ${XC_ARCH[@]}; do
    echo "==> Building binary for $OS/$ARCH..."
    out="dist/"$OS"_"$ARCH"/anchor"
    if [[ "${OS}" == "windows" ]]; then
       out="$out.exe"
    fi
    CGO_ENABLED=0 GOARCH=${ARCH} GOOS=${OS} go build -o ${out} -a -ldflags '-extldflags "-static"' .
    zip dist/"${OS}_${ARCH}"/"anchor_${OS}_${ARCH}".zip ${out}
  done
done