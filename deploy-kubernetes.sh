#!/usr/bin/env bash
GC_PROJECT=$1
ZONE=$2
CLUSTER_NAME=$3
VERSION=$4
SECRET_FILE=$5

PATH=$PATH:${HOME}/google-cloud-sdk/bin
CLOUDSDK_CORE_DISABLE_PROMPTS=1

echo "Installing Google Cloud SDK..."
curl https://sdk.cloud.google.com | bash;
gcloud components update
gcloud components install kubectl

gcloud auth activate-service-account --key-file $SECRET_FILE
gcloud config set project $GC_PROJECT
gcloud config set compute/zone $ZONE
gcloud container clusters get-credentials $CLUSTER_NAME

echo "Deploying..."
kubectl set image deployment backend backend=nginx:$VERSION
