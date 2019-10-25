# Cyphernode admin app

## Test hydra

### Create network
````bash
docker network create hydra
````

### Create database container
````bash
docker run --network hydra \
  --name hydra-test--postgres \
  -e POSTGRES_USER=hydra \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=hydra \
  -d postgres:9.6
````

### Export some environment variables
````bash
export SECRETS_SYSTEM="this_needs_to_be_the_same_always_and_also_very_$3cuR3-._"
````

```bash
export DSN=postgres://hydra:secret@hydra-test--postgres:5432/hydra?sslmode=disable
```

### Run hydra migration to create database structure
```bash
docker run -it --rm \
  --network hydra \
  oryd/hydra:v1.0.0 \
  migrate sql --yes $DSN
```

### Run hydra service

```bash
docker run -d \
  --name hydra-test--hydra \
  --network hydra \
  -p 9000:4444 \
  -p 9001:4445 \
  -e SECRETS_SYSTEM=$SECRETS_SYSTEM \
  -e DSN=$DSN \
  -e URLS_SELF_ISSUER=http://127.0.0.1:9000/ \
  -e URLS_CONSENT=http://127.0.0.1:3030/hydra/consent \
  -e URLS_LOGIN=http://127.0.0.1:3030/hydra/login \
  oryd/hydra:v1.0.0 serve all --dangerous-force-http
```

### Create oauth2 client

````bash
docker run --rm -it \
  --network hydra \
  oryd/hydra:v1.0.0 \
  clients create \
    --endpoint http://hydra-test--hydra:4445 \
    --id another-consumer \
    --secret consumer-secret \
    -g authorization_code,refresh_token \
    -r token,code,id_token \
    --scope openid,offline \
    --callbacks http://127.0.0.1:9010/callback
````

### List oauth2 clients

```bash
docker run --rm -it \
  --network hydra \
  oryd/hydra:v1.0.0 \
  clients list \
    --endpoint http://hydra-test--hydra:4445
```

### Delete oauth2 client

```bash
docker run --rm -it \
  --network hydra \
  oryd/hydra:v1.0.0 \
  clients delete another-consumer \
    --endpoint http://hydra-test--hydra:4445
```

### Start login/consent process 

Note: start cyphernode admin beforehand with env CNA_DISABLE_HYDRA_SYNC

```bash
docker run --rm -it \
  --network hydra \
  -p 9010:9010 \
  oryd/hydra:v1.0.0 \
  token user \
    --port 9010 \
    --auth-url http://127.0.0.1:9000/oauth2/auth \
    --token-url http://hydra-test--hydra:4444/oauth2/token \
    --client-id another-consumer \
    --client-secret consumer-secret \
    --scope openid,offline \
    --redirect http://127.0.0.1:9010/callback
```