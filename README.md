``` sh
curl https://api.lowkey.energiesandsuch.com/hey
```

### Local Development

This repo is powered by [dagger](https://dagger.io/).
Install it [here](https://docs.dagger.io/quickstart/cli).

#### Building

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

#### Integration Tests

``` sh
dagger call integration-test --source=.
```

#### Spin Up Local Development Instance

``` sh
dagger call build-test-service --source=. up
```

Test in another shell

``` sh
curl localhost:6670/hey
```

