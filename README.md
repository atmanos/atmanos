> *NOTE*: AtmanOS is highly experimental, and not particularly featureful; for
> example, it does not yet have console, network, or storage drivers.

# AtmanOS

AtmanOS allows you to compile ordinary Go programs into standalone unikernels
that run under the Xen hypervisor.

AtmanOS is implemented as a series of patches and additional files for Go's
runtime and standard library. AtmanOS is a `GOOS` available from the modified
`go` command, meaning programs can be cross-compiled for AtmanOS in the normal
manner:

```
GOOS=atman go build
```

## Building AtmanOS

The parent directory of AtmanOS should look like:

```
$ ls ../
atmanos/
go/
go_bootstrap/
```

Where:

  * `go/` should be a clone of [Go](https://github.com/golang/go). The `go1.5.2`
    tag should be checked out.
  * `go_bootstrap/` is a precompiled version of [Go 1.4.3](https://golang.org/dl/#go1.4.3)

Then run `make patch` and `make build` to build.
