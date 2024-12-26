# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Kit < Formula
  desc "Crazy fast local dev loop."
  homepage "https://github.com/kitproj/kit"
  version "0.1.27"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/kitproj/kit/releases/download/v0.1.27/kit_0.1.27_Darwin_arm64.tar.gz"
      sha256 "610b987c50bcb32818394758d6e3c3144eabfb6e157835ed87c4bbdc78cec975"

      def install
        bin.install "kit"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/kitproj/kit/releases/download/v0.1.27/kit_0.1.27_Darwin_x86_64.tar.gz"
      sha256 "d6c49f9c39b00390fb2335c7255675fda7d0bb56bb3137608ede2ae3ae1e3aa4"

      def install
        bin.install "kit"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/kitproj/kit/releases/download/v0.1.27/kit_0.1.27_Linux_arm64.tar.gz"
      sha256 "4525f84fcaa32921086ec8e5bb5b566f142dd2f4cd3c6f7d9e5b28dd64a37e8f"

      def install
        bin.install "kit"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/kitproj/kit/releases/download/v0.1.27/kit_0.1.27_Linux_x86_64.tar.gz"
      sha256 "e41831b7648be3eec4d44ed4e9c551d37d2ef9e89dea7eb6612a90fd3f99e9e6"

      def install
        bin.install "kit"
      end
    end
  end
end
