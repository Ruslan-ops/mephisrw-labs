
FROM node:latest AS builder

WORKDIR /app

COPY . .


# Define build arguments
ARG BASE_URL
ARG LAB12_URL
# Set environment variables inside the container
ENV BASE_URL=${BASE_URL}
ENV LAB12_URL=${LAB12_URL}

RUN npm install && npm run build


FROM nginx:alpine

WORKDIR /usr/share/nginx/html

RUN rm -rf ./*

COPY --from=builder /app/dist/spa .

ENTRYPOINT ["nginx", "-g", "daemon off;"]
