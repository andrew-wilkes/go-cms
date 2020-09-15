b:
	php scripts/build.php

deploy:
	php scripts/deploy.php

test:
	go test gocms/cmd/main

copyfiles:
	php scripts/copytestfiles.php

g:
	go run scripts/generate/generate.go
