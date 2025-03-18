#!/usr/bin/env bash
 
hotfix() {
    sudo sed -i 's/proxies=self.proxy/mounts=self.proxy/g' /usr/lib/python3/dist-packages/nala/downloader.py
}
