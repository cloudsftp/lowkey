# Operating

extension id: 0c3f42aa-536b-4562-abaf-ccbd91f53bda
contributor id: e77eb8f9-9b25-4e13-b0ab-92f6eaa8d3f7

## Nats

### Get extension data

``` sh
nats --server=localhost:6675 kv get extension_instances 40391011-7db4-4a9b-b55d-7d6f4c7de360 --raw | bsondump --quiet | jq
```
