use std::{env, process::Command};

use anyhow::{bail, Result};

fn main() -> Result<()> {
    let cwd = env::var("CARGO_MANIFEST_DIR")?;

    // Build the test image
    let output = Command::new("docker")
        .arg("build")
        .arg("--file")
        .arg(&format!("{cwd}/Dockerfile"))
        .arg("--force-rm")
        .arg("--tag")
        .arg("lowkey:test")
        .arg(".")
        .output()?;
    if !output.status.success() {
        eprintln!("stderr: {}", String::from_utf8(output.stderr)?);
        bail!("unable to build lowkey:test");
    }
    eprintln!("Built lowkey:test");

    // trigger recompilation when dockerfiles are modified
    println!("cargo:rerun-if-changed=Dockerfile");
    println!("cargo:rerun-if-changed=.dockerignore");

    Ok(())
}
