# abs-tract

This is an "all-in-one" book metadata provider for AudioBookShelf that can currently pull metadata from Goodreads and Kindle Store.

Current metadata providers plan to be improved and other metadata providers are on the roadmap. 

Please open issues for feature requests or new provider suggestions!
#### Why the name abs-tract?

I'm glad you asked - It's a fun play on words. AudioBookShelf is often abbreviated to ABS.  Books can have abstracts. Plus we are abstracting away the providing of metadata from AudioBookShelf. Therefore abs-tract seemed suitably punny 💩
## Providers

### Goodreads

#### Pros:
- A very good provider of general metadata for books.
#### Cons:
- Covers are quite low quality and often missing completely (see FAQ for reasoning)
- Goodreads has poor search functionality making finding a matching book sometimes very difficult. Your mileage may vary. 90% of my library is able to be matched. See FAQ for tips on finding a match
#### Metadata Provided:
- Title
- Author
- Cover - Low quality. **Sometimes missing** (see FAQ)
- Original Publish Year
- Series Name
- Series Position
- Description - of "best" edition chosen by goodreads
- Genres - Top 3 chosen by goodreads users
- ISBN - of "best" edition chosen by goodreads. **Sometimes missing**
- Publisher - of "best" edition chosen by Goodreadns
- Language - of "best" edition chosen by Goodreadns

### Kindle

#### Pros:
- Provides extremely high quality covers. Will make your eyes cry with happiness.
#### Cons:
- Currently not much metadata is provided
#### Metadata Provided:
- Title
- Author
- Cover - Crazy high quality
- Publish Year - of edition chosen by Amazon. **Not original publish year**
- ASIN
## Running

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

### Goodreads

```bash
ADDRESS=localhost
curl --request GET \
    --url "http://$ADDRESS:5555/goodreads/search?query=The+Hobbit&author=J.R.R.+Tolkien"
```
### Kindle

```bash
ADDRESS=localhost
curl --request GET \
    --url "http://$ADDRESS:5555/kindle/uk/search?query=The+Hobbit&author=J.R.R.+Tolkien"
```
## Setup with AudioBookShelf

You can then set up abs-goodreads in AudioBookShelf.

```
Settings -> Item Metadata Utils -> Custom Metadata Providers -> Add
```

and entering the following details:

### Goodreads

- Name: **GoodReads**
- URL: **http://\<your_address\>:5555/goodreads**
- Authorization Header Value: **\<leave this unset\>**

### Kindle

- Name: **Kindle**
- URL: **http://\<your_address\>:5555/kindle/\<your_region\>**
- Authorization Header Value: **\<leave this unset\>

Region can be one of the following:

- au - Australia 
- ca - Canada 
- de - Germany 
- es - Spain
- fr - Franch
- in - Indian 
- it - Italy 
- jp - Japan
- uk - United Kingdom 
- us - United States

### Setup video walkthrough:

https://github.com/ahobsonsayers/abs-goodreads/assets/32173585/54437af6-a17c-4458-bb82-479b183171da

## FAQ

### Why is Goodreads not returning some covers?

The Goodreads API will sometimes not return cover images due to a variety of factors, usually related to copyright and the requirement to honour terms and conditions. Unfortunately, this is not something we can fix while using the Goodreads API - covers will need to be sourced from elsewhere.

See the FAQ below for a possible future alternative way to retrieve Goodreads metadata, which could solve this issue.

### How are you using the Goodreads API? Hasn't it been shut down?

Yes, since 8th December 2020, Goodreads no longer issues new API keys and plans to retire their API entirely in the future.

However, thanks to [LazyLibrarian](https://gitlab.com/LazyLibrarian/LazyLibrarian), we can use the read-only key they provide until Goodreads shuts down the API for good.

See the FAQ below for a possible future alternative way to retrieve Goodreads metadata, which could be used if/when Goodreads fully shuts down their API.

### Is there an alternative/better way to retrieve Goodreads metadata?

There are several issues with using the Goodreads API to retrieve book metadata.

To solve some of these issues, there is a plan at some point in the future to migrate to scraping metadata from Goodreads instead of using their API ([BiblioReads does this](https://github.com/nesaku/BiblioReads)). This is still in the works, but any help towards this goal is much appreciated.
