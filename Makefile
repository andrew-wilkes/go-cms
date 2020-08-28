build:
	go build -o /build /cmd/main/gocms.go

deploy:
	php scripts/deploy.php
