
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

