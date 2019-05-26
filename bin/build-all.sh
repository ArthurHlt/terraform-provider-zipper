#!/bin/bash

set -e

OUTDIR=$(dirname $0)/../out
BINARYNAME="terraform-provider-zipper"

GOARCH=amd64 GOOS=windows $(dirname $0)/build && cp $OUTDIR/$BINARYNAME "$OUTDIR/${BINARYNAME}_windows_amd64.exe"
GOARCH=amd64 GOOS=linux $(dirname $0)/build && cp $OUTDIR/$BINARYNAME "$OUTDIR/${BINARYNAME}_linux_amd64"
GOARCH=amd64 GOOS=darwin $(dirname $0)/build && cp $OUTDIR/$BINARYNAME "$OUTDIR/${BINARYNAME}_darwin_amd64"