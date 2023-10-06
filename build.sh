#!/bin/bash

# Ensure that script stops on first error
set -e

# Ensure Go is installed
if ! command -v go &> /dev/null
then
    echo "Go could not be found, install it to continue"
    exit
fi

# Validate arguments
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <version_tag>"
    exit 1
fi

VERSION=$1

# Get the latest tag as version
# VERSION=$(git describe --tags $(git rev-list --tags --max-count=1))
# if [ -z "$VERSION" ]; then
#     echo "No git tags found, please create a tag and try again."
#     exit 1
# fi

OUTPUT_DIR="releases/download/${VERSION}"
PROJECT_PATH="main.go"  # adjust to your project's main package
PROJECT_NAME="canopy-cli"       # adjust to your project's output binary name

# Supported GOOS/GOARCH pairs
PLATFORMS="darwin/amd64 linux/amd64 windows/amd64"

# Ensure clean build dir
rm -rf $OUTPUT_DIR
mkdir -p $OUTPUT_DIR

# Cross-compiling for multiple platforms
for PLATFORM in $PLATFORMS; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    
    # Define output path
    OUTPUT="${OUTPUT_DIR}/${PROJECT_NAME}"
    
    # Add exe extension for Windows
    if [ "$GOOS" == "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
    fi
    
    echo "Building for $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build -o $OUTPUT $PROJECT_PATH
    
    # Zip the binary with the version in the filename
    zip "${OUTPUT_DIR}/${PROJECT_NAME}_${VERSION}_${GOOS}_${GOARCH}.zip" $OUTPUT
    
    # Remove the binary after zipping
    rm $OUTPUT
done

# At this point, binaries should be zipped in the "build" directory
# If you want to upload to GitHub releases, consider using GitHub Actions or gh cli
# gh release create $VERSION build/*.zip

git checkout main
git add .
git commit -m "build: ${VERSION} release binaries" 
git tag -d "${VERSION}"
git tag "${VERSION}"
