#!/bin/bash

set -e
bash build.sh

DOCKER_USER="devkahar99"
IMAGE_NAME="incident-tracker"
TAG="latest"

echo "🚀 Building Docker image..."
docker build -t $DOCKER_USER/$IMAGE_NAME:$TAG .

echo "🔑 Logging in to Docker Hub..."
docker login

echo "📤 Pushing image to Docker Hub..."
docker push $DOCKER_USER/$IMAGE_NAME:$TAG

echo "📦 Deploying on AWS..."
ssh ec2-dev "~/apps/incident-tracker/stop.sh && ~/apps/incident-tracker/start.sh"

echo "✅ Deployment complete!"
