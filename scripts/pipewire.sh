#!/usr/bin/env bash

mkdir -p /tmp/audiofix && cd /tmp/audiofix
deps=(
  "libcanberra-pulse"
  "libconfig++9v5"
  "libffado2"
  "libglibmm-2.4-1t64"
  "libxml++2.6-2v5"
  "libpipewire-0.3-common"
  "pipewire-audio"
)
pacdeps=(
  "libwebrtc-audio-processing1-deb"
  "libspa-0.2-modules-deb"
  "libspa-0.2-bluetooth-deb"
  "libpipewire-0.3-0t64-deb"
  "libpipewire-0.3-modules-deb"
  "libpipewire-0.3-modules-x11-deb"
  "pipewire-bin-deb"
  "pipewire-deb"
  "pipewire-alsa-deb"
  "pipewire-jack-deb"
  "pipewire-pulse-deb"
  "pipewire-v4l2-deb"
  "gstreamer1.0-pipewire-deb"
)
sudo apt update && sudo apt install "${deps[@]}" -y
echo N | for i in "${pacdeps[@]}"; do pacstall -QaB "${i}"#6926; done
sudo dpkg -i *.deb
cd ~/
sudo rm -r /tmp/audiofix
