use std::fmt::Debug;

use anyhow::Result;
use async_trait::async_trait;
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
struct ExtensionInstance {
    context_id: String,
}

#[async_trait]
pub trait Repository: Debug {
    async fn create_extension_instance(&self, instance_id: &str, context_id: &str) -> Result<()>;
    async fn delete_extension_instance(&self, instance_id: &str) -> Result<()>;
}

pub mod nats {
    use anyhow::Result;
    use async_nats::jetstream::kv::Store;
    use async_trait::async_trait;

    #[derive(Debug, Clone)]
    pub struct NatsRepository {
        pub extension_instances: Store,
    }

    #[async_trait]
    impl super::Repository for NatsRepository {
        async fn create_extension_instance(
            &self,
            instance_id: &str,
            context_id: &str,
        ) -> Result<()> {
            let instance = super::ExtensionInstance {
                context_id: context_id.into(),
            };

            let instance = bson::ser::to_document(&instance)?;

            let mut encoded: Vec<u8> = Vec::new();
            instance.to_writer(&mut encoded)?;

            self.extension_instances
                .create(instance_id, encoded.into())
                .await?;

            Ok(())
        }

        async fn delete_extension_instance(&self, instance_id: &str) -> Result<()> {
            self.extension_instances.delete(instance_id).await?;

            Ok(())
        }
    }
}
