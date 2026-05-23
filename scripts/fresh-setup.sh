#!/bin/bash
set -e

echo "--- Installing dependencies ---"
sudo apt-get update
sudo apt-get install -y docker.io git make curl python3 jq

echo "--- Starting Docker ---"
sudo systemctl enable docker
sudo systemctl start docker

echo "--- Adding user to docker group ---"
sudo usermod -aG docker $USER

echo "--- Cloning repo ---"
git clone https://github.com/googleboy-byte/goboxd
cd goboxd
git checkout team/silverex

echo ""
echo "--- Setup complete ---"
echo "Run the following to start:"
echo "  newgrp docker"
echo "  cd goboxd"
echo "  make run"