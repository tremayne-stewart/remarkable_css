host=10.11.99.1


build:
	go get ./...
	env GOOS=linux GOARCH=arm GOARM=7 DEBUG=0 go build -o local_binary

install: build
	scp local_binary root@$(host):local_binary
	scp ./scripts/local_binary.service root@$(host):/etc/systemd/system/local_binary.service
	ssh root@$(host) systemctl daemon-reload
	ssh root@$(host) systemctl enable local_binary
	ssh root@$(host) systemctl restart local_binary

	echo "Installed"


debug:
	ssh root@$(host) journalctl --unit local_binary -f 

