FROM node:alpine3.12

COPY ./dist ~/run
COPY node_modules ~/run/node_modules

WORKDIR ~/run
EXPOSE 8000
CMD ["node","index.js"]
