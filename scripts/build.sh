#!/bin/bash

# Version
VERSION="0.0.1"

# Platforms
PLATFORMS=("linux/amd64")

# Binary name
BINARY="janusknock"

mkdir -p builds

# Build for all platforms (specified in list "PLATFORMS")
for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    
    # Output name
    output_name=$BINARY'_'$VERSION'_'$GOOS'_'$GOARCH
    
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    echo "Building for $GOOS/$GOARCH..."
    
    # Run build process
    env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build -o builds/$output_name ./cmd/janusknock
    
    if [ $? -eq 0 ]; then
        echo "✓ Build successful for $GOOS/$GOARCH"
    else
        echo "✗ Build failed for $GOOS/$GOARCH"
        exit 1
    fi
done

# Make binaries executable
chmod +x builds/*

echo "Build proceso completado!"
