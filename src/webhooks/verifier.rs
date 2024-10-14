use actix_web::{http::header::HeaderMap, web};
use anyhow::Result;
use async_nats::jetstream::kv::Store;
use base64::prelude::*;
use futures::StreamExt;
use log::info;
use mittlife_cycles::verification::{
    fetcher::{KeyFetcher, PublicKeyResponse},
    headers::SignatureHeaders,
    Cache, ED25519PublicKey, Ed25519Verifier, KeyCollection, MappedHeaders, PublicKey, Verifier,
};
use rusttwald::apis::configuration::Configuration;

pub struct WebhookVerifier {
    key_collection: KeyCollection<NatsCache, ED25519PublicKey, RusttwaldFetcher>,
    verifier: Ed25519Verifier,
}

impl WebhookVerifier {
    pub fn new(api_config: Configuration, keys: Store) -> Self {
        WebhookVerifier {
            key_collection: KeyCollection::new(NatsCache { keys }, RusttwaldFetcher { api_config }),
            verifier: Ed25519Verifier {},
        }
    }

    pub async fn verify_request(&mut self, body: web::Bytes, headers: &HeaderMap) -> Result<()> {
        info!("verifying request signature");

        let headers: MappedHeaders = headers.try_into()?;
        let serial = headers.get_serial();
        let public_key = self.key_collection.get_or_fetch_key(serial).await?;

        self.verifier
            .verify_signature(&headers, &body.to_vec(), &public_key)
    }
}

struct NatsCache {
    keys: Store,
}

impl NatsCache {
    async fn try_retire_keys(&mut self) -> Result<()> {
        let mut serials = self.keys.keys().await?;
        while let Some(serial) = serials.next().await {
            self.keys.delete(serial?).await?
        }

        Ok(())
    }
}

#[async_trait::async_trait]
impl Cache<ED25519PublicKey> for NatsCache {
    async fn get(&self, serial: &str) -> Option<ED25519PublicKey> {
        let bytes = *self
            .keys
            .get(serial)
            .await
            .inspect_err(|err| info!("error while reading public key from nats: {}", err))
            .ok()??
            .first_chunk()?;

        Some(ED25519PublicKey::new(bytes))
    }

    async fn set(&mut self, serial: String, value: ED25519PublicKey) -> Result<()> {
        let bytes = value.get_bytes().to_vec().into();
        self.keys.create(serial, bytes).await?;

        Ok(())
    }

    async fn retire_keys(&mut self) {
        let _ = self
            .try_retire_keys()
            .await
            .inspect_err(|err| info!("error while retiring keys in nats: {}", err));
        // ignore error after logging
    }
}

struct RusttwaldFetcher {
    api_config: Configuration,
}

#[async_trait::async_trait]
impl KeyFetcher for RusttwaldFetcher {
    async fn fetch(&self, serial: &str) -> Result<PublicKeyResponse> {
        info!("api config for public key provider: {:?}", self.api_config);

        let public_key_response =
            rusttwald::apis::marketplace_api::extension_get_public_key(&self.api_config, serial)
                .await?;

        Ok(PublicKeyResponse {
            key_base64: BASE64_STANDARD.encode(public_key_response.key),
            serial: public_key_response.serial.to_string(),
        })
    }
}
