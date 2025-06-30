#!/usr/bin/env bash

hotfix() {
  pkgname="nix-bin"
  badver="2.26.3+dfsg-1ubuntu2"
  goodver="2.26.3+dfsg-1+b1"
  CARCH="$(dpkg --print-architecture)"

  mapfile -t status < <(dpkg-query -W --showformat='${db:Status-Status}\n${Version}' "${pkgname}" 2> /dev/null)
  if [[ ${status[0]} == "installed" ]]; then
    if [[ ${status[1]} == "${badver}" ]]; then
      wget -q "http://ftp.debian.org/debian/pool/main/n/nix/${pkgname}_${goodver}_${CARCH}.deb"
      sudo apt install -f ./"${pkgname}_${goodver}_${CARCH}.deb" -y
    else
      echo "Hotfix for '${pkgname}' not compatible or not required!"
    fi
  else
    echo "'${pkgname}' is not installed!"
  fi
}
