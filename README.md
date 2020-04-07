# architecture

* 500 bytes x 10,000,000 = 5gb
* synchronization is handled

pros of pool:
* handle high throughput of creation
* can do more in background (validate, etc)

# Dependencies

* Use go mod to manage dependencies
* Gorilla Mux for building the web service

# installation

```bash
go get -u github.com/golang/dep/cmd/dep
dep ensure
```

# future

* validate long urls (redirect loops, etc)
* dedicated background workers for pool generation
* sharding on lexigraphic order of shorturl
