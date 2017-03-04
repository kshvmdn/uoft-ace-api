www: main.go app.go model.go controller.go
	go build -o www $^

clean:
	rm www
