start /b go run main.go 1111 1 1 10000
start /b go run main.go 1112 1 2 10
start /b go run main.go 1113 1 3 100
start /b go run main.go 1114 1 4 10000
start /b go run main.go 1115 1 5 10000

timeout /t 1

start /b go run main.go 2111 2 1 1
start /b go run main.go 2112 2 2 1

timeout /t 1

start /b go run main.go 0111 0 1 1
start /b go run main.go 0112 0 2 1
start /b go run main.go 0113 0 3 1
start /b go run main.go 0114 0 4 1
start /b go run main.go 0115 0 5 1
start /b go run main.go 0116 0 6 1
start /b go run main.go 0117 0 7 1
start /b go run main.go 0118 0 8 1
start /b go run main.go 0119 0 9 1