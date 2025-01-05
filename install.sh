#!/usr/bin/env sh
set -eux

os=$(uname | tr '[A-Z]' '[a-z]')
arch=$(uname -m | sed 's/aarch64/arm64/' | sed 's/x86_64/amd64/')
tag=${tag:-}

if [ "$tag" = "" ]; then
  tag=$(curl --fail --retry 5 --silent --location https://api.github.com/repos/kitproj/kit/releases/latest | jq -r .tag_name)
fi

sudo curl --fail --retry 5 --silent --location --output /usr/local/bin/kit https://github.com/kitproj/kit/releases/download/${tag}/kit_${tag}_${os}_${arch}
sudo chmod +x /usr/local/bin/kit
