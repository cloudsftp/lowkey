use anyhow::Error;
use reqwest::Client;
use rusttwald::apis::configuration::{ApiKey, Configuration};

#[tokio::main]
async fn main() -> Result<(), Error> {
    let api_token = env!("MITTWALD_API_TOKEN");

    let client = Client::new();

    let config = Configuration {
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
    };

    //project_list_servers(&config, customer_id, limit, page, skip)

    //project_list_projects(&config, customer_id, server_id, None, None, None);

    let projects = model::get_customers(&config).await?;

    println!("{:?}", projects);

    Ok(())
}

mod model {
    use anyhow::Error;
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
        async fn from_id(config: &Configuration, customer_id: String) -> Result<Self, Error> {
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

    pub async fn get_customers(config: &Configuration) -> Result<Vec<Customer>, Error> {
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
