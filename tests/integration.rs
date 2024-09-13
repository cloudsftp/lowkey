mod containers;

#[tokio::test]
async fn test_with_nats() {
    let containers = containers::setup_containers().await;

    println!("{containers:?}");
}

#[tokio::test]
async fn second_test() {
    containers::setup_containers().await;
}
