# docker112rc3-runtimefix

## First and foremost, this tool is a HACK. It is offered without guarantee or support! It modifies files generally considered private to the docker daemon!

Docker 1.12-rc3 introduced a breaking change from docker 1.12-rc2. Breaking changes are fully expected between pre-release versions (such as rc2 to rc3 here). However in some cases it can be painful. In this case to fix the compatibility problem you either need to re-create all your containers or manually update the container configuration on disk.

See [docker/docker#24343](github.com/docker/docker/issues/24343) for more details on the issue.  **If you do not see the error reported by the original poster, this tool will not help you**

This tool will manually update all the container configurations on your docker host so that you won't see this error again.
**Once this tool is run you must restart docker.**

## Usage

Run the `docker112rc3-runtimefix` on your docker host file passing in the path to the containers directory for docker
```bashtext
$ ./docker112rc3-runtimefix /var/lib/docker
```

Alternatively using the minimal docker image:
```bashtext
$ docker run --rm -v /var/lib/docker:/docker cpuguy83/docker112rc3-runtimefix:rc3
```

Then restart docker.
