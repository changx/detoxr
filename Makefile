x64:
	GOOS=linux GOARCH=amd64 buffalo build -v
all:
	-rm .env
	GOOS=linux GOARCH=amd64 buffalo build -v
	tar czvf release/detoxr-linux-amd64.tgz bin/detoxr
	GOOS=windows GOARCH=amd64 buffalo build -v
	tar czvf release/detoxr-windows-amd64.tgz bin/detoxr
	GOARCH=amd64 buffalo build -v
	tar czvf release/detoxr-macOS-amd64.tgz bin/detoxr
	GOARCH=mips GOOS=linux buffalo build -v
	tar czvf release/detoxr-linux-mips.tgz bin/detoxr
	buffalo build -v
	tar czvf release/detoxr-macOS-arm64.tgz bin/detoxr
	-ln -s .env.devel .env
