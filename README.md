# file-compare-move

Compare files in two directories and move duplicate files from destination directory to a new directory.

## Usage

```bash
Usage of file-compare-move:
  -d string
        destination directory to compare
  -o string
        output directory for duplicate files from dst that exist in src
  -s string
        source directory to compare
```

```
# Example:
./file-compare-move -s test/src -d test/dst -o test/duplicate
```

## Build

```bash
go build
```

