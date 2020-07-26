prepare:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.24.0
	go get golang.org/x/tools/cmd/goimports

check:
	goimports -w .
	golangci-lint run --no-config --issues-exit-code=0 --deadline=30m \
  --disable-all --enable=deadcode  --enable=gocyclo --enable=golint --enable=varcheck \
  --enable=structcheck --enable=errcheck --enable=ineffassign \
  --enable=unconvert --enable=goconst --enable=gosec --enable=megacheck --enable=maligned --enable=dupl --enable=interfacer \
  --skip-files ".*_test.go"
