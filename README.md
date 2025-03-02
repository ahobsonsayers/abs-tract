# abs-tract

This is an "all-in-one" book metadata provider for AudiobookShelf that can currently pull metadata from Goodreads and Kindle Store.

Current metadata providers plan to be improved, and other metadata providers are on the roadmap.

Please open issues for feature requests or new provider suggestions!

#### Why the name abs-tract?

I'm glad you asked - It's a fun play on words. AudiobookShelf is often abbreviated to ABS. Books can have abstracts. Plus, we are abstracting away the provision of metadata from AudiobookShelf. Therefore, abs-tract seemed suitably punny

## Providers

### Goodreads

#### Pros:

- A very good provider of general metadata for books.

#### Cons:

- Covers are quite low quality and often missing completely (see FAQ for reasoning)
- Goodreads has poor search functionality, making finding a matching book sometimes very difficult. Your mileage may vary. 90% of my library is able to be matched. See FAQ for tips on finding a match

#### Metadata Provided:

- Title
- Author
- Cover - Low quality. **Sometimes missing** (see FAQ)
- Original Publish Year
- Series Name
- Series Position
- Description - of "best" edition chosen by Goodreads
- Genres - Top 3 chosen by Goodreads users
- ISBN - of "best" edition chosen by Goodreads. **Sometimes missing**
- Publisher - of "best" edition chosen by Goodreads
- Language - of "best" edition chosen by Goodreads

### Kindle

#### Pros:

- Provides extremely high quality covers. Will make your eyes cry with happiness.

#### Cons:

- Currently, not much metadata is provided

#### Metadata Provided:

- Title
- Author
- Cover - Crazy high quality
- Publish Year - of edition chosen by Amazon. **Not original publish year**
- ASIN

## Running

The best way to run abs-tract is to use Docker. To run abs-tract using Docker, use the following command:

```bash
docker run -d \
    --name abs-tract \
    -p 5555:5555 \
    --restart unless-stopped \
    arranhs/abs-tract:latest
```

## Test

Test if abs-tract is working using curl.

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

## Setup with AudiobookShelf

You can then set up abs-tract in AudiobookShelf.

```
Settings -> Item Metadata Utils -> Custom Metadata Providers -> Add
```

and enter the following details:

### Goodreads

- Name: **GoodReads**
- URL: `http://<your_address>:5555/goodreads`
  - e.g. `192.168.1.100:5555/goodreads`
- Authorization Header Value: **Leave this unset**

### Kindle

- Name: **Kindle**
- URL: `http://<your_address>:5555/kindle/<your_region>`
  - e.g. `192.168.1.100:5555/kindle/uk`
- Authorization Header Value: **Leave this unset**

Region can be one of the following:

- au - Australia
- ca - Canada
- de - Germany
- es - Spain
- fr - France
- in - India
- it - Italy
- jp - Japan
- uk - United Kingdom
- us - United States

## FAQ

### Why is Goodreads not returning covers?

The Goodreads API sometimes does not return cover images due to a variety of factors, usually related to copyright and the requirement to honor terms and conditions. Unfortunately, this is not something we can fix while using the Goodreads API - covers will need to be sourced from elsewhere.

### How are you using the Goodreads API? Hasn't it been shut down?

Yes, since 8th, December 2020, Goodreads no longer issues new API keys and plans to retire their API entirely in the future.

However, thanks to [LazyLibrarian](https://gitlab.com/LazyLibrarian/LazyLibrarian), we can use the read-only key they provide until Goodreads shuts down the API for good.

See the FAQ below for a possible future alternative way to retrieve Goodreads metadata, which could be used if/when Goodreads fully shuts down their API.

### Is there an alternative/better way to retrieve Goodreads metadata?

There are several issues with using the Goodreads API to retrieve book metadata.

To solve some of these issues, there is a plan at some point in the future to migrate to scraping metadata from Goodreads instead of using their API ([BiblioReads does this](https://github.com/nesaku/BiblioReads)). This is still in the works, but any help towards this goal is much appreciated.
