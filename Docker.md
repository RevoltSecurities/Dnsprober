# Docker Usage Guide

## Building the Docker Image

```sh
sudo docker build -t dnsprober .
```

## Running dnsprober in Docker

```sh
sudo docker run --rm dnsprober --help
```

## Using a File as Input in Docker

When using files as input, you need to mount the file into the container using `-v`.

```sh
sudo docker run --rm -v /path/to/input/input.txt:/input.txt dnsprober -l /input.txt
```

## Example Commands
### Using Subdominator Output as Input
```sh
subdominator -d hackerone.com -s | sudo docker run --rm -i dnsprober --response -s
```
```sh
subdominator -d hackerone.com -s | sudo docker run --rm -i dnsprober --dns-response -s
```
```sh
subdominator -d hackerone.com -s | sudo docker run --rm -i dnsprober --response -s --cname
```

### Wildcard Filtering
```sh
sudo docker run --rm dnsprober -d freshworks.com -w wordlist.txt --wildcard-domain wildcard-domain.freshworks.com
```

### PTR Record Lookup for an IP Range
```sh
prips 157.240.19.0/24 | sudo docker run --rm -i dnsprober --ptr --response
```

### Running with Custom Resolvers
```sh
sudo docker run --rm -v /path/to/resolvers.txt:/resolvers.txt dnsprober -d example.com -r /resolvers.txt
```

### Saving Output to a File
```sh
sudo docker run --rm -v $(pwd):/output dnsprober -d example.com -o /output/results.txt
```

## Notes
- Use `-i` when piping data into `docker run`.
- Mount files using `-v` if they are needed inside the container.
- The `--rm` flag ensures the container is removed after execution.

