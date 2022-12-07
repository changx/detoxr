all:
	-rm .env
	GOOS=linux GOARCH=amd64 buffalo build -v
	-ln -s .env.devel .env
