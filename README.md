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
