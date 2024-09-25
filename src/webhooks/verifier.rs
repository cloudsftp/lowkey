use actix_web::{http::header::HeaderMap, web};
use anyhow::Result;
use log::info;
use mittlife_cycles::verification::MappedHeaders;

#[derive(Debug)]
pub struct WebhookVerifier {}

impl WebhookVerifier {
    pub fn new() -> Self {
        WebhookVerifier {}
    }

    pub async fn verify_request(&self, body: web::Bytes, headers: &HeaderMap) -> Result<()> {
        info!("verifying request signature");

        let headers: MappedHeaders = headers.try_into()?;

        info!("headers: {:?}", headers);

        Ok(())
    }
}
