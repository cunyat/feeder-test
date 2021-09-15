# Feeder

## How to run

In order to run the app you will need `go1.16+` and optionally `make`

In `Makefile` you will find command to run, test and lint the application.

- `make run`: to execute the server
- `make test`: to execute tests
- `make lint`: to execute linter

The server, also accepts the output file as a parameter, to indicate it you should either build the application and execute the binary or use go run command:
```bash
make build
feeder -output out/skus-for-today.txt
-- or --
go run cmd/feeder/main.go -output file-to-store-skus.txt
```


I also provide a simple script that sends random skus to the server (`make run-client`). It accepts two parameters for customization but should be run directly with `go run`:
- `-count n`: indicate the number of skus to send
- `-terminate`: will send a terminate sequence when all skus are sent

For example (server must be up to work):
```bash
go run cmd/client/main.go -count 123 -terminate
```

## Direcory structure

```text
feeder
├── cmd/ -> application entrypoints
├── internal/feeder/ -> application specific code
├── out/ -> file for app output files
└── pkg/ -> application shareable code
    ├── server/ -> server that listens for incoming connections
    ├── store/ -> store that deduplicates strings 
    └── utils/ -> some utility code
```

I tried to follow [Go standard layout](https://github.com/golang-standards/project-layout):

Application entrypoints are placed in `cmd/` directory, and one folder for each entrypoint.

In `pkg/` folder, i placed code that could be shared among different internal service or maybe promoted to a library. It doesn't contain any business logic and is reusable.

in `internal/` folder, there is the "business logic" of the application, this code would not be shared between services, for that it keeps in internal "private" folder.


## Assumptions

### Client messages

As I understand, a client will only send ONE sku or terminate sequence, at first moment I supposed that a client could send many skus,
but reading many times I understand it's not.

### SKUs format

Test doesn't specify if sku is case-sensitive, or if all characters must be uppercase (like in examples). 
I decided to allow uppercase and lowercase, but match duplicated in case-insensitive way.
For this, I transform all skus to uppercase when they are stored

### Deduplication & "database" (`pkg/store`)
It's obious we need some sort of storage in order to check for duplicates and finally count valid, duplicateds and bad skus.
I decided to implement a simple in memory store. It keeps skus sorted to optimize searching for duplicates.

Also, it offers an implementation of io.Reader, so we could make use of os.Copy or other methods to write directly into a file a large set of skus.

