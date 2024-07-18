# RCPD

This is a modern re-implementation of [rcp](https://linux.die.net/man/1/rcp) (remote copy protocol) daemon, originally part [berkeley r-commands](https://en.wikipedia.org/wiki/Berkeley_r-commands)

## Motivation

Copy files between vintage computer clients and modern server/NAS. R-commands are ubiquitous on old operating systems, even on Windows NT. However, in recent years both r-daemons and [inetd](https://en.wikipedia.org/wiki/Inetd) have
been deprecated, leaving no alternatives, maybe except [tftp](https://en.wikipedia.org/wiki/Trivial_File_Transfer_Protocol).

## Usage

### Server

```sh
./rcpd -root_dir=/some/path
```

The server must bind to port `514/tcp`, which may require elevated privileges.

### Docker

https://hub.docker.com/r/tenox7/rcpd

Inside docker container root dir is `/srv`:

```sh
docker run -d --name rcpd -v /some/dir:/srv -p 514:514 tenox7/rcpd:latest
```

### Client

Just like `scp`:

```sh
rcp host:/dir/file.txt .
rcp file.txt host:
rcp file.txt host:/path
```

## Legal

- This code has been writen entirely by Claude
- License: Public Domain
