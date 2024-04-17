# abs-goodreads

A Goodreads Custom Metadata Provider for AudioBookShelf

## Run

The best way to run abs-goodreads is to use docker.

```bash
docker run -d \
    --name abs-goodreads \
    -p 5555:5555 \
    --restart unless-stopped \
    arranhs/abs-goodreads:latest
```

## Test

Test abs-goodreads is working using curl.

```bash
ADDRESS=localhost
curl --request GET \
    --url "http://$ADDRESS:5555/search?query=The%20Hobbit&author=J.R.R.%20Tolkien"
```

## Setup

You can then set up abs-goodreads in AudioBookShelf.

```
Settings -> Item Metadata Utils -> Custom Metadata Providers -> Add
```

and entering the following details:

- Name: **Goodreads**
- URL: **\<your_address\>:5555**
- Authorization Header Value: **\<leave this unset\>**

See video below for a walkthrough:

https://github.com/ahobsonsayers/abs-goodreads/assets/32173585/54437af6-a17c-4458-bb82-479b183171da
