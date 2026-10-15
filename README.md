# Reusable Command Line Options

*for [spf13/cobra](https://github.com/spf13/cobra)-based CLI tools*

This package provides reusable command line options and configuration patterns for Carabiner CLI tools. It defines the `OptionsSet` interface for composable CLI flag management and provides common option sets that can be shared across multiple commands and applications.

## Available OptionsSets

### `log` - Logging Configuration

Package: `github.com/carabiner-dev/command/log`

Provides structured logging configuration with `--log-level` flag supporting debug, info, warn, and error levels. Includes `WithLogger(ctx)` method to initialize and inject logger into context.

**Usage:**

```go
import "github.com/carabiner-dev/command/log"

logOpts := &log.Options{}

rootCmd := &cobra.Command{
    Use: "myapp",
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        ctx, err := logOpts.WithLogger(cmd.Context())
        if err != nil {
            return err
        }
        cmd.SetContext(ctx)
        return nil
    },
}

logOpts.AddFlags(rootCmd)
```

### `keys` - Public Key File Management

Package: `github.com/carabiner-dev/command/keys`

Provides public key file path configuration with `--key/-k` flag. Validates key file existence and parses keys into `key.PublicKeyProvider` instances.

**Usage:**

```go
import "github.com/carabiner-dev/command/keys"

keyOpts := &keys.Options{}

verifyCmd := &cobra.Command{
    Use: "verify",
    RunE: func(cmd *cobra.Command, args []string) error {
        providers, err := keyOpts.ParseKeys()
        if err != nil {
            return err
        }
        // Use providers for verification...
        return nil
    },
}

keyOpts.AddFlags(verifyCmd)
```

### `output` - Output File Management

Package: `github.com/carabiner-dev/command/output`

Provides output file path configuration with `--output/-o` flag. Returns an `io.Writer` that writes to the specified file or defaults to STDOUT.

**Usage:**

```go
import "github.com/carabiner-dev/command/output"

outOpts := &output.Options{}

exportCmd := &cobra.Command{
    Use: "export",
    RunE: func(cmd *cobra.Command, args []string) error {
        writer, err := outOpts.GetWriter()
        if err != nil {
            return err
        }
        // Write output to file or STDOUT...
        fmt.Fprintln(writer, "data")
        return nil
    },
}

outOpts.AddFlags(exportCmd)
```

### `subject` - Attestation Subjects

Package: `github.com/carabiner-dev/command/subject`

Provides the `--subject/-s` flag to read attestation subjects from the command
line. Each value is either a digest in `algorithm:value` form or the path to a
file, which gets hashed. Every algorithm name from the
[in-toto DigestSet spec](https://github.com/in-toto/attestation/blob/main/spec/v1/digest_set.md)
is accepted (`sha256`, `sha1`, `sha512`, `sha3_256`, `gitCommit`, `gitTree`,
`gitBlob`, `gitTag`, `dirHash`, ...), matched case-insensitively and normalized
to the canonical spelling. Strings shaped like a digest are never looked up as
files, so a malformed digest is reported instead of being mistaken for a path.

`subject.ParseDigest` can also be used on its own to parse a single
`algorithm:value` string.

**Usage:**

```go
import "github.com/carabiner-dev/command/subject"

subjectOpts := &subject.Options{}

attestCmd := &cobra.Command{
    Use: "attest",
    PreRunE: func(cmd *cobra.Command, args []string) error {
        return subjectOpts.Validate()
    },
    RunE: func(cmd *cobra.Command, args []string) error {
        descriptors, err := subjectOpts.ResourceDescriptors()
        if err != nil {
            return err
        }
        // Use the in-toto resource descriptors as subjects...
        return nil
    },
}

subjectOpts.AddFlags(attestCmd)
```

## Contributing

This goal of this package is to to keep reusable sets of flags that have
a high reuse incidence. We try to keep the dependency tree as small as
possible. It is expected to mature slowly.

Feel free to file issues or suggestions if you find it useful.

This software is released under the Apache-2.0 license, paches are also welcome.
