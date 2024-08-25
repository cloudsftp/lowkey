use anyhow::Error;
use reqwest::Client;
use rusttwald::apis::{
    configuration::{ApiKey, Configuration},
    customer_api::customer_list_customers,
};

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

    let customers = get_customers(&config).await?;

    println!("{:?}", customers);

    Ok(())
}

#[derive(Debug, Clone)]
pub struct Customer {
    pub id: String,
}

pub async fn get_customers(config: &Configuration) -> Result<Vec<Customer>, Error> {
    let customer_list = customer_list_customers(config, None, None, None, None, None).await?;
    Ok(customer_list
        .iter()
        .map(|customer| Customer {
            id: customer.customer_id.to_string(),
        })
        .collect())
}
