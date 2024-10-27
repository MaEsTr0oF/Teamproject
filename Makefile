pkgs = $(shell go list ./...)

lint:
	# go install golang.org/x/lint/golint@latest
	# go install honnef.co/go/tools/cmd/staticcheck@latest
	# go install github.com/kisielk/errcheck

	golint $(pkgs)
	go vet $(pkgs)
	staticcheck $(pkgs)
	errcheck $(pkgs)