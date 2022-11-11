FROM debian:stable-slim
ARG filename
RUN echo "${filename}"
COPY ./bin/${filename} /app/

WORKDIR /app/

EXPOSE 8000
EXPOSE 9000
