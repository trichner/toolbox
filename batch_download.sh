#!/bin/bash

cat matter_* | jq -r '.Id' | xargs -L 1 kraki download-export --matter-id
