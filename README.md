
# Installation & Usage

This service uses Docker Compose to manage its build and dependencies. To run the service, install Docker and run:

```bash
docker-compose up
```

You should be able to access the web service via http://localhost:8080.
Tests can be run via:

```
go test ./...
```


# Endpoints

`POST /create`

Endpoint used to create a shortened url

Parameters:
* `longurl`: URL to encode into a short URL. Must be a valid URL

Sample request:

```
curl -d "longUrl=http://www.example.com" -X POST http://localhost:8080/create
```

Sample response:

```
{
  "shortUrl": "UJuPI2jWg"  // not a fully qualified URL
}
```

`GET /{shortUrl}`

Redirects to the long URL, if it exists. Otherwise, returns 404

Sample request:

```
curl http://localhost:8080/UJuPI2jWg
```


`GET /{shortUrl}/stats`

Get the hit stats for a short URL

Sample request:

```
curl http://localhost:8080/UJuPI2jWg/stats
```

Sample response:

```
{
  "all_time": 24,
  "weekly": 12,
  "daily": 10,
  "minute": 0
}
```


# Architecture

This link shortener is designed for high speed and throughput.

## Creation

When creating a short URL, it selects the first entry from a pool of pre-generated candidate words. Access to the pool is synchronized to avoid any race conditions so that a single short URL does not get used multiple times.

## Persistence

When getting a short URL, we use an in-memory cache with a fall back to a durable KVS (Redis is used in this case). The size of a typical long URL might be around 500-1000 bytes, so even storing 10,000,000 URLs in memory would consume 5-10GB, so this approach should be fairly scalable. Redis is a great solution for our persistent storage as it is fast on its own even without a cache ([<= 3ms for most GET requests](https://redis.io/topics/benchmarks)), but allows us to horizontally scale rather well. We could shard based on lexigraphically ordering of short URLs (e.g. a separate redis node for short urls starting with `a-d`, `e-h`, etc.)

## Stats

Hit/view stats are tracked whenever a link is visited. Counts are stored in a per-minute, daily, weekly, and all-time basis. Redis hashes are used to store these values. When a link is visited, we bucket the current timestamp by each of these time ranges.

For example, a Unix timestamp of 1586312189 would get bucketed into:

* minute: 1586312160
* daily: 1586304000
* weekly: 1585785600
* all time: 0 (special case)

We then increment the counters at each of the keys associated with each of these values.

Currently, counters are not cleaned up which enables us to look at the hit counter at any time in the past. However, this would come at a high storage cost. In a real life scenario, we'd want to include a way to clean up historical stats.

## Pool

The keyword pool is a queue of minimum size that is continuously replenished when it falls below this minimum size. Currently this runs in a goroutine in the main process, but could be moved to a completely separate worker instance if desired. Using a pool allows short urls to be created faster as we don't have to worry about the time to generate a keyword, as we have to check to ensure the values we generate do not conflict with already assigned values. Additionally, using a pool allows us to handle a high throughput of creation were we to require that by pre-generating a larger pool size.


# Future

Some future improvements/potential enhancements include:

* Validating long urls (e.g. ensuring there are no redirect loops)
* Creating a dedicated background worker for pool generation
* Sharding on lexigraphic order of short URLs
* Pruning old counters to save storage space
* Using an LRU / memory cache of fixed size in front of redis instead of a simple struct
* Using a different algorithm or data structure for the pool as creating a short URL is currently an O(n) operation on the size of the pool queue, but this can be improved to just O(1)
