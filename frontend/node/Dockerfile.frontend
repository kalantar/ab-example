
FROM node:18
# RUN mkdir -p /usr/app
WORKDIR /usr/app
COPY frontend/node/package*.json .

RUN npm install

COPY frontend/node .

EXPOSE 8091
CMD ["node", "main.js"]