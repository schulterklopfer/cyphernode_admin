#!/usr/bin/env bash

mkdir -p build_context/Users/jash/go/src/github.com/SatoshiPortal/

rsync -avr --exclude '.git' /Users/jash/go/src/github.com/SatoshiPortal/cam build_context/Users/jash/go/src/github.com/SatoshiPortal/ &&\
  docker-compose build &&\
  rm -rf build_context/Users/jash/go/src/github.com/SatoshiPortal/cam
