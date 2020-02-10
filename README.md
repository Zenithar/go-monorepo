# Golang Monorepository pattern

## Requirements

  * Golang 1.13
  * Docker to build images

## Mage (Local build)

> Recommended

You can invoke directly `magefile`from their respective directory. It will use
your environment to build artefacts.

```sh
> go run mage.go -d <directory>
```

* `tools`: install tools;
* `cmd/foo`: build the `bin/foo` artefact;

## Go-get

> Not recommended, build information will be missing from build.
> :warning: DON'T USE IT FOR PRODUCTION. :warning:

```sh
> go get -u -v github.com/Zenithar/go-monorepo/cmd/foo
> go install github.com/Zenithar/go-monorepo/cmd/foo
```

It will works with all commands.

## Targets

> All builds are made in a docker container to be the most reproductible
> to it could be.

Invoke `mage` using following targets :

```sh
> go run mage.go <target>
```

* API
  * `api:generate`: regenerates Protobuf files from descriptors;

* Docker images
  * `docker:foo` : build the foo docker container;

* Code maintenance (monorepo wide)
  * `code:format` : format all the code;
  * `code:lint` : lint all the code;
  * `code:licenser` : add license banner to all sources;
