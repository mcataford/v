#!/bin/bash

# Scenario: User installs a specific version of python.

echo "Scenario: User installs a specific version of Python"

go build .

TARGET_VERSION="3.10.0"

V_ROOT=/tmp/v ./v init
V_ROOT=/tmp/v ./v python install $TARGET_VERSION --no-cache

INSTALLED_VERSIONS=$(V_ROOT=/tmp/v ./v python ls)

if [ -z "$(echo $INSTALLED_VERSIONS | grep $TARGET_VERSION)" ]; then
    echo "FAIL: Could not find target version."
    exit 1
else
    echo "SUCCESS"
fi
