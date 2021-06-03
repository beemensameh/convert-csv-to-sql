Command = go
Main = .

build:
	${Command} build ${Main}

run:
	${Command} run ${Main}

dev:
	reflex -r '\.go' -R '\_test.go' -s -- sh -c "${Command} run ${Main}"

start:build
	./main

tidy:
	${Command} mod tidy
