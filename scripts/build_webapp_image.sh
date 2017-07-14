#!/usr/bin/env bash

kubectl delete deployment stockswebapp

docker rmi itagir/stocks_wa

cd ../docker/webapp_image
if [ $? != 0 ];
then
  exit 1
fi


cp ../../web/webapp .
cp -R ../../web/docs .

docker build -t itagir/stocks_wa .
docker push itagir/stocks_wa
cd ../../k8s
kubectl create -f webapp_deployment.yaml
kubectl get pods
