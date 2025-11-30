# eq
A fast CLI utility that detects duplicate files

---

## Usage

```
eq [paths...] [options]
```

You can provide any mix of files and directories.
Directories are scanned recursively. Symlinks are skipped.

---

## Description

`eq` scans files and directories, hashes their contents, and reports which files are exact duplicates.  
It supports multiple hashing algorithms and can optionally verify matches using a byte-by-byte comparison.


---

## Options

| Flag | Description |
|------|-------------|
| `--hash string` | Hash algorithm to use. Supported: `sha256`, `blake3`, `xxh64`. Defaults to `blake3`. |
| `-h`, `--help` | Show help. |
| `-v`, `--version` | Print version information. |

---

## Examples

Detect duplicates in two directories:

```
eq ~/Pictures ~/Downloads
```

Use a specific hashing algorithm:

```
eq --hash sha256 /mnt/data
```

Scan the current directory:

```
eq
```

