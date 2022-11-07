![Gopher](https://user-images.githubusercontent.com/29982395/200129257-71a4c973-78ad-43e6-a1b3-353427fee487.jpeg)
# Money Manager

Микросервис для работы с балансом пользователей


## Содержание
- [Запуск](#запуск)
- [Предоставляемый API](#предоставляемый-api)
- [Функционал](#функционал)
- [Вопросы к ТЗ](#вопросы-к-тз-и-их-решения)

## Запуск
- Склонировать репозиторий
```sh
git clone git@github.com:Gusyatnikova/money-manager.git   
```
- (Опционально) изменить переменные окружения в файле [`.env`](.env)
- Запустить сервис
```sh
docker-compose -f money-manager/docker-compose.yml up
```
- Проверить соединение (http://<HTTP_HOST>:<HTTP_PORT>/healthcheck)
```sh
curl http://0.0.0.0:8888/healthcheck
```
- При необходимости сохранять данные между запусками запустить сервис командой
```sh
docker-compose -f money-manager/docker-compose-persistent.yml up
```

## Предоставляемый API
- Документация в [Postman](https://documenter.getpostman.com/view/17987701/2s8YYBRmmE)

## Функционал
- [x] Зачисление денег на баланс пользователя
- [x] Получение баланса пользователя
- [x] Списание денег с баланса пользователя
- [x] Перевод денег от пользователя к пользователю
- [x] Резервирование денег с основного баланса на отдельном счету
- [x] Признание выручки из зарезервированных денег с занесением записи в отчет для бухгалтерии
- [x] Разрезервирование денег, если услугу применить не удалось
- [x] Формирование сводного отчета по всем пользователям в разрезе каждой услуги

## Вопросы к [ТЗ](https://github.com/avito-tech/internship_backend_2022) и их решения
1. С какими денежными единицами работать?
    * Внутри микросервиса все расчеты для точности вычислений производятся в копейках
    * Для удобства использования API поддерживает работу как с рублями, так и с копейками.
  ```sh
    "money": {
        "amount": "101.57", //также допустимо: 101.1, 101, 101.00, значение 0 недопустимо
        "unit": "rub"
    }
```  
```sh
    "money": {
        "amount": "57", //только целое положительное число, отличное от 0
        "unit": "kop"
    }
```
2. Как представлен id пользователя и как его проверять?
   * Т.к. формат хранения id пользователя может быть разным (число, строка, uuid, ulid и т.д.), выбран самый универсальный формат - строка.
   * Считается, что есть мастер-система, работающая с пользователями, к ней и нужно обращаться при проверках. Т.к. это сторонний сервис, проверка выполняется только на пустоту
3. Возможно ли, что запрос на резервирование(признание) средств по одной тройке (id пользователя, id услуги, id заказа)
   будет выслан в систему несколько раз?
   * 

5. Насколько много будет пользователей в системе?
   * Атомарность операций с базой данных в моей реализации достигается за счет использования транзакций. Если предположить, что 
   пользователей и операций станет настолько много, что база потребует шардирования, поддерживать транзакции в представленном виде станет невозможно.
   * Возможным решением такой ситуации станут двухфазные комиты, но ввиду ограниченный по времени и моим компетенциям, есть только транзакции =)
##