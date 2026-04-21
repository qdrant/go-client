#!/bin/bash

docker build -t qdrant-simple-example .

docker run --rm -it --network host qdrant-simple-example

