prod:
	git pull
	go build
	chmod +x gored
	./gored
	
