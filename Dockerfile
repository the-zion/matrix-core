FROM debian:stable-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y
ARG filename
COPY ./bin/${filename} /home/app/
RUN useradd matrix && chown -R matrix:matrix /home/app/ && chmod 700 /home/app/
WORKDIR /home/app/
USER matrix
EXPOSE 8000
EXPOSE 9000
