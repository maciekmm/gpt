all: delays.bin stoptimes.bin

delays.bin: delays/delay.go
	go build ./delays/delay.go -o delays.bin

stoptimes.bin: stoptimes/stoptimes.go
	go build ./stoptimes/stoptimes.go -o stoptimes.bin

update-stoptimes: stoptimes.bin
	env $(cat ./credentials.env) ./stoptimes.bin

update-delays: delays.bin
	env $(cat ./credentials.env) ./delays.bin

clean:
	rm stoptimes.bin delays.bin
