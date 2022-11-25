FROM debian:stable-slim
ARG filename
COPY ./bin/${filename} /home/app/
#RUN useradd matrix && chown -R matrix:matrix /home/app/ && chmod 700 /home/app/
WORKDIR /home/app/
#USER matrix
EXPOSE 8000
EXPOSE 9000
