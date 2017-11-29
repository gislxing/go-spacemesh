## go-unruly
UnrulyOS p2p node

### Build

Compile the .proto files using the protobufs go compiler:

```
cd pb
protoc --go_out=. ./*.proto
```
#### Vendoring 3rd party GO packages
We use [govendor](https://github.com/kardianos/govendor) for all 3rd party packages.
We commit to git all 3rd party packages in the vendor folder so we have our own copy of versioned releases.
To update a 3rd party package use vendor.json and govendor commands.

Installing govendor:
```
go get -u github.com/kardianos/govendor
```

To get the vendor packages use:
```
govendor init
govendor sync
```

To build the node use:

```
go build
```

### Running

```
./go-unruly
```

### Contributing

- go-unruly is part of [The UnrulyOS open source Project](https://unrulyos.io), and is MIT licensed open source software.
- We welcome contributions big and small! 
- We welcome major contributors to the unruly core dev team.
- Please make sure to scan the [issues](https://github.com/UnrulyOS/go-unruly/issues). 
- Search the closed ones before reporting things, and help us with the open ones.

Guidelines:

- Read the UnrulyOS project white paper
- Please make branches + pull-request, even if working on the main repository
- Ask questions or talk about things in [Issues](https://github.com/UnrulyOS/go-unruly/issues) or #unruly on freenode.
- Ensure you are able to contribute (no legal issues please)
- Run `go fmt` before pushing any code
- Run `golint` and `go vet` too -- some things (like protobuf files) are expected to fail.
- Get in touch with @avive about how best to contribute
- Have fun hacking away our blockchain future!

There's a few things you can do right now to help out:
 - **check out existing issues**. This would be especially useful for modules in active development.
 - **Perform code reviews**.
 - **Add tests**. There can never be enough tests.


### Tests

### tasks

- Get rid of libp2p GX deps (asap) and them to vendor folder
- Support command line args in a robust way 
- Support basic account and keys ops

