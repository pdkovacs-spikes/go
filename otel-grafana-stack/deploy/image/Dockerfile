FROM ubuntu

RUN apt-get update && apt-get install iproute2 netcat-traditional dnsutils curl -y
COPY ./otel-grafana-stack /bin
# We want to use the same log-processing pipeline that the default quickpizza env
# provides for containers labeled "applicatioin". The alloy configuration must,
# nevertheless, be shaped such as to keep metric-processing pipelines separate
# between legacy scraping and the OTLP push collection, which this test service uses
# in order to avoid collisions between relabeling and transformations:
LABEL service.type="application" 
EXPOSE 2223
CMD [ "/bin/otel-grafana-stack" ]
