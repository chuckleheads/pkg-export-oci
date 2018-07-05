# Habitat OCI Exporter

The Habitat OCI exporter wraps the Habitat client to install a package to a ROOT_FS, then creates a RunC config for ease of use.

Artisanally crafted, shade grown, bespoke container solution for RunC

## Usage

For single run services
```
sudo pkg-export-oci chef/inspec --entrypoint="inspec --flags"
```

For long running services
```
sudo pkg-export-oci core/redis
```

This will give you a bzip tarball that contains everyting you need to run your service.

## Running

```
tar xvjf chef-inspec.tar.bz2
sudo runc run inspec
```
