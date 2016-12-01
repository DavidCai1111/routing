test:
	go test -v

cover:
	rm -rf *.coverprofile
	go test -coverprofile=routing.coverprofile
	gover
	go tool cover -html=routing.coverprofile