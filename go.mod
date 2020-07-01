module github.com/pivotal/scdf-k8s-prel

go 1.14

require (
	github.com/Sirupsen/logrus v1.4.0 // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/magiconair/properties v1.8.0
	github.com/moby/moby v1.13.1
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/pivotal/go-ape v0.0.0-20200224111603-3ada71e48e45
	github.com/pivotal/image-relocation v0.0.0-20200331082614-4cc801b4f8e7
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.4.0
)

// work around an import using the wrong case (see https://github.com/sirupsen/logrus/blob/master/README.md)
replace github.com/Sirupsen/logrus v1.4.0 => github.com/sirupsen/logrus v1.4.0
