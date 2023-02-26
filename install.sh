#!/usr/bin/env sh
set -eux

tag=$(curl -s "https://api.github.com/repos/kitproj/kit/releases/latest" | jq -r '.tag_name')
version=$(echo $tag | cut -c 2-)
url="https://github.com/kitproj/kit/releases/download/${tag}/kit_${version}_$(uname)_$(uname -m | sed 's/aarch64/arm64/').tar.gz"
curl -L $url | tar -zxvf - kit
chmod +x kit
sudo mv kit /usr/local/bin/kit
