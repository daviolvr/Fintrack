FROM node:20-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

# Build da aplicação
RUN npm run build

# Servir os arquivos estáticos
RUN npm install -g serve
EXPOSE 5173

CMD ["serve", "-s", "dist", "-l", "5173"]