#!/usr/bin/env bash

hotfix() {
  local name pkgver base_url packages pacdeps
  name="pipewire-bundle"
  pkgver="1.2.7-1+b2"
  base_url="http://ftp.debian.org/debian/pool/main/p/pipewire"
  packages=(
    "libspa-0.2-modules"
    "libspa-0.2-bluetooth"
    "libpipewire-0.3-0t64"
    "libpipewire-0.3-modules"
    "libpipewire-0.3-modules-x11"
    "pipewire-bin"
    "pipewire"
    "pipewire-alsa"
    "pipewire-jack"
    "pipewire-pulse"
    "pipewire-v4l2"
    "gstreamer1.0-pipewire"
  )

  for i in "${!packages[@]}"; do
    wget -q "${base_url}/${packages[i]}_${pkgver}_${CARCH}.deb" -O "${i}_${name}.deb"
  done

  sudo apt-get install --allow-downgrades -f ./[0-9]*"${name}".deb -y
}
