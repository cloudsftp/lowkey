use actix_web::{web, App, HttpServer};
use anyhow::Result;
use async_nats::jetstream;
use dotenv::dotenv;
use reqwest::Client;
use rusttwald::apis::configuration::{ApiKey, Configuration};
use std::env;

#[tokio::main]
async fn main() -> Result<()> {
    dotenv().expect("could not load variables from .env");
    let state = bootstrap().await?;

    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(state.clone()))
            .route("/hey", web::get().to(hello_mittwald))
    })
    .bind(("0.0.0.0", 6670))?
    .run()
    .await?;

    Ok(())
}

#[derive(Debug, Clone)]
struct State {
    _jetstream: jetstream::Context,
    mittwald_api_configuration: Configuration,
}

async fn bootstrap() -> Result<State> {
    let jetstream = setup_nats().await?;

    let mittwald_api_configuration = build_config();

    Ok(State {
        _jetstream: jetstream,
        mittwald_api_configuration,
    })
}

async fn setup_nats() -> Result<jetstream::Context> {
    let nats_host = env::var("NATS_HOST").expect("could not get NATS_HOST from the environment");
    let nats_client = async_nats::connect(nats_host).await?;
    let jetstream = jetstream::new(nats_client);

    let _ = get_or_create_key_value(&jetstream, "extension_instances").await;

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
    let api_token = env::var("MITTWALD_API_TOKEN")
        .expect("could not get MITTWALD_API_TOKEN from the environment");

    let client = Client::new();

    Configuration {
        base_path: "https://api.mittwald.de".to_string(),
        user_agent: Some("lowkey via rusttwald".to_string()),
        client,
        basic_auth: None,
        oauth_access_token: None,
        bearer_access_token: None,
        api_key: Some(ApiKey {
            prefix: None,
            key: api_token.to_string(),
        }),
    }
}

async fn hello_mittwald(data: web::Data<State>) -> Result<String, actix_web::Error> {
    let projects = model::get_customers(&data.mittwald_api_configuration)
        .await
        .expect("TODO: actix error");

    Ok(format!("{:?}", projects))
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
