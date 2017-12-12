# Make sure your setup is done..
# All you need to check is gopath and gobin exist. and once the go file is ready run
# go install -> this will install this app to your bin dir
# then cp that to /usr/local/bin
# Then enjoy
#!/usr/bin/env bash
go install
cp $GOBIN/MiniBlockChain /usr/local/bin
