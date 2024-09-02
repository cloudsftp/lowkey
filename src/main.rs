use actix_web::{web, App, HttpServer};
use dotenv::dotenv;
use reqwest::Client;
use rusttwald::apis::configuration::{ApiKey, Configuration};
use std::env;

#[tokio::main]
async fn main() -> std::io::Result<()> {
    dotenv().expect("could not load variables from .env");

    HttpServer::new(|| {
        App::new()
            .app_data(web::Data::new(State {
                config: build_config(),
            }))
            .route("/hey", web::get().to(hello_mittwald))
    })
    .bind(("0.0.0.0", 6670))?
    .run()
    .await
}

struct State {
    config: Configuration,
}

fn build_config() -> Configuration {
    let api_token = env::var("MITTWALD_API_TOKEN").expect("could not load the mittwald API token");

    let client = Client::new();

    Configuration {
        base_path: "https://api.mittwald.de".to_string(),
        user_agent: Some("rusttwald - Unofficial Rust API Client".to_string()),
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
    let projects = model::get_customers(&data.config)
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
        pub id: String,
    }

    #[derive(Debug, Clone)]
    pub struct Customer {
        pub id: String,
        pub projects: Vec<Project>,
    }

    impl Customer {
        async fn from_id(config: &Configuration, customer_id: String) -> Result<Self> {
            let projects =
                project_list_projects(config, Some(&customer_id), None, None, None, None)
                    .await?
                    .iter()
                    .map(|project_response| Project {
                        id: project_response.id.to_string(),
                    })
                    .collect();

            Ok(Customer {
                id: customer_id.to_string(),
                projects,
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
