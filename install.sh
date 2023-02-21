#!/usr/bin/env sh
set -eux

tag=$(curl -s "https://api.github.com/repos/alexec/kit/releases/latest" | jq -r '.tag_name')
version=$(echo $tag | cut -c 2-)
url="https://github.com/alexec/kit/releases/download/${tag}/kit_${version}_$(uname)_$(uname -m).tar.gz"
curl -L $url | tar -xvf - kit
chmod +x kit
sudo mv kit /usr/local/bin/kit
