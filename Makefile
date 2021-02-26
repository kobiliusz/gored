prod:
	killall gored
	git pull
	go build
	chmod +x gored
	GIN_MODE=release ./gored &
