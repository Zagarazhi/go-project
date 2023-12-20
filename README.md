# Описание  
Этот проект представляет собой два сервиса, соединенных с помощью gRPC, которые обращаются к API Binance  

# Запуск  
Для старта заполните [конфигурационный файл](config.yaml).  
Далее в папке проекта: go run main.go  

# Путь  
GET "serverAdress:serverPort/process" вместе с JSON.  
Пример:  
{  
    "symbol": "ETHBTC",  
    "interval": 1,  
}  
symbol - пары монет  
interval - тип интервала, где  
- "INTERVAL_1S":     0,
- "INTERVAL_1M":     1,
- "INTERVAL_3M":     2,
- "INTERVAL_5M":     3,
- "INTERVAL_15M":    4,
- "INTERVAL_30M":    5,
- "INTERVAL_1H":     6,
- "INTERVAL_2H":     7,
- "INTERVAL_4H":     8,
- "INTERVAL_6H":     9,
- "INTERVAL_8H":     10,
- "INTERVAL_12H":    11,
- "INTERVAL_1D":     12,
- "INTERVAL_3D":     13,
- "INTERVAL_1W":     14,
- "INTERVAL_1MONTH": 15,

# Примеры работы  
[Здесь](examples) можно увидеть скриншоты выполнения запросов