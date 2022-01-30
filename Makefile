host=10.11.99.1
binary_name=remarkable_change_suspend_screen


build:
	go get ./...
	env GOOS=linux GOARCH=arm GOARM=7 DEBUG=0 go build -o $(binary_name)

install: build
	scp $(binary_name) root@$(host):/usr/share/remarkable/scripts/$(binary_name)
	scp ./scripts/$(binary_name).service root@$(host):/etc/systemd/system/$(binary_name).service
	ssh root@$(host) systemctl daemon-reload
	ssh root@$(host) systemctl enable $(binary_name)
	ssh root@$(host) systemctl restart $(binary_name)

	echo "Installed"


debug:
	ssh root@$(host) journalctl --unit $(binary_name) -f 

stop:
	ssh root@$(host) systemctl stop $(binary_name)

