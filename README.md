# shorturl
Сервису 2 года многое требуется переделать
### Rest api укорачивалка ссылок

Главная.
```js
Get http://example.org/
```

Производит редирект(301) на урл указанный в источнике
```js
Get http://example.org/{shorturl}
```
Выводит статистику по переходам на коротую ссылку

```js
Get http://example.org/get/3705ce87a6
```

Генерирует короткую ссылку.
``` js
Request
Post http://example.org/add{sourceurl:"Исходная ссылка","ip":Адрес отправителя}

Response

{
    hurl:"Rgfsjtggg34"
}
```



### API для работы лк

Сервис для работы с магазином.
Включает следующий функционал
1. Авторизация,регистрация пользователя
    
  ``` js
    Request
    Post http://example.org/user/register
    json{email:'',password:''}

    Response 
    id в монго
  ```
  ``` js
    Request
    Post http://example.org/user/login
    json{email:'',password:''}

    Response
    json{"_id","email","date_add","token"}
  ```

2. (не сделано)Создание магазина(магазинов может быть несколько) 

3. Аналитика по магазину
    1) Количество переходов по конкретном qr,url
    2) Данные по переходу(браузер,время)