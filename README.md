
## Assumptions

### Deduplication & "database" (`pkg/store`)
It's obious we need some sort of storage in order to check for duplicates and finally count valid, duplicateds and bad skus.
I decided to implement a simple in memory store. It keeps skus sorted to optimize searching for duplicates.

Also it offers an implementation of io.Reader, so we could make use of os.Copy or other methods to write directly into a file a large set of skus.

