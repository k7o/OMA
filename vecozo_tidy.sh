#!/bin/bash

cp go.mod go.mod.bak
cat go.mod.vecozo > go.mod
go mod tidy
mv go.mod.bak go.mod
