# Protoc Serilization Plugin

## Overview

This library aims to help with language agnostic data encoding.

## Features

### Go

to help with serilization and deserialization operations , It parses the protobuf generated code and adds `xml` and `mapstructure` tags

- **JSON** : Generates Message specific JSON serilizer and deserilizer.

- **XML** :  Generates Message specific XML serilizer and deserilizer.

- **Mapstructure** : it generates serilizer and deserilizer to/from `map[string]interface{}` based on `github.com/mitchellh/mapstructure` library.

## Directory Structure

### example

This directory contains the example proto file

### modules

This directory is contains protobuf compiler modules based on different target language implementation .
there is also a `proto` directory in which shared data structure models used by different compiler modules.
These data structures fields are populated by the modules for each language and then passed to the approportiate template.  

### templates

this directory contains language specific templates that the plugin uses to generate the code.

## Roadmap

- [x] Generating Go code from proto files
  - [x] JSON
  - [x] XML
  - [x] mapstructure
- [ ] Generating TypeScript code from proto files
- [ ] Generating vanilla Javascript code from proto files
- [ ] Generating Rust code from proto files
- [ ] Generating Java code from proto files
- [ ] Generating C++ code from proto files

## Build

### Dependencies

We use [`mage`](https://github.com/magefile/mage) as our main build tool and [`packr`](https://github.com/gobuffalo/packr/) to embed static files (the templates).

install `mage` by running  :

```shell
go get -u -d -v github.com/magefile/mage
cd $GOPATH/src/github.com/magefile/mage
go run bootstrap.go
```

install `packr` by running :

```shell
go get -u -v github.com/gobuffalo/packr/v2/...
go get -u -v github.com/gobuffalo/packr/v2/packr2
```

### Commands

- `mage build` : rebuild the plugin and install the binary
- `mage example` : rebuild the plugin and generate examples in `/example`
- `mage clean` : remove generated go files in example folder and clean up the files packer generates

## Usage

you can use `example/example.proto` and makefile in `example/` as a reference point for your own proto definitions.  

### Google WellKnown Types

Most often we found ourselves using Google Well-Known Types as defined in [`https://developers.google.com/protocol-buffers/docs/reference/google.protobuf`](https://developers.google.com/protocol-buffers/docs/reference/google.protobuf). the following table is the filename (proto import) to package (go import) for `github.com/gogo/protobuf`

| File Name                               | Package                                             |
|-----------------------------------------|-----------------------------------------------------|
| `google/protobuf/any.proto`             | github.com/gogo/protobuf/types                      |
| `google/protobuf/api.proto`             | google.golang.org/genproto/protobuf/api             |
| `google/protobuf/compiler/plugin.proto` | github.com/gogo/protobuf/protoc-gen-gogo/plugin     |
| `google/protobuf/descriptor.proto`      | github.com/gogo/protobuf/protoc-gen-gogo/descriptor |
| `google/protobuf/duration.proto`        | github.com/gogo/protobuf/types                      |
| `google/protobuf/empty.proto`           | github.com/gogo/protobuf/types                      |
| `google/protobuf/field_mask.proto`      | github.com/gogo/protobuf/types                      |
| `google/protobuf/source_context.proto`  | google.golang.org/genproto/protobuf/source_context  |
| `google/protobuf/struct.proto`          | github.com/gogo/protobuf/types                      |
| `google/protobuf/timestamp.proto`       | github.com/gogo/protobuf/types                      |
| `google/protobuf/type.proto`            | google.golang.org/genproto/protobuf/ptype           |
| `google/protobuf/wrappers.proto`        | github.com/gogo/protobuf/types                      |
