#!/usr/bin/env bash

VERSION=${TRAVIS_TAG}


package=cerberus-cli

package_split=(${package//\// })
package_name=${package_split[*]: -1}

platforms=("windows/amd64" "darwin/amd64" "linux/amd64")
ldflag="-X cerberus-cli/cmd.version=${VERSION}"

echo 'building executable for local machine'
go build -ldflags "$ldflag"

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name='./dist/'${package_name}'-'${GOOS}'-'${GOARCH}
    if [[ ${GOOS} = "windows" ]]; then
        output_name+='.exe'
    fi
    echo "building executable for ${GOOS}-${GOARCH}"
    env GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build -ldflags "${ldflag}" -o ${output_name} ${package}
    if [[ $? -ne 0 ]]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
