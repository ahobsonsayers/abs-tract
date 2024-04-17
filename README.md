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

You can then set up abs-goodreads in AudioBookShelf by going to:

Settings -> Item Metadata Utils -> Custom Metadata Providers -> Add

and entering the following details

- Name: **Goodreads**
- URL: **\<your_address\>:5555**
- Authorization Header Value: **\<leave this unset\>**

See video below for walkthrough
