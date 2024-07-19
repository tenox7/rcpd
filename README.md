# RCPD

This is a modern re-implementation of [rcp](https://linux.die.net/man/1/rcp) (remote copy protocol) daemon, originally part [berkeley r-commands](https://en.wikipedia.org/wiki/Berkeley_r-commands).

## Motivation

Easily copy files between vintage computer clients and a modern server/NAS. The r-commands are ubiquitous on old operating systems, even Windows NT. Also rcp is NAT/QEMU/Docker frienly as it uses only a single TCP port, unlike for example `ftp` or `tftp`.

In recent years both r-daemons and [inetd](https://en.wikipedia.org/wiki/Inetd) have been deprecated on modern OSes, leaving little alternatives.

## Implementation

This implementaion of rcpd is fully self contained, stand-alone, statically linked binary, with zero dependencies. Doesn't require `rshd`, `inetd`, `shell` or anything else. It can be run as a Docker container with ease.

## Security considerations

This implementation is fully open, with no security what so ever. Ignoring user names, authentication, `.rhosts`, `/etc/hosts.equiv` and all that nonsense. It's intended to be used on a secure LAN only.

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
