build: 
	go build -o bin/deployer
	cp config-example.json bin/config.json
	cp contract-manifest-example.json bin/contract-manifest.json
	cp setup-output-example.json bin/setup-output.json

