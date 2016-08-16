DEPS=go list -f '{{range .TestImports}}{{.}} {{end}}' ./...

update-deps:
	rm -rf Godeps
	rm -rf vendor
	go get github.com/tools/godep
	godep save ./...

install-deps:
	go get github.com/tools/godep
	godep restore
	$(DEPS) | xargs -n1 go get -d

fmt:
	bash -c 'go list ./... | grep -v vendor | xargs -n1 go fmt'

test-oauth:
	bash -c 'go test -timeout=30s github.com/RichardKnop/example-api/oauth'

test-accounts:
	bash -c 'go test -timeout=30s github.com/RichardKnop/example-api/accounts'

test-facebook:
	bash -c 'go test -timeout=60s github.com/RichardKnop/example-api/facebook'

test:
	bash -c 'go list ./... | grep -v vendor | xargs -n1 go test -timeout=60s'
