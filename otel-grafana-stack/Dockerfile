FROM ubuntu

RUN apt-get update && apt-get install iproute2 -y
COPY ./cmd/otel-grafana-stack /bin
LABEL service.type="application"
EXPOSE 2223
CMD [ "/bin/otel-grafana-stack" ]
