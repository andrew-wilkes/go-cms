b:
	go build -o build cmd/main/gocms.go

deploy:
	go run scripts/deploy/deploy.go /var/www/gocms/gocms.com

copyfiles:
	go run scripts/copyfiles/copyfiles.go /var/www/gocms/gocms.com

g:
	go run scripts/generate/generate.go

test:
	go test ./pkg/... ./cmd/...
