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

Build AtmanOS by running `bin/setup`,
which will download the required dependencies to the build directory
and then build AtmanOS itself.

## Example

Here's the output when running the [concurrent prime sieve][sieve] on a local
Xen host:

  [sieve]: https://golang.org/doc/play/sieve.go

```
(d87) Atman OS
(d87)      ptr_size:  8
(d87)    start_info:  0x502000
(d87)         magic:  xen-3.0-x86_64
(d87)      nr_pages:  8192
(d87)   shared_info:  393093120
(d87)    siff_flags:  0
(d87)     store_mfn:  68117
(d87)     store_evc:  1
(d87)   console_mfn:  68116
(d87)   console_evc:  2
(d87)       pt_base:  5263360
(d87)  nr_pt_frames:  7
(d87)      mfn_list:  5185536
(d87)     mod_start:  0
(d87)       mod_len:  0
(d87)      cmd_line:  [1024/1024]0x502080
(d87)     first_pfn:  0
(d87) nr_p2m_frames:  0
(d87)
(d87) 2
(d87) 3
(d87) 5
(d87) 7
(d87) 11
(d87) 13
(d87) 17
(d87) 19
(d87) 23
(d87) 29
```

Contributing
------------

We love pull requests from everyone.
By participating in this project,
you agree to abide by the [Go Community Code of Conduct][code of conduct].

[code of conduct]: https://golang.org/conduct

We expect everyone to follow the code of conduct
anywhere in the project's codebases,
issue trackers, chatrooms, and mailing lists.

License
-------

AtmanOS is Copyright (c) 2015 Bernerd Schaefer. It is free software,
and may be redistributed under the terms specified in the [LICENSE] file.

[LICENSE]: /LICENSE
