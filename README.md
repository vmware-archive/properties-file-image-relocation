# scdf-k8s-prel
Image relocation for Java properties files.

# Usage

The `prel` utility takes a properties file along with the images it refers to and packages them up in a zipped archive.
The archive may be passed to a system which is disconnected from the internet and then `prel` may be used to relocate
the images to a private registry and generate a modified properties file in which the images have been replaced with their
relocated counterparts.

_Relocating_ an image means copying it to another repository, possibly in a private registry. `prel` preserves digests when
packaging and relocating images, so you can be sure that each relocated image is identical to the original.
For more information, refer to the image relocation repository [README](https://github.com/pivotal/image-relocation#what-is-image-relocation).

## Installing and upgrading

To install `prel` either download and unpack a suitable binary release for your platform from the
[releases](https://github.com/pivotal/scdf-k8s-prel/releases) page or, if you have Go and make installed, clonse this repository,
change directory into it, and then issue the following:
```
make prel
```
Add the binary to your path if necessary.

You can check which version you have installed by issuing:
```
prel --version
```

## Properties file format

A Java properties file is a collection of named properties with a _de facto_ standard [format](https://en.wikipedia.org/wiki/.properties).

By default, `prel` interprets any property value prefixed with `docker:` or `docker://` as denoting an image reference. The image reference
immediately follows the prefix and extends up to, but does not include, the next whitespace character or end of file, whichever occurs first.
So, for example, the following property:
```
sink.cassandra=docker:springcloudstream/cassandra-sink-rabbit:2.1.2.RELEASE
```
declares a property named `sink.cassandra` with value denoting an image with reference `springcloudstream/cassandra-sink-rabbit:2.1.2.RELEASE`.

## Configuring

If you need to use other property value prefixes, you can create a [TOML](https://toml.io) configuration file named `.prel.toml` in your home directory.
This file contains a list of prefixes. If `docker:` or `docker://` are required, these must be listed in the configuration file.
For example, the following configuration will cause `prel` to recognize the prefixes `docker:` and `image:` (but _not_ `docker://`):
```
property_value_prefixes = [ "docker:", "image:" ]
```

## Packaging up a properties file and its images

The `prel package` command takes a properties file and packages it up, along with any images it refers to, in a `.tgz` zipped archive.
The images are stored in a standard [OCI image layout](https://github.com/opencontainers/image-spec/blob/master/image-layout.md) inside the archive.

The properties file may be specified using a file system path or as a HTTP/HTTPS URL. For example, the following command takes a
properties file at `https://dataflow.spring.io/rabbitmq-docker-latest` and packages it up into a zipped archive named `rabbitmq-docker.tgz`:
```
prel package https://dataflow.spring.io/rabbitmq-docker-latest --archive rabbitmq-docker.tgz
```

For detailed help on the command, issue:
```
prel package --help
```

## Relocating a packaged properties file and its images

The `prel relocate` command takes a `.tgz` zipped archive produced by `prel package`, relocates its images by pushing them to a
specified registry, and creates a relocated version of the original properties file in which the image references have been replaced by their
relocated counterparts.

`prel relocate` takes a parameter of the path to the zipped archive and a flag _repository prefix_ whose value is used to prefix the
relocated image references. For example, the following command relocates the zipped archive `rabbitmq-docker.tgz` with a repository prefix of
`example.com/user` and creates the relocated properties file in `rabbitmq-docker.properties`:
```
prel relocate rabbitmq-docker.tgz --repository-prefix example.com/user --output rabbitmq-docker.properties
```

As an example of how `prel relocate` maps image references, the reference `springcloudstream/cassandra-sink-rabbit:2.1.2.RELEASE` above
would be mapped by this command to:
```
example.com/user/springcloudstream-cassandra-sink-rabbit-bec30c995c2e67ec4a914d3acce0ef57:2.1.2.RELEASE
```

For more information on how image names are mapped, see [Relocating image names](https://github.com/pivotal/image-relocation#relocating-image-names).

For detailed help on the command, issue:
```
prel relocate --help
```

## Authentication and access control

Authentication of `prel` to registries is provided by a regular docker configuration file. For details, see the
[authn README](https://github.com/google/go-containerregistry/blob/master/pkg/authn/README.md) of the Go Container Registry repository (a transitive
dependency of `prel`).

## Docker daemon

`prel` accesses image registries directly and not via the Docker daemon. This is primarily because the daemon doesn't guarantee to provide the 
same digest of an image as when the image has been pushed to a registry. Consequently, you don't need to have a Docker daemon running in order to use `prel`.

# Development

To run the tests, [install Go](https://golang.org/doc/install) and `make` and then issue the following from the root directory of this repository:
```
make test
```

Check linting (so you don't get caught out by CI), after installing [golangci-lint](https://golangci-lint.run/):
```
make lint
```

To create a release on github, merge a commit which removes "-snapshot" from [VERSION](VERSION) (and, optionally,
bumps the major version if there has been an incompatible change), then push a tag beginning with "v".
To continue development of the next release, merge a commit which bumps the minor version in [VERSION](VERSION) and adds
"-snapshot" back in.