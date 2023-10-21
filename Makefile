localstack:
	docker run --rm -p 4566:4566 -p 4510-4559:4510-4559 -v /var/run/docker.sock:/var/run/docker.sock --name localstack localstack/localstack
