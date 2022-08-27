.SILENT:

recv:
	go run cmd/reciver/reciver.go

send:
	go run cmd/sender/sender.go

e:
	go run cmd/experimental/main.go