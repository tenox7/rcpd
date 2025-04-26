NAME := rcpd

all: $(NAME)
$(NAME): *.go
	go build -o $(NAME) .

cross:
	GOOS=linux GOARCH=amd64 go build -a -o $(NAME)-amd64-linux .
	GOOS=linux GOARCH=arm go build -a -o $(NAME)-arm-linux .
	GOOS=linux GOARCH=arm64 go build -a -o $(NAME)-arm64-linux .
	GOOS=darwin GOARCH=amd64 go build -a -o $(NAME)-amd64-macos .
	GOOS=darwin GOARCH=arm64 go build -a -o $(NAME)-arm64-macos .
	GOOS=freebsd GOARCH=amd64 go build -a -o $(NAME)-amd64-freebsd .
	GOOS=freebsd GOARCH=arm64 go build -a -o $(NAME)-arm64-freebsd .
	GOOS=openbsd GOARCH=amd64 go build -a -o $(NAME)-amd64-openbsd .
	GOOS=netbsd GOARCH=amd64 go build -a -o $(NAME)-amd64-netbsd .
	GOOS=solaris GOARCH=amd64 go build -a -o $(NAME)-adm64-solaris .
	GOOS=aix GOARCH=ppc64 go build -a -o $(NAME)-ppc64-aix .
	GOOS=plan9 GOARCH=amd64 go build -a -o $(NAME)-amd64-plan9 .
	GOOS=plan9 GOARCH=arm go build -a -o $(NAME)-arm-plan9 .
	GOOS=windows GOARCH=amd64 go build -a -o $(NAME)-amd64-win64.exe
	GOOS=windows GOARCH=arm64 go build -a -o $(NAME)-arm64-win64.exe

docker-local:
	GOOS=linux GOARCH=amd64 go build -a -o $(NAME)-amd64-linux .
	GOOS=linux GOARCH=arm64 go build -a -o $(NAME)-arm64-linux .
	docker buildx build --platform linux/amd64,linux/arm64 -t tenox7/$(NAME):latest --load .

clean:
	rm -f $(NAME) $(NAME)-*
