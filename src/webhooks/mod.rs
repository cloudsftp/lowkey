pub mod verifier;

use actix_web::{
    error::{ErrorBadRequest, ErrorInternalServerError},
    post, web, HttpRequest,
};
use log::{error, info};
use serde::Deserialize;

use crate::{persistence::ExtensionInstance, WrappedState};

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

#[derive(Debug, Deserialize)]
struct AddedToContextBody {
    secret: String,
}

#[post("/added/{instance_id}/{context_id}")]
async fn added(
    payload: web::Bytes,
    request: HttpRequest,
    state: web::Data<WrappedState>,
    path: web::Path<AddedToContextPath>,
) -> Result<String, actix_web::Error> {
    info!("Received webhook to add an extension instance ({:?})", path);
    verify_webhook(&state, &payload, &request).await?;

    let body: AddedToContextBody = serde_json::de::from_slice(&payload).map_err(ErrorBadRequest)?;

    let instance = ExtensionInstance {
        id: path.instance_id.clone(),
        context_id: path.context_id.clone(),
        secret: body.secret.clone(),
    };

    state
        .repository
        .create_extension_instance(&instance)
        .await
        .map_err(ErrorInternalServerError)?;

    info!("Added extension instance ({:?})", path);

    Ok("Ok".to_string())
}

#[derive(Debug, Deserialize)]
struct UpdatedPath {
    instance_id: String,
}

#[post("/updated/{instance_id}")]
async fn updated(
    payload: web::Bytes,
    request: HttpRequest,
    state: web::Data<WrappedState>,
    path: web::Path<UpdatedPath>,
) -> Result<String, actix_web::Error> {
    info!(
        "Received webhook to update an extension instance ({:?})",
        path
    );
    verify_webhook(&state, &payload, &request).await?;

    let _ = path.instance_id;

    info!("Updated extension instance ({:?})", path);

    Ok("Ok".to_string())
}

#[derive(Debug, Deserialize)]
struct SecretRotatedPath {
    instance_id: String,
}

/*
#[derive(Debug, Deserialize)]
struct SecretRotatedBody {
    secret: String,
}
*/

#[post("/secret-rotated/{instance_id}")]
async fn secret_rotated(
    payload: web::Bytes,
    request: HttpRequest,
    state: web::Data<WrappedState>,
    path: web::Path<SecretRotatedPath>,
    //body: web::Json<SecretRotatedBody>,
) -> Result<String, actix_web::Error> {
    info!(
        "Received webhook to rotate an extension instance secret ({:?})",
        path,
    );
    verify_webhook(&state, &payload, &request).await?;

    let _ = path.instance_id;
    //let _ = body.secret;

    info!("Rotated extension instance secret ({:?})", path);

    Ok("Ok".to_string())
}

#[derive(Debug, Deserialize)]
struct RemovedPath {
    instance_id: String,
}

#[post("/removed/{instance_id}")]
async fn removed(
    payload: web::Bytes,
    request: HttpRequest,
    state: web::Data<WrappedState>,
    path: web::Path<RemovedPath>,
) -> Result<String, actix_web::Error> {
    info!(
        "Received webhook to remove an extension instance ({:?})",
        path,
    );
    verify_webhook(&state, &payload, &request).await?;

    state
        .repository
        .delete_extension_instance(&path.instance_id)
        .await
        .map_err(ErrorInternalServerError)?;

    info!("Removed extension instance ({:?})", path);

    Ok("Ok".to_string())
}

async fn verify_webhook(
    state: &web::Data<WrappedState>,
    payload: &web::Bytes,
    request: &HttpRequest,
) -> Result<(), actix_web::Error> {
    state
        .verifier
        .verify_request(payload.clone(), request.headers())
        .await
        .inspect_err(|err| error!("Could not verify webhook: {}", err))
        .map_err(ErrorBadRequest)
}
