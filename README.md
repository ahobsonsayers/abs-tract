# abs-goodreads

A Goodreads Custom Metadata Provider for AudioBookShelf

## Usage

The best way to run is to use docker

```bash
docker run -d \
    --name abs-goodreads \
    -p 5555:5555 \
    --restart unless-stopped \
    arranhs/abs-goodreads:latest
```
