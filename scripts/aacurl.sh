#!/usr/bin/env bash

hotfix() {
  sudo apt update && sudo apt install apparmor-utils -y
  sudo aa-disable curl
}
