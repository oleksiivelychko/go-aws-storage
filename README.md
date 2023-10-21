### go-aws-storage

### (local) usage of AWS S3 via Cobra and/or AWS CLI.

- list buckets via AWS CLI
```
aws s3api list-buckets --endpoint-url http://localhost:4566 --profile localstack
```
- list buckets via CLI
```
go run main.go list-buckets
```
---

---
※ References:
- [s3api](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/s3api/index.html)
