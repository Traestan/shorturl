FROM node:lts-alpine

# устанавливаем простой HTTP-сервер для статики
RUN npm install -g http-server

# делаем каталог 'app' текущим рабочим каталогом
WORKDIR /devs/github.com/shorturlshop/frontend

# копируем оба 'package.json' и 'package-lock.json' (если есть)
COPY package*.json ./

# устанавливаем зависимости проекта за исключением зависимостей для разработки
RUN npm install -g @vue/cli
RUN npm install
RUN vue ui

# копируем файлы и каталоги проекта в текущий рабочий каталог (т.е. в каталог 'app')
COPY . .

# собираем приложение для production с минификацией
RUN npm run serve

EXPOSE 4000
CMD [ "http-server", "dist" ]