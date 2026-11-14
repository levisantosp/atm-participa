FROM node:26.13.1-alpine AS build

WORKDIR /app

ARG VITE_CDN_URL
ENV VITE_CDN_URL=$VITE_CDN_URL

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml turbo.json /app/
COPY apps/web/package.json /app/apps/web/package.json
COPY pkg/api-client/package.json /app/pkg/api-client/package.json
COPY pkg/ui/package.json /app/pkg/ui/package.json
COPY apps/admin/package.json /app/apps/admin/package.json

RUN npm i -g pnpm
RUN pn ci

COPY . .

RUN printf 'VITE_CDN_URL=%s\n' "$VITE_CDN_URL" > apps/admin/.env
RUN pn build:admin

FROM nginx:alpine
COPY --from=build /app/apps/admin/dist /usr/share/nginx/html
COPY apps/admin/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
