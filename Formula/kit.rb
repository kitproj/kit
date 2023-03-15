# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Kit < Formula
  desc "Crazy fast local dev loop."
  homepage "https://github.com/kitproj/kit"
  version "0.1.1"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/kitproj/kit/releases/download/v0.1.1/kit_0.1.1_Darwin_x86_64.tar.gz"
      sha256 "39caccf4a4ef1746352a4751c239689dcb2e00c553052785689a43edefea61e3"

      def install
        bin.install "kit"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/kitproj/kit/releases/download/v0.1.1/kit_0.1.1_Darwin_arm64.tar.gz"
      sha256 "2a0a20fbb8baf6aa3b66a51f4092e5bf0933e956832c153c4f5f14a151e9db13"

      def install
        bin.install "kit"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/kitproj/kit/releases/download/v0.1.1/kit_0.1.1_Linux_arm64.tar.gz"
      sha256 "4d4015e1f3cf24acd6deebea5a8b855fa99b0def5204b871d4c4521f31bc30b8"

      def install
        bin.install "kit"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/kitproj/kit/releases/download/v0.1.1/kit_0.1.1_Linux_x86_64.tar.gz"
      sha256 "91869e3b8dff58d29d04fd3b39f0cd48485eff17525023f60dc53048d45a2161"

      def install
        bin.install "kit"
      end
    end
  end
end
