# syntax=docker/dockerfile:1

##
## Build the static files
##

FROM node:22-alpine AS build-stage

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

# Vite inlines this at build time. The browser calls the backend directly,
# so this must be a URL reachable from the user's machine, not a container name.
ARG VITE_API_URL=http://localhost:8080
ENV VITE_API_URL=$VITE_API_URL

RUN npm run build

##
## Serve the static files with nginx
##

FROM nginx:1.29-alpine AS release-stage

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build-stage /app/dist /usr/share/nginx/html

EXPOSE 80
