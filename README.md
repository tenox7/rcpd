# RCPD

This is a modern re-implementation of rcp (remote copy protocol) daemon, originally part [berkeley r-commands](https://en.wikipedia.org/wiki/Berkeley_r-commands)

## Motivation

Copy files between vintage computer clients and modern server/NAS. R-commands are ubiquitous on old operating systems, even on Windows NT. However, in recent years both r-daemons and [inetd](https://en.wikipedia.org/wiki/Inetd) have
been deprecated, leaving no alternatives, maybe except [tftp](https://en.wikipedia.org/wiki/Trivial_File_Transfer_Protocol).

## Usage

### Server

./rcpd -root_dir=/some/path

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
