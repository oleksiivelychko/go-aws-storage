### go-aws-storage

### (local) usage of AWS S3 via Cobra and/or AWS CLI.

- create bucket via AWS CLI
```
aws s3api create-bucket --bucket my-bucket --endpoint-url http://localhost:4566 --profile localstack
```
- create bucket via CLI
```
go run main.go --config=config.yaml create-bucket --name=my-bucket [--isPublic=true]
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

---
※ References:
- [s3api](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/s3api/index.html)
- [Bucket naming rules](https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html)
