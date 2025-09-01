# ghcr-badge-image-builder
![Image Size](https://ghcr-badge.linuxnet.io/kerwood/ghcr-badge/latest_tag?color=%2344cc11&ignore=latest&label=version) ![Image Size](https://ghcr-badge.linuxnet.io/kerwood/ghcr-badge/size?color=%2344cc11&tag=latest&label=image+size)

This repository builds the [eggplants/ghcr-badge](https://github.com/eggplants/ghcr-badge) application, but produces a smaller, more secure image compared to the original.

The official `eggplants/ghcr-badge:0.5.1` image is built in a single stage using the standard Python base image.
While functional, this approach results in a very large image that includes unnecessary build-time dependencies and contains numerous fixable vulnerabilities.

In contrast, this repository’s image is built on Google’s Distroless Python image, which is specifically optimized for security and minimalism.
This significantly reduces both the size and the number of vulnerabilities.

## Comparison

### Compressed image size
- `eggplants/ghcr-badge:0.5.1`: `381.76 MiB`
- `kerwood/ghcr-badge:0.5.1`: `32.77 MiB`

### Uncompressed image size
- `eggplants/ghcr-badge:0.5.1`: `1.05 GiB`
- `kerwood/ghcr-badge:0.5.1`: `86.2 MiB`

### Fixable vulnerabilities
- `eggplants/ghcr-badge:0.5.1`: **436** Total (Unknown: 3, Low: 54, Medium: 241, High: 138, Critical: 0)
- `kerwood/ghcr-badge:0.5.1`: **3** Total (Unknow: 0, Low: 0, Medium: 0, High: 3, Critical: 0)

This repository uses Dagger for building, but a Dockerfile is also included for those who prefer building the image manually.

## Deployment

The instance of `eggplants/ghcr-badge` is deployed on [Render](https://render.com/) using the free tier.
Because of this, you may occasionally see the message “This service has been suspended” at the end of each month when the free quota is exhausted.
You can, of course, deploy your own instance of the `ghcr-badge` application on Render as well.

If you prefer to self-host but don’t want to expose the service for everyone to use, you can restrict which repositories it works with by limiting the allowed paths.

Below is an example using my standard [Traefik](https://linuxblog.xyz/posts/traefik-3-docker-compose/) setup, configured to serve badges only for my GitHub repository by adding the ``(Path(`/`) || PathPrefix(`/kerwood/`))`` in the rule label.
```yaml
networks:
  traefik-proxy:
    external: true

services:
  ghcr_badge:
    image: ghcr.io/kerwood/ghcr-badge:0.5.1
    container_name: ghcr-badge
    restart: unless-stopped
    expose:
      - 5000
    networks:
      - traefik-proxy
    labels:
      - traefik.enable=true
      - traefik.http.services.ghcr-badge.loadbalancer.server.port=5000
      - traefik.http.routers.ghcr-badge.rule=Host(`ghcr-badge.example.org`) && (Path(`/`) || PathPrefix(`/kerwood/`))
      - traefik.http.routers.ghcr-badge.tls.certresolver=le
      - traefik.http.routers.ghcr-badge.entrypoints=websecure
      - com.centurylinklabs.watchtower.enable=true
```

## Building the image
You can build and publish the image in several ways.

### Docker
Build the image directly with Docker.
```sh
docker build -t <registry>/ghcr-badge:0.5.1 --build-arg VERSION=v0.5.1 .
```

### Dagger
Export the image to your local Docker using the Dagger module.

```sh
dagger -c 'build v0.5.1 | export-image ghcr-badge'
```

Authenticate and push the image to a registry (e.g., Docker Hub).
```sh
export PASSWD=<your-docker-hub-password>
dagger -c 'build v0.5.1 | with-registry-auth docker.io <username> env://PASSWD | publish docker.io/<username>/ghcr-badge'
```
