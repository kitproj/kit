#!/usr/bin/env sh
set -eu

os=$(uname | tr '[A-Z]' '[a-z]')
arch=$(uname -m | sed 's/aarch64/arm64/')
tag=${tag:-}

if [ "$tag" = "" ]; then
  curl --fail --retry 5 --silent --location https://api.github.com/repos/kitproj/kit/releases/latest --output /tmp/latest

  tag=$(jq -r .tag_name /tmp/latest)
fi

version=$(echo $tag | cut -c 2-)

url="https://github.com/kitproj/kit/releases/download/${tag}/kit_${os}_${arch}"

sudo curl --fail --retry 5 --silent --location $url --output /usr/local/bin/kit
sudo chmod +x /usr/local/bin/kit
