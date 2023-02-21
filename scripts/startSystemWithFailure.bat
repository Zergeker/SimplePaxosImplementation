start /b go run main.go 1111 1 1 1
start /b go run main.go 1112 1 2 1
start /b go run main.go 1113 1 3 100000
start /b go run main.go 1114 1 4 100000
start /b go run main.go 1115 1 5 100000

timeout /t 1

start /b go run main.go 2111 2 1 1
start /b go run main.go 2112 2 2 1

timeout /t 1

start /b go run main.go 0111 0 1 1
start /b go run main.go 0112 0 2 1