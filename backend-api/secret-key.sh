#!/bin/bash

KEY=$(head -c 32 /dev/urandom | base64)

echo $KEY > secret.txt
