use actix_web::{put, web};
use serde::Deserialize;

use crate::WrappedState;

pub fn build_service() -> actix_web::Scope {
    web::scope("extension").service(register_extension_instance)
}

#[derive(Debug, Deserialize)]
struct AddedToContextPath {
    instance_id: String,
    context_id: String,
}

#[put("/added/{instance_id}/{context_id}")]
async fn register_extension_instance(
    data: web::Data<WrappedState>,
    path: web::Path<AddedToContextPath>,
) -> Result<String, actix_web::Error> {
    data.repository
        .register_extension_instance(&path.instance_id, &path.context_id)
        .await
        .map_err(actix_web::error::ErrorInternalServerError)?;

    Ok("Ok".to_string())
}
