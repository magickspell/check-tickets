# Используем официальный образ Node.js с предустановленным yarn
FROM node:22.14-alpine3.21 as build

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем package.json и yarn.lock
COPY package.json yarn.lock ./

# Устанавливаем зависимости через yarn
RUN yarn install

# Копируем исходный код
COPY . .

# Собираем проект через yarn
RUN yarn build

# Используем легковесный образ Nginx для сервинга статики
FROM nginx:1.27.4-alpine3.21 as rd-nginx

# Копируем собранные файлы из предыдущего этапа
COPY --from=build /app/dist /usr/share/nginx/html

# Копируем конфигурацию Nginx
COPY nginx/default.conf /etc/nginx/conf.d/default.conf

# Открываем порт 80
EXPOSE 80

# # Запускаем Nginx
CMD ["nginx", "-g", "daemon off;"]