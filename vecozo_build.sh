#!/bin/bash

cp go.mod go.mod.bak
cat go.mod.vecozo > go.mod
go build
mv go.mod.bak go.mod
