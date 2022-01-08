setup-tools:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/vektra/mockery/v2@v2.5.1
	go install github.com/uudashr/gocognit/cmd/gocognit@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/mcubik/goverreport@latest
