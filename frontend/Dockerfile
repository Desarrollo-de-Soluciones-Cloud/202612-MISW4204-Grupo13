FROM node:20-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY index.html ./
COPY src ./src
COPY tsconfig.json ./
COPY tsconfig.app.json ./
COPY tsconfig.node.json ./
COPY vite.config.ts ./
RUN chown -R node:node /app

USER node
EXPOSE 5173
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
