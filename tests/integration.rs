mod containers;

#[tokio::test]
async fn test_with_nats() {
    containers::setup_containers().await;
}
