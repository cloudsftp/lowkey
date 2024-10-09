pub mod verifier;

use actix_web::{error::ErrorBadRequest, post, web, HttpRequest};
use log::info;
use serde::Deserialize;

use crate::WrappedState;

pub fn build_service() -> actix_web::Scope {
    web::scope("extension")
        .service(added)
        .service(updated)
        .service(removed)
        .service(secret_rotated)
}

#[derive(Debug, Deserialize)]
struct AddedToContextPath {
    instance_id: String,
    context_id: String,
}

#[allow(clippy::await_holding_lock)]
#[post("/added/{instance_id}/{context_id}")]
async fn added(
    body: web::Bytes,
    request: HttpRequest,
    state: web::Data<WrappedState>,
    path: web::Path<AddedToContextPath>,
) -> Result<String, actix_web::Error> {
    state
        .verifier
        .lock() // TODO: lock not needed when using nats
        .unwrap()
        .verify_request(body, request.headers())
        .await
        .map_err(ErrorBadRequest)?;

    info!("Extension instance added: {:?}", path);

    state
        .repository
        .create_extension_instance(&path.instance_id, &path.context_id)
        .await
        .map_err(actix_web::error::ErrorInternalServerError)?;

    Ok("Ok".to_string())
}

#[derive(Debug, Deserialize)]
struct UpdatedPath {
    instance_id: String,
}

#[post("/updated/{instance_id}")]
async fn updated(
    _: web::Data<WrappedState>,
    path: web::Path<UpdatedPath>,
) -> Result<String, actix_web::Error> {
    info!("Extension instance updated: {:?}", path);

    let _ = path.instance_id;

    Ok("Ok".to_string())
}

#[derive(Debug, Deserialize)]
struct SecretRotatedPath {
    instance_id: String,
}

#[derive(Debug, Deserialize)]
struct SecretRotatedBody {
    secret: String,
}

#[post("/secret-rotated/{instance_id}")]
async fn secret_rotated(
    _: web::Data<WrappedState>,
    path: web::Path<SecretRotatedPath>,
    body: web::Json<SecretRotatedBody>,
) -> Result<String, actix_web::Error> {
    info!("Extension instance secret rotated: {:?}, {:?}", path, body);

    let _ = path.instance_id;
    let _ = body.secret;

    Ok("Ok".to_string())
}

#[derive(Debug, Deserialize)]
struct RemovedPath {
    instance_id: String,
}

#[post("/removed/{instance_id}")]
async fn removed(
    state: web::Data<WrappedState>,
    path: web::Path<RemovedPath>,
) -> Result<String, actix_web::Error> {
    info!("Extension instance removed: {:?}", path);

    state
        .repository
        .delete_extension_instance(&path.instance_id)
        .await
        .map_err(actix_web::error::ErrorInternalServerError)?;

    Ok("Ok".to_string())
}
