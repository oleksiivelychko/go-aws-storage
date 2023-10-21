### go-aws-storage

### (local) usage of AWS S3 via Cobra and/or AWS CLI.

- create bucket via AWS CLI
```
aws s3api create-bucket --bucket my-bucket --endpoint-url http://localhost:4566 --profile localstack
```
- create bucket via CLI
```
go run main.go --config=config.yaml create-bucket --name=my-bucket [--public=true]
```
---

- list buckets via AWS CLI
```
aws s3api list-buckets --endpoint-url http://localhost:4566 --profile localstack
```
- list buckets via CLI
```
go run main.go --config=config.yaml list-buckets
```
---

- put object via AWS CLI
```
aws s3api put-object --bucket my-bucket --key sample.png --body sample.png \
    --endpoint-url http://localhost:4566 --profile localstack
```
- put object via CLI
```
go run main.go put-objects --bucket=my-bucket --key=sample_1.png --key=sample_2.png
```
---

- list objects via AWS CLI
```
aws s3api list-objects --bucket my-bucket --endpoint-url http://localhost:4566 --profile localstack
```
- list objects via CLI
```
go run main.go --config=config.yaml list-objects --bucket=my-bucket
```
---

- get object via AWS CLI
```
aws s3api get-object --bucket my-bucket --key sample.png new-sample.png \
    --endpoint-url http://localhost:4566 --profile localstack OUT_FILENAME
```
- get object via CLI
```
go run main.go get-object --config=config.yaml --bucket=my-bucket --key=sample.png --path=.
```
---

- pre-sign URL via AWS CLI
```
aws s3 presign s3://BUCKET/KEY --endpoint-url http://localhost:4566 --profile localstack
```
- pre-sign URL via CLI
```
go run main.go assign-url --config=config.yaml --bucket=my-bucket --key=sample.png
```
---

- delete object via AWS CLI
```
aws s3api delete-object --bucket my-bucket --key sample.png --endpoint-url http://localhost:4566 --profile localstack
```
- delete object via CLI
```
go run main.go delete-object --config=config.yaml --bucket=my-bucket --key=sample.png
```
---

- delete bucket via AWS CLI
```
aws s3api delete-bucket --bucket my-bycket --endpoint-url http://localhost:4566 --profile localstack
```
- delete bucket via CLI
```
go run main.go delete-bucket --config=config.yaml --bucket=my-bucket
```
---

---
※ References:
- [s3api](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/s3api/index.html)
