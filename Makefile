run:
	go run cmd/main.go

install-pkg:
	go get -u github.com/go-sql-driver/mysql
	go get -u github.com/gorilla/mux
	go get -u github.com/qnstdx/shu.rl	