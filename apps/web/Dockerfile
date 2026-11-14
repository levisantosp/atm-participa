FROM node:26.13.1-alpine

WORKDIR /app

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml turbo.json /app/
COPY apps/web/package.json /app/apps/web/package.json
COPY pkg/api-client/package.json /app/pkg/api-client/package.json
COPY pkg/ui/package.json /app/pkg/ui/package.json
COPY apps/admin/package.json /app/apps/admin/package.json

RUN npm i -g pnpm
RUN pn ci

COPY . .

RUN pn build:web

CMD ["pn", "start:web"]
