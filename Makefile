http:
	go run ./cmd/http/main.go

simulation:
	go run ./cmd/simulation/main.go

automaton:
	go run ./cmd/automaton/main.go

test:
	go test -v -race -covermode=atomic ./...

lint: misspell staticcheck vet fmt optimise
# The disabled linters are deprecated
	golangci-lint run --enable-all --sort-results --tests --fix \
		--disable maligned --disable interfacer --disable scopelint --disable golint --disable exhaustivestruct \
		--disable varcheck --disable ifshort --disable nosnakecase --disable structcheck --disable deadcode \
		--disable forbidigo --disable depguard --disable ireturn --disable goerr113 ./...

misspell:
	misspell -locale UK .

staticcheck:
	staticcheck ./...

vet:
	go vet -all ./...

fmt:
	go fmt ./cmd/*
	go fmt ./internal/*
	go fmt ./pkg/*
	gofumpt -l -w .

optimise: structs

structs:
	structslop -fix -apply ./...

install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/client9/misspell/cmd/misspell@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/orijtech/structslop/cmd/structslop@latest