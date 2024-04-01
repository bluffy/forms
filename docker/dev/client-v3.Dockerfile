# base image
#FROM node:10.18.0-slim
FROM node:alpine3.18 AS vue 
RUN apk add --no-cache bash su-exec
COPY VERSION /.
WORKDIR /webapp/app

COPY docker/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh
CMD bash -c "npm install && npm run build-dev"
#CMD bash -c "tail -f /dev/null"

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
