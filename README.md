# abs-goodreads

A Goodreads Custom Metadata Provider for AudioBookShelf.

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
    --url "http://$ADDRESS:5555/search?query=The+Hobbit&author=J.R.R.+Tolkien"
```

## Setup

You can then set up abs-goodreads in AudioBookShelf.

```
Settings -> Item Metadata Utils -> Custom Metadata Providers -> Add
```

and entering the following details:

- Name: **GoodReads**
- URL: **\<your_address\>:5555**
- Authorization Header Value: **\<leave this unset\>**

See video below for a walkthrough:

https://github.com/ahobsonsayers/abs-goodreads/assets/32173585/54437af6-a17c-4458-bb82-479b183171da

## FAQ

### Why are covers not being returned?

The Goodreads API will sometimes not return cover images due to a variety of factors, usually related to copyright and the requirement to honour terms and conditions. Unfortunately, this is not something we can fix while using the Goodreads API - covers will need to be sourced from elsewhere.

See the FAQ below for a possible future alternative way to retrieve Goodreads metadata, which could solve this issue.

### How are you using the Goodreads API? Hasn't it been shut down?

Yes, since 8th December 2020, Goodreads no longer issues new API keys and plans to retire their API entirely in the future.

However, thanks to [LazyLibrarian](https://gitlab.com/LazyLibrarian/LazyLibrarian), we can use the read-only key they provide until Goodreads shuts down the API for good.

See the FAQ below for a possible future alternative way to retrieve Goodreads metadata, which could be used if/when Goodreads fully shuts down their API.

### Is there an alternative/better way to retrieve Goodreads metadata?

There are several issues with using the Goodreads API to retrieve book metadata.

To solve some of these issues, there is a plan at some point in the future to migrate to scraping metadata from Goodreads instead of using their API ([BiblioReads does this](https://github.com/nesaku/BiblioReads)). This is still in the works, but any help towards this goal is much appreciated.
