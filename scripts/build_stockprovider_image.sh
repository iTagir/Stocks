#!/usr/bin/env bash

cd ../docker/stockprovider_image
if [ $? != 0 ];
then
  exit 1
fi


cp ../../stockprovider/stockprovider .

docker build -t itagir/stocks_provider .
