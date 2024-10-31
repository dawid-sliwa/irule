# Use a single stage for simplicity
FROM node:20-alpine
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install
COPY . .
RUN npm run build

# Serve with Vite preview
EXPOSE 4173
CMD ["npm", "run", "preview"]