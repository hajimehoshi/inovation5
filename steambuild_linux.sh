set -e

name=innovation2007
STEAM_RUNTIME_VERSION=0.20210721.1
GO_VERSION=$(go env GOVERSION)
mkdir -p bin/${STEAM_RUNTIME_VERSION}

cd bin

# Download binaries for 386.
if [ ! -f ${STEAM_RUNTIME_VERSION}/com.valvesoftware.SteamRuntime.Sdk-i386-scout-sysroot.Dockerfile ]; then
    curl -L --output ${STEAM_RUNTIME_VERSION}/com.valvesoftware.SteamRuntime.Sdk-i386-scout-sysroot.Dockerfile https://repo.steampowered.com/steamrt-images-scout/snapshots/${STEAM_RUNTIME_VERSION}/com.valvesoftware.SteamRuntime.Sdk-i386-scout-sysroot.Dockerfile
fi
if [ ! -f ${STEAM_RUNTIME_VERSION}/com.valvesoftware.SteamRuntime.Sdk-i386-scout-sysroot.tar.gz ]; then
    curl -L --output ${STEAM_RUNTIME_VERSION}/com.valvesoftware.SteamRuntime.Sdk-i386-scout-sysroot.tar.gz https://repo.steampowered.com/steamrt-images-scout/snapshots/${STEAM_RUNTIME_VERSION}/com.valvesoftware.SteamRuntime.Sdk-i386-scout-sysroot.tar.gz
fi
if [ ! -f ${GO_VERSION}.linux-386.tar.gz ]; then
    curl -L --output ${GO_VERSION}.linux-386.tar.gz https://golang.org/dl/${GO_VERSION}.linux-386.tar.gz
fi

# Download binaries for amd64.
if [ ! -f ${STEAM_RUNTIME_VERSION}/com.valvesoftware.SteamRuntime.Sdk-amd64,i386-scout-sysroot.Dockerfile ]; then
    curl -L --output ${STEAM_RUNTIME_VERSION}/com.valvesoftware.SteamRuntime.Sdk-amd64,i386-scout-sysroot.Dockerfile https://repo.steampowered.com/steamrt-images-scout/snapshots/${STEAM_RUNTIME_VERSION}/com.valvesoftware.SteamRuntime.Sdk-amd64,i386-scout-sysroot.Dockerfile
fi
if [ ! -f ${STEAM_RUNTIME_VERSION}/com.valvesoftware.SteamRuntime.Sdk-amd64,i386-scout-sysroot.tar.gz ]; then
    curl -L --output ${STEAM_RUNTIME_VERSION}/com.valvesoftware.SteamRuntime.Sdk-amd64,i386-scout-sysroot.tar.gz https://repo.steampowered.com/steamrt-images-scout/snapshots/${STEAM_RUNTIME_VERSION}/com.valvesoftware.SteamRuntime.Sdk-amd64,i386-scout-sysroot.tar.gz
fi
if [ ! -f ${GO_VERSION}.linux-amd64.tar.gz ]; then
    curl -L --output ${GO_VERSION}.linux-amd64.tar.gz https://golang.org/dl/${GO_VERSION}.linux-amd64.tar.gz
fi

# Build for 386
(cd ${STEAM_RUNTIME_VERSION}; docker build -f com.valvesoftware.SteamRuntime.Sdk-i386-scout-sysroot.Dockerfile -t steamrt_scout_i386:latest .)
docker run --rm --workdir=/work --volume $(pwd)/..:/work steamrt_scout_i386:latest /bin/sh -c "
export PATH=\$PATH:/usr/local/go/bin
export CGO_CFLAGS=-std=gnu99

rm -rf /usr/local/go && tar -C /usr/local -xzf bin/${GO_VERSION}.linux-386.tar.gz

go build -tags=steam -o bin/${name}_linux_386 .
"

# Build for amd64
(cd ${STEAM_RUNTIME_VERSION}; docker build -f com.valvesoftware.SteamRuntime.Sdk-amd64,i386-scout-sysroot.Dockerfile -t steamrt_scout_amd64:latest .)
docker run --rm --workdir=/work --volume $(pwd)/..:/work steamrt_scout_amd64:latest /bin/sh -c "
export PATH=\$PATH:/usr/local/go/bin
export CGO_CFLAGS=-std=gnu99

rm -rf /usr/local/go && tar -C /usr/local -xzf bin/${GO_VERSION}.linux-amd64.tar.gz

go build -tags=steam -o bin/${name}_linux_amd64 .
"
