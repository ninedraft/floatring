#!/usr/bin/env bash

tag=$(echo $1 | awk -F"/" '{print $NF}')
curl https://proxy.golang.org/github.com/ninedraft/floatring/@v/$tag.info 