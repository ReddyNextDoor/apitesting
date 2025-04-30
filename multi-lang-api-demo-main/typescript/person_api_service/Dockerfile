# Use official Node.js LTS image for build
FROM node:20-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
COPY openapi.yaml ./openapi.yaml
RUN npm run build

# Production image
FROM node:20-alpine AS runtime
WORKDIR /app
COPY --from=build /app/dist ./dist
COPY package*.json ./
COPY openapi.yaml ./openapi.yaml
RUN npm install --omit=dev
EXPOSE 8000
ENV NODE_ENV=production
CMD ["node", "dist/index.js"]
