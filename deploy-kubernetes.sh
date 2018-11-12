#!/usr/bin/env bash
GC_PROJECT=$1
ZONE=$2
CLUSTER_NAME=$3
VERSION=$4
SECRET_FILE=$5

echo "Deploying to Kubernetes"
CLOUDSDK_CORE_DISABLE_PROMPTS=1
if [ ! -d ${HOME}/google-cloud-sdk ]; then
    curl https://sdk.cloud.google.com | bash;
fi

echo "Setting up Google Cloud SDK"
gcloud auth activate-service-account --key-file $SECRET_FILE
gcloud config set project $GC_PROJECT
gcloud config set compute/zone $ZONE
gcloud container clusters get-credentials $CLUSTER_NAME

echo "Deploying..."
kubectl set image deployment backend backend=nginx:$VERSION
