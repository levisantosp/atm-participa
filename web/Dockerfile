FROM node:24.19.0-alpine

WORKDIR /app

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml turbo.json /app/
COPY web/package.json /app/web/package.json
COPY api-client/package.json /app/api-client/package.json
COPY ui/package.json /app/ui/package.json
COPY admin/package.json /app/admin/package.json

RUN npm i -g pnpm
RUN pn ci

COPY . .

RUN pn build:web

CMD ["pn", "start:web"]
