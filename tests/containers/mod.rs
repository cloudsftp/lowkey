use anyhow::Result;
use testcontainers_modules::{
    nats,
    testcontainers::{core::ContainerPort, runners::AsyncRunner, GenericImage},
};
use tokio::sync::OnceCell;

struct ContainerInfo {
    host: String,
    port: String,
}

static ONCE: OnceCell<Result<()>> = OnceCell::const_new();

pub async fn setup_containers() -> &'static Result<()> {
    ONCE.get_or_init(setup_constainers_once).await
}

const LOWKEY_TESTIMAGE_NAME: &str = "lowkey";
const LOWKEY_TESTIMAGE_TAG: &str = "test";

async fn setup_constainers_once() -> Result<()> {
    images::compile_image("Dockerfile", LOWKEY_TESTIMAGE_NAME, LOWKEY_TESTIMAGE_TAG)?;

    let nats_info = start_nats().await;
    let lowkey_info = start_lowkey(nats_info, LOWKEY_TESTIMAGE_NAME, LOWKEY_TESTIMAGE_TAG).await;
    let local_dev_info = start_local_dev(lowkey_info).await;

    let local_dev_container = GenericImage::new("mittwald/marketplace-local-dev-server", "1.3.0")
        .pull_image()
        .await
        .expect("could not pull image for local dev server")
        .start()
        .await
        .expect("could not start local dev server");

    Ok(())
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

async fn start_lowkey(
    nats_info: ContainerInfo,
    image_name: &str,
    image_tag: &str,
) -> ContainerInfo {
    let port = 6670;

    let container = GenericImage::new(image_name, image_tag)
        .with_exposed_port(ContainerPort::Tcp(port))
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
            .get_host_port_ipv4(port)
            .await
            .expect("could not get port of lowkey container")
            .to_string(),
    }
}

async fn start_local_dev(lowkey_info: ContainerInfo) -> ContainerInfo {
    let port = 8080;

    let container = GenericImage::new("mittwald/marketplace-local-dev-server", "1.3.0")
        .with_exposed_port(ContainerPort::Tcp(port))
        .pull_image()
        .await
        .expect("could not pull local dev server")
        .start()
        .await
        .expect("could not start local dev server");

    ContainerInfo {
        host: container
            .get_host()
            .await
            .expect("could not get host of local dev server container")
            .to_string(),
        port: container
            .get_host_port_ipv4(port)
            .await
            .expect("could not get port of local dev server container")
            .to_string(),
    }
}

mod images {
    use std::{env, process::Command};

    use anyhow::{bail, Result};

    pub fn compile_image(docker_file: &str, image_name: &str, image_tag: &str) -> Result<()> {
        let cwd = env::var("CARGO_MANIFEST_DIR")?;

        let output = Command::new("docker")
            .arg("build")
            .arg("--file")
            .arg(&format!("{cwd}/{docker_file}"))
            .arg("--force-rm")
            .arg("--tag")
            .arg(&format!("{image_name}:{image_tag}"))
            .arg(".")
            .output()?;

        if !output.status.success() {
            bail!("unable to build lowkey:test");
        }

        Ok(())
    }
}
