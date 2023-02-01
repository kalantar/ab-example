---
template: main.html
title: "Iter8 Version"
hide:
- toc
---
## iter8 version

Print Iter8 CLI version

### Synopsis


Print the version of Iter8 CLI.

	iter8 version

The output may look as follows:

	$ cmd.BuildInfo{Version:"v0.13.0", GitCommit:"f24e86f3d3eceb02eabbba54b40af2c940f55ad5", GoVersion:"go1.19.3"}

In the sample output shown above:

- Version is the semantic version of the Iter8 CLI.
- GitCommit is the SHA hash for the commit that this version was built from.
- GoVersion is the version of Go that was used to compile Iter8 CLI.


```
iter8 version [flags]
```

### Options

```
  -h, --help    help for version
      --short   print abbreviated version info
```

### Options inherited from parent commands

```
  -l, --loglevel string   trace, debug, info, warning, error, fatal, panic (default "info")
```

### SEE ALSO

* [iter8](iter8.md)	 - Kubernetes release optimizer

###### Auto generated by spf13/cobra on 25-Jan-2023