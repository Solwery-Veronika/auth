gen:
	export PATH=$PATH:$(go env GOPATH)/bin 
	go generate ./...

init-db:
	cat 0001_create_table.sql | docker exec -i auth-db psql -h 0.0.0.0 -U master -f-