# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Kit < Formula
  desc "Crazy fast local dev loop."
  homepage "https://github.com/kitproj/kit"
  version "0.0.64"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/kitproj/kit/releases/download/v0.0.64/kit_0.0.64_Darwin_x86_64.tar.gz"
      sha256 "1995d60971dbe2767faf95cd5afd9f505330fbed942b2f1432d6dba60f36b8b0"

      def install
        bin.install "kit"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/kitproj/kit/releases/download/v0.0.64/kit_0.0.64_Darwin_arm64.tar.gz"
      sha256 "48b3087fd9314f509e744095a0ad44402a5103b39c535ec0c213c910e8cfb32a"

      def install
        bin.install "kit"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/kitproj/kit/releases/download/v0.0.64/kit_0.0.64_Linux_arm64.tar.gz"
      sha256 "46dcd7051b786d4edc42d98511843fa6870c8d3ca844204ac69cb29a6e9f7b3c"

      def install
        bin.install "kit"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/kitproj/kit/releases/download/v0.0.64/kit_0.0.64_Linux_x86_64.tar.gz"
      sha256 "30a471c14c42e18fa36f0f4ed7e94ad6f13775f636f2dff2d384dddb285ff1ec"

      def install
        bin.install "kit"
      end
    end
  end
end
