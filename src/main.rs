use reqwest::Client;
use rusttwald::apis::{
    configuration::{ApiKey, Configuration},
    customer_api::customer_get_customer,
    project_api::{project_list_projects, project_list_servers},
};

#[tokio::main]
async fn main() {
    let api_token = env!("MITTWALD_API_TOKEN");

    let client = Client::new();

    let config = Configuration {
        base_path: "https://api.mittwald.de/v2/".to_string(),
        user_agent: Some("rusttwald - Unofficial Rust API Client".to_string()),
        client,
        basic_auth: None,
        oauth_access_token: None,
        bearer_access_token: None,
        api_key: Some(ApiKey{
            prefix: None,
            key: api_token.to_string(),
        }),
    };

    let result = customer_get_customer(&config, "self").await;

    //project_list_servers(&config, customer_id, limit, page, skip)

    //project_list_projects(&config, customer_id, server_id, None, None, None);

    println!("Hello, world!\n{:?}", result);
}
