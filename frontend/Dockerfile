FROM node:18-alpine as build

WORKDIR /app

COPY package.json ./
COPY package-lock.json* ./

RUN npm cache clean --force
RUN npm install

COPY . .

RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"] 