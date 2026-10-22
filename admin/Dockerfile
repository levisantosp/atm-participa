FROM node:24.19.0-alpine AS build

WORKDIR /app

ARG VITE_CDN_URL
ENV VITE_CDN_URL=$VITE_CDN_URL

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml turbo.json /app/
COPY web/package.json /app/web/package.json
COPY api-client/package.json /app/api-client/package.json
COPY ui/package.json /app/ui/package.json
COPY admin/package.json /app/admin/package.json

RUN npm i -g pnpm
RUN pn ci

COPY . .

RUN printf 'VITE_CDN_URL=%s\n' "$VITE_CDN_URL" > admin/.env
RUN pn build

FROM nginx:alpine
COPY --from=build /app/admin/dist /usr/share/nginx/html
COPY admin/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
