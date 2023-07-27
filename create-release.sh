#!/bin/bash

go build -o deployer
tar -czvf latest-linux.tar.gz deployer examples
