# brotli

If you want to use `brotli` compress, you should add build flag `-tags brotli` for go build.


## build library

### alpine

```bash
$ apk update
$ apk add git make g++ bash cmake
$ git clone --depth=1 https://github.com/google/brotli.git
$ cd brotli && mkdir out && cd out
$ ../configure-cmake --disable-debug
$ make
$ make test
$ make install
```
