FROM ubuntu

RUN apt-get update && apt-get install iproute2 -y
COPY ./cmd/otel-grafana-stack /bin
# LABEL service.type="application" - used for scraped targets in the quickpizza env
LABEL service.type="my-otlp-push-sample-service"
EXPOSE 2223
CMD [ "/bin/otel-grafana-stack" ]
