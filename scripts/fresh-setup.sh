#!/bin/bash
set -e

echo "--- Installing dependencies ---"
sudo apt update
sudo apt install -y docker.io git make curl python3

echo "--- Adding user to docker group ---"
sudo usermod -aG docker $USER

echo "--- Cloning repo ---"
git clone https://github.com/googleboy-byte/goboxd
cd goboxd
git checkout team/silverex

echo "--- Done. Log out and back in for docker group to take effect ---"
echo "--- Then run: cd goboxd && make run ---"
