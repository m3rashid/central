### Central

Golang micro-services for the distributed architectures

This project was created only to check the viability (POC) of the idea. It is not intended to be used in production. Checkout [Central V2](https://github.com/m3rashid/central-v2) for the full implementation.

---

<br />

**How to run**

1. with docker

```bash
docker-compose up -d
```

2. without docker
   - install [golang](https://golang.org/doc/install), [postgres](https://www.postgresql.org/download/), [redis](https://redis.io/download), [air](github.com/cosmtrek/air@latest), [templ cli](github.com/a-h/templ/cmd/templ@latest)
   - put the required environment variables in the `sample.env` file into a `.env` file
   - run `cd auth && templ generate && cd ../campaigns && templ generate && cd -`
   - run `cd auth && air`
   - run `cd campaigns && air`

---

<br />

**Objective**

Trying to build a complete ecosystem powered by

- home-grown OAuth2 system
- OIDC
- scoped data sharing
- permissions
- and much more ...
