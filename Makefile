test:
	go test -v

cover:
	rm -rf *.coverprofile
	go test -v -coverprofile=routing.coverprofile
	gover
	go tool cover -html=routing.coverprofile