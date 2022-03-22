build:
	go build -o spam -v ./cmd/spam/main.go
	@echo "\033[0;32mBinary built successfully\033[m"
run:
	./spam
test:
	 go test -v -count=1 ./.../apiserver -cover
run_inmem:
	./urlShortener
run_psql:
	./urlShortener -mem psql
