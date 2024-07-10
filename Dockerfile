FROM golang:1.22
ADD ./bin/weather_app /weather_app
WORKDIR /app
COPY . .
COPY config/config_docker.yaml config/config.yaml
CMD ["/weather_app"]