# eq
A CLI utility that detects duplicate files

---

## Usage

```
eq [paths...] [options]
```

You can provide any mix of files and directories.
Directories are scanned recursively. Symlinks and zero-size files are skipped.

---

## Description

`eq` scans files and directories, hashes their contents, and reports which files are exact duplicates.  
It supports multiple hashing algorithms like sha256, blake3 and xxh64.

---

## Options

| Flag | Description |
|------|-------------|
| `--hash string` | Hash algorithm to use. Supported: `sha256`, `blake3`, `xxh64`. Defaults to `blake3`. |
| `-h`, `--help` | Show help. |
| `-v`, `--version` | Print version information. |

---

## Examples

Scan the current directory:

```
eq
```

Detect duplicates in two directories:

```
eq ~/Pictures ~/Downloads
```

Use a specific hashing algorithm:

```
eq --hash sha256 /mnt/data
```

