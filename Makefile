##Tail
prebuild:
	smdcatalog
debug:

qrun:
	go run cmd/pickcheck/main.go
test:

install:
	cd cmd/pickcheck && go install
	cd cmd/pickinfo && go install
clean:

