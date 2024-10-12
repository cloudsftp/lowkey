mod persistence;
mod webhooks;

use actix_web::{get, web, App, HttpServer};
use anyhow::Result;
use async_nats::jetstream::{self, Context};
use dotenv::dotenv;
use log::info;
use reqwest::Client;
use rusttwald::apis::configuration::Configuration;
use std::{
    env,
    sync::{Arc, Mutex},
};

use crate::webhooks::verifier::WebhookVerifier;

struct State {
    repository: Box<dyn persistence::Repository + Send + Sync>, // TODO: rename repository
    api_configuration: Configuration, // TODO: wrap in some kind of repository
    verifier: Mutex<WebhookVerifier>, // TODO: remove mutex after using immutable webhook verifier
}

type WrappedState = Arc<State>;

#[tokio::main]
async fn main() -> Result<()> {
    dotenv().expect("could not load variables from .env");

    info!("starting server");

    env_logger::init();
    let state = bootstrap().await?;

    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(state.clone()))
            .service(webhooks::build_service())
            .service(hey)
    })
    .bind(("0.0.0.0", 6670))?
    .run()
    .await?;

    Ok(())
}

async fn bootstrap() -> Result<WrappedState> {
    let jetstream = setup_nats().await?;

    let extension_instances = get_or_create_key_value(&jetstream, "extension_instances").await?;
    let keys = get_or_create_key_value(&jetstream, "public_signing_keys_mittwald").await?;

    let repository = persistence::nats::NatsRepository {
        extension_instances,
    };

    let api_configuration = build_config();

    Ok(Arc::new(State {
        repository: Box::new(repository),
        api_configuration: api_configuration.clone(),
        verifier: Mutex::new(WebhookVerifier::new(api_configuration, keys)),
    }))
}

async fn setup_nats() -> Result<Context> {
    let nats_host = env::var("NATS_HOST")?;
    let nats_client = async_nats::connect(nats_host).await?;
    let jetstream = jetstream::new(nats_client);

    Ok(jetstream)
}

async fn get_or_create_key_value(
    jetstream: &jetstream::Context,
    bucket: &str,
) -> Result<jetstream::kv::Store> {
    let store = jetstream.get_key_value(bucket).await;
    match store {
        Ok(store) => Ok(store),
        Err(err) => {
            if let jetstream::context::KeyValueErrorKind::GetBucket = err.kind() {
                Ok(jetstream
                    .create_key_value(jetstream::kv::Config {
                        bucket: bucket.into(),
                        history: 10,
                        ..Default::default()
                    })
                    .await?)
            } else {
                Err(err.into())
            }
        }
    }
}

fn build_config() -> Configuration {
    Configuration {
        base_path: "https://api.mittwald.de".to_string(),
        user_agent: Some("lowkey via rusttwald".to_string()),
        client: Client::new(),
        basic_auth: None,
        oauth_access_token: None,
        bearer_access_token: None,
        api_key: None,
    }
}

#[get("/hey")]
async fn hey() -> Result<String, actix_web::Error> {
    Ok("hi :)\n".to_string())
}

#[get("/hey-mittwald")]
async fn hey_mittwald(data: web::Data<WrappedState>) -> Result<String, actix_web::Error> {
    let _projects = model::get_customers(&data.api_configuration)
        .await
        .map_err(actix_web::error::ErrorInternalServerError)?;

    Ok("hi :)\n".to_string())
}

mod model {
    use anyhow::Result;
    use futures::stream::FuturesUnordered;
    use futures::TryStreamExt;
    use rusttwald::apis::{
        configuration::Configuration, customer_api::customer_list_customers,
        project_api::project_list_projects,
    };

    #[derive(Debug, Clone)]
    pub struct Project {
        pub _id: String,
    }

    #[derive(Debug, Clone)]
    pub struct Customer {
        pub _id: String,
        pub _projects: Vec<Project>,
    }

    impl Customer {
        async fn from_id(config: &Configuration, customer_id: String) -> Result<Self> {
            let projects =
                project_list_projects(config, Some(&customer_id), None, None, None, None)
                    .await?
                    .iter()
                    .map(|project_response| Project {
                        _id: project_response.id.to_string(),
                    })
                    .collect();

            Ok(Customer {
                _id: customer_id.to_string(),
                _projects: projects,
            })
        }
    }

    pub async fn get_customers(config: &Configuration) -> Result<Vec<Customer>> {
        FuturesUnordered::from_iter(
            customer_list_customers(config, None, None, None, None, None)
                .await?
                .iter()
                .map(|customer| customer.customer_id.to_string())
                .map(|id| Customer::from_id(config, id)),
        )
        .try_collect()
        .await
    }
}
