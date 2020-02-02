# exer

A tool that execute build tool commands.

## Installation
### Requirement
- git

### go get
```
$ go get github.com/4afs/exer
```

## Usage

```
$ exer -run
```

execute run command at the root directory of the project.

```
$ exer -build
```

execute build command same as `exer -run`.

## Supported languages and build tools

| Build tool | Match | Run | Build |
|:--|:--|:--|:--|
| Stack | stack.yaml | `stack run` | `stack build` |
| Cargo | Cargo.toml | `cargo run` | `cargo build` |
| Spago | .spago | `spago run` | `spago build` |
| Elm | elm.json | `elm reactor` | - |
| sbt | build.sbt | `sbt run` | `sbt build` |
