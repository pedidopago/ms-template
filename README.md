# xyzservice


## Testing

### Run Locally (self hosted database)

Ensure that a database "xyzservice" exists on your test server;
Ensure that a mariadb user exists with full privileges on the "ms_xyzservice" database;

```sh
export DEV_DB_CS="testuser:123456789@tcp(localhost)/ms_xyzservice?parseTime=true"
mage run
```
### Run with docker-compose

```sh
mage composerun
```

## Prerequisites

### MAGE
#### Installation:

```sh
mkdir -p ~/tmp
cd ~/tmp
git clone https://github.com/magefile/mage
cd mage
go run bootstrap.go
cd ..
rm -rf mage
```

#### List Targets

```sh
mage -l
```