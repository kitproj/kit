#!/usr/bin/env sh
set -eu

tag=${tag:-}

if [ "$tag" = "" ]; then
  # retry for all errors up to 10m
  curl --fail --retry 99 --silent --location https://api.github.com/repos/kitproj/kit/releases/latest --output /tmp/latest

  tag=$(jq -r .tag_name /tmp/latest)
fi

version=$(echo $tag | cut -c 2-)

echo "Downloading Kit $version"

url="https://github.com/kitproj/kit/releases/download/${tag}/kit_${version}_$(uname)_$(uname -m | sed 's/aarch64/arm64/').tar.gz"

# retry for all errors up to 10m
curl --fail --retry 99 --silent --location $url --output /tmp/kit.tar.gz
tar -zxvf /tmp/kit.tar.gz kit
chmod +x kit
sudo mv kit /usr/local/bin/kit

echo "Installed /usr/local/bin/kit"