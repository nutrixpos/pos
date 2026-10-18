# Stage 1: Build
FROM node:20-alpine AS build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm add -D vitepress@next
RUN npm run docs:build

# Stage 2: Serve
FROM nginx:stable-alpine AS production-stage
COPY --from=build-stage /app/.vitepress/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]