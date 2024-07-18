all: rcpd

rcpd: *.go
	go build .

cross:
	GOOS=linux GOARCH=amd64 go build -a -o rcpd-amd64-linux .
	GOOS=linux GOARCH=arm go build -a -o rcpd-arm-linux .
	GOOS=linux GOARCH=arm64 go build -a -o rcpd-arm64-linux .
	GOOS=darwin GOARCH=amd64 go build -a -o rcpd-amd64-macos .
	GOOS=darwin GOARCH=arm64 go build -a -o rcpd-arm64-macos .
	GOOS=freebsd GOARCH=amd64 go build -a -o rcpd-amd64-freebsd .
	GOOS=openbsd GOARCH=amd64 go build -a -o rcpd-amd64-openbsd .
	GOOS=netbsd GOARCH=amd64 go build -a -o rcpd-amd64-netbsd .
	GOOS=solaris GOARCH=amd64 go build -a -o rcpd-adm64-solaris .
	GOOS=aix GOARCH=ppc64 go build -a -o rcpd-ppc64-aix .
	GOOS=plan9 GOARCH=amd64 go build -a -o rcpd-amd64-plan9 .
	GOOS=plan9 GOARCH=arm go build -a -o rcpd-arm-plan9 .
	GOOS=windows GOARCH=amd64 go build -a -o rcpd-amd64-win64.exe
	GOOS=windows GOARCH=arm64 go build -a -o rcpd-arm64-win64.exe

clean:
	rm rcpd rcpd-*
