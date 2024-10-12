## Local Development

Requirements:

- libopenssl-3-devel (OpenSUSE Tumbleweed)

This repo is powered by [dagger](https://dagger.io/).
Install it [here](https://docs.dagger.io/quickstart/cli).

### Building

``` sh
dagger call build --source=.
```

To export the result

``` sh
dagger call build --source=. export --path=bin/lowkey
```

Cargo build also works of course

``` sh
cargo build --release
```

#### Unit Tests

``` sh
cargo test
```

Or

``` sh
dagger call test --source=.
```

### Integration Tests

Run these commands:

``` sh
dagger call integration-lowkey-service --source=. --local-dev-service=tcp://localhost:8080
```

``` sh
dagger call integration-local-dev-service --source=. --lowkey-service=tcp://localhost:6670 up --ports 8080:8080
```

``` sh
dagger call integration-drive-tests --source=integration --lowkey-service=tcp://localhost:6670 --local-dev-service=tcp://localhost:8080
```

Broken:

``` sh
dagger call integration-test --source=.
```

### Spin Up Local Development Instance

``` sh
dagger call build-test-service --source=. up
```

Test in another shell

``` sh
curl localhost:6670/hey
```


