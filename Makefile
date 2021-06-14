prod:
	killall gored || echo "[gored] not running"
	git pull
	go1.15 build
	chmod +x gored
	echo "[gored] starting server"
	GIN_MODE=release ./gored &
