#!/usr/bin/env sh
set -eux

tag=${tag:-}

if [ "$tag" = "" ]; then
  # retry for all errors up to 10m
  curl --retry 99 -vfsL "https://api.github.com/repos/kitproj/kit/releases/latest" -o /tmp/latest

  tag=$(jq -r .tag_name /tmp/latest)
fi

version=$(echo $tag | cut -c 2-)

url="https://github.com/kitproj/kit/releases/download/${tag}/kit_${version}_$(uname)_$(uname -m | sed 's/aarch64/arm64/').tar.gz"

# retry for all errors up to 10m
curl --retry 99 -vfsL $url -o /tmp/kit.tar.gz
tar -zxvf /tmp/kit.tar.gz kit
chmod +x kit
sudo mv kit /usr/local/bin/kit
