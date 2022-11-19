FROM debian:stable-slim
ARG filename
RUN echo "${filename}"
COPY ./bin/${filename} /home/app/

WORKDIR /home/app/

EXPOSE 8000
EXPOSE 9000
