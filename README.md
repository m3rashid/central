### Central

Golang micro-services for the distributed architectures

I was researching a lot about Google OAuth2, openID connect, scopes and permissions these days. Building this project, I learnt a lot about openID and OAuth2 and its importance and the security it gives, especially interacting with 3rd party services.

I created an OIDC server and a sample client (contacts)
Like in Google cloud, you first register your app to get your client_id and client_secret.

Then, init the flow
- check if a session exists, redirects to select user
- else prompts to login/register
- after successful flow, gives a callback on the app's success redirect url 
- and vice versa with failure
- scopes are not yet handled on the frontend

There are a lot of nuances and flaws in this flow as of now (v1 quickly built in 2 days).
Currently, it can run on its own to manage users in any external app. Other required things to be added soon

This project was created only to check the viability (POC) of the idea. It is not intended to be used in production. Checkout [Central V2](https://github.com/m3rashid/central-v2) for the full implementation.

---

**How to run**

1. with docker, run `docker-compose up -d`
2. without docker
   - install [golang](https://golang.org/doc/install), [postgres](https://www.postgresql.org/download/), [redis](https://redis.io/download), [air](github.com/cosmtrek/air@latest), [templ cli](github.com/a-h/templ/cmd/templ@latest)
   - put the required environment variables in the `sample.env` file into a `.env` file
   - run `cd auth && templ generate && cd ../campaigns && templ generate && cd -`
   - run `cd auth && air`
   - run `cd campaigns && air`

---

**Objective**

Trying to build a complete ecosystem powered by

- home-grown OAuth2 system
- OIDC
- scoped data sharing
- permissions
- and much more ...
