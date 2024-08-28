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

    let projects = model::get_customers(&config).await?;

    println!("{:?}", projects);

    Ok(())
}

mod model {
    use anyhow::Error;
    use rusttwald::apis::{
        configuration::Configuration, customer_api::customer_list_customers,
        project_api::project_list_projects,
    };
    use tokio::sync::mpsc;
    use tokio_stream::StreamExt;

    #[derive(Debug, Clone)]
    pub struct Project {
        pub id: String,
    }

    #[derive(Debug, Clone)]
    pub struct Customer {
        pub id: String,
        pub projects: Vec<Project>,
    }

    pub async fn get_projects(
        config: &Configuration,
        customer_id: &str,
    ) -> Result<Vec<Project>, Error> {
        let customers = get_customers(config).await?;
        Ok(customers
            .iter()
            .map(|customer| Project {
                id: "bla".to_string(),
            })
            .collect())
    }

    pub async fn get_customers(config: &Configuration) -> Result<Vec<Customer>, Error> {
        let customer_ids: Vec<_> = customer_list_customers(config, None, None, None, None, None)
            .await?
            .iter()
            .map(|customer| customer.customer_id.to_string())
            .collect();

        let customers = tokio_stream::iter(customer_ids.clone())
            .map(|id| async move { () })
            .collect::<Vec<_>>()
            .await;

        let mut customers = Vec::new();
        /*
        while let Some(customer) = customer_receiver.recv().await {
            customers.push(customer)
        }
         */
        for id in customer_ids {
            let projects = project_list_projects(config, Some(&id), None, None, None, None)
                .await?
                .iter()
                .map(|project_response| Project {
                    id: project_response.id.to_string(),
                })
                .collect();

            customers.push(Customer { id, projects })
        }

        Ok(customers)
    }
}
