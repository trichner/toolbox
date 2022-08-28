#!/bin/bash

cat matter_* | jq -r '.Id' | xargs -L 1 kraki describe-matter --matter-id | jq -r '(.Id + " " + .Exports[].Status)'
