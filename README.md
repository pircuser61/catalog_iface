Сервис интерфейс для работы с gitlab.ozon.dev/pircuser61/catalog
Осуществляет валидацию данных и передает их дальше

Swagger http://localhost:8080/swagger/

curl -X POST localhost:8080/v1/good -d '{"name":"name1", "unit_of_measure":"uom1", "country":"country1"}'
curl localhost:8080/v1/goods
curl localhost:8080/v1/good/1
curl -X PUT localhost:8080/v1/good -d '{"good":{"code":1, "name":"name2", "unit_of_measure":"uom1", "country":"Country1"}}'
curl -X DELETE localhost:8080/v1/good -d '{"code":1}'

GRPC: localhost:8081
Команды:
    add name3 uom3 country3
    list
    get 4
    update 4 name4 uom4 country4
    delete 4