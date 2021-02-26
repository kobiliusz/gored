prod:
	killall gored || echo "[gored] not running"
	git pull
	go build
	chmod +x gored
	echo "[gored] starting server"
	GIN_MODE=release ./gored &
