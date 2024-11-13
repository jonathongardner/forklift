# Forklift
A command line tool to extract (unload) archives to the simplest form.

## Usage
```sh
forklift --help
```

## Dev
```sh
go run main.go --help
go run build
```
### Dockerfile
```
docker build -t forklift:dev -f Dockerfile.dev .
docker run --rm -v $PWD:/app -it forklift:dev /bin/sh
```