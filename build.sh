#!/bin/sh

APP="gblog"

packr2 
gox -tags="embed" -osarch="linux/amd64" -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}"
gox -tags="embed" -osarch="windows/amd64" -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}"
     
