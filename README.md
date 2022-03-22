# :heavy_check_mark: MailingService
**MailingService** - сервис управления рассылками API администрирования и получения статистики.

# Подробная документация по API 
https://app.swaggerhub.com/apis-docs/DeedsBaron/MailingListService/1.0.0
# Конфиг  
![image](https://user-images.githubusercontent.com/80648065/159502935-55dd39ad-91ce-4aff-9aae-0d7048d8b056.png)
# Флаги
![image](https://user-images.githubusercontent.com/80648065/159503293-6eea7882-c92e-4cba-8d9a-4ffd5e179107.png)
# Схема базы данных
![image](https://user-images.githubusercontent.com/80648065/159508544-f96f401e-4c8f-4db1-86c7-3ed113ce4ccd.png)
# Usage
По умолчанию поднимается контейнер в котором работает сервис

    make

Поднимается контейнер с сервисом и контейнер с базой данных postgresql

    make containers
    
Выполняются тесты, обязательно должен быть поднят контейнер с базой

    make test
# Дополнительные задания
1. Организовать тестирование написанного кода
3. Подготовить docker-compose для запуска всех сервисов проекта одной командой
5. Сделать так, чтобы по адресу /docs/ открывалась страница со Swagger UI и в нём отображалось описание разработанного API.

    
## Other
**Author:**  
:vampire:*[Deeds Baron](https://github.com/DeedsBaron)*  

