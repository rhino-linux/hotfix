#!/usr/bin/env bash

hotfix() {
  sudo apt-get update && sudo apt-get install --reinstall libsnappy1v5 -y
}
