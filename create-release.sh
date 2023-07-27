#!/bin/bash

go build -o deployer
tar -czvf latest-linux.tar.gz deployer contract-manifest-example.json setup-output-example.json config-example.json
