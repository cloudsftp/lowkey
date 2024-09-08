use actix_web::{post, web};
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

#[post("/added/{instance_id}/{context_id}")]
async fn register_extension_instance(
    data: web::Data<WrappedState>,
    path: web::Path<AddedToContextPath>,
) -> Result<String, actix_web::Error> {
    data.repository
        .create_extension_instance(&path.instance_id, &path.context_id)
        .await
        .map_err(actix_web::error::ErrorInternalServerError)?;

    Ok("Ok".to_string())
}

#[derive(Debug, Deserialize)]
struct UpdatedPath {
    _instance_id: String,
}

#[post("/updated/{instance_id}")]
async fn extension_instance_updated(
    _: web::Data<WrappedState>,
    _: web::Path<UpdatedPath>,
) -> Result<String, actix_web::Error> {
    Ok("Ok".to_string())
}

#[derive(Debug, Deserialize)]
struct SecretRotatedPath {
    _instance_id: String,
}

#[post("/secret-rotated/{instance_id}")]
async fn secret_rotated(
    _: web::Data<WrappedState>,
    _: web::Path<SecretRotatedPath>,
) -> Result<String, actix_web::Error> {
    Ok("Ok".to_string())
}

#[derive(Debug, Deserialize)]
struct RemovedPath {
    instance_id: String,
}

#[post("/removed/{instance_id}")]
async fn extension_instance_removed(
    data: web::Data<WrappedState>,
    path: web::Path<RemovedPath>,
) -> Result<String, actix_web::Error> {
    data.repository
        .delete_extension_instance(&path.instance_id)
        .await
        .map_err(actix_web::error::ErrorInternalServerError)?;

    Ok("Ok".to_string())
}
