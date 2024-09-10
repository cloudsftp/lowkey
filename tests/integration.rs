use testcontainers_modules::{
    nats,
    testcontainers::{runners::AsyncRunner, GenericImage},
};

struct ContainerInfo {
    host: String,
    port: String,
}

#[tokio::test]
async fn test_with_nats() {
    let nats_info = start_nats().await;
    let lowkey_info = start_lowkey(nats_info).await;
    let local_dev_info = start_local_dev(lowkey_info).await;

    let local_dev_container = GenericImage::new("mittwald/marketplace-local-dev-server", "1.3.0")
        .pull_image()
        .await
        .expect("could not pull image for local dev server")
        .start()
        .await
        .expect("could not start local dev server");
}

async fn start_nats() -> ContainerInfo {
    let container = nats::Nats::default()
        .start()
        .await
        .expect("could not start nats");

    ContainerInfo {
        host: container
            .get_host()
            .await
            .expect("could not get host of nats container")
            .to_string(),
        port: container
            .get_host_port_ipv4(8222)
            .await
            .expect("could not get port of nats container")
            .to_string(),
    }
}

async fn start_lowkey(nats_info: ContainerInfo) -> ContainerInfo {
    let container = GenericImage::new("lowkey", "test")
        .start()
        .await
        .expect("could not start lowkey");

    ContainerInfo {
        host: container
            .get_host()
            .await
            .expect("could not get host of lowkey container")
            .to_string(),
        port: container
            .get_host_port_ipv4(6670)
            .await
            .expect("could not get port of lowkey container")
            .to_string(),
    }
}

async fn start_local_dev(lowkey_info: ContainerInfo) -> ContainerInfo {
    let container = GenericImage::new("mittwald/marketplace-local-dev-server", "1.3.0")
        .pull_image()
        .await
        .expect("could not pull lowkey image")
        .start()
        .await
        .expect("could not start lowkey");

    ContainerInfo {
        host: container
            .get_host()
            .await
            .expect("could not get host of  container")
            .to_string(),
        port: container
            .get_host_port_ipv4(6670)
            .await
            .expect("could not get port of lowkey container")
            .to_string(),
    }
}
