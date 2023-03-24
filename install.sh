#!/usr/bin/env sh
set -eux

tag=
while [ "$tag" = "" ]; do
  tag=$(curl -fsL "https://api.github.com/repos/kitproj/kit/releases/latest" | jq -r '.tag_name')
done

version=$(echo $tag | cut -c 2-)
url="https://github.com/kitproj/kit/releases/download/${tag}/kit_${version}_$(uname)_$(uname -m | sed 's/aarch64/arm64/').tar.gz"

while [ ! -e kit ]; do
  curl -fsL $url | tar -zxvf - kit
done

chmod +x kit
sudo mv kit /usr/local/bin/kit
