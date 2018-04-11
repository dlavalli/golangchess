# Smshubbeat

Welcome to Smshubbeat.

Ensure that this folder is at the following location:
`${GOPATH}/github.com/tata/smshubbeat`

## Getting Started with Smshubbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Init Project
To get running with Smshubbeat and also install the
dependencies, run the following command:

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Smshubbeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/tata/smshubbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Smshubbeat run the command below. This will generate a binary
in the same directory with the name smshubbeat.

```
make
```


### Run

To run Smshubbeat with debugging output enabled, run:

```
./smshubbeat -c smshubbeat.yml -e -d "*"
```


### Test

To test Smshubbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/smshubbeat.template.json and etc/smshubbeat.asciidoc

```
make update
```


### Cleanup

To clean  Smshubbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Smshubbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/tata/smshubbeat
cd ${GOPATH}/github.com/tata/smshubbeat
git clone https://github.com/tata/smshubbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The hole process to finish can take several minutes.
