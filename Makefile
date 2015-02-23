test:
	go test -v ./elbrewritessl

cover:
	go test -v ./elbrewritessl  -coverprofile=/tmp/coverage.out
	go tool cover -html=/tmp/coverage.out
