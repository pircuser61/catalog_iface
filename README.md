Сервис интерфейс для работы с gitlab.ozon.dev/pircuser61/catalog
Осуществляет валидацию данных и передает их дальше

docker-comoose - для jaeger, остальное в gitlab.ozon.dev/pircuser61/catalog

Jaeger http://localhost:16686
Swagger http://localhost:8080/swagger/

curl -X POST localhost:8080/v1/good -d '{"name":"name1", "unit_of_measure":"uom1", "country":"country1"}'
curl localhost:8080/v1/goods
curl localhost:8080/v1/good/1
curl -X PUT localhost:8080/v1/good -d '{"good":{"code":1, "name":"name2", "unit_of_measure":"uom1", "country":"Country1"}}'
curl -X DELETE localhost:8080/v1/good -d '{"code":1}'

GRPC: localhost:8081
Команды:
addCountry country3
listCountry
getCountry 1
updateCountry 1 country1
deleteCountry 1

    addUom uom1
    listUom
    getUom 1
    updateUom 1 Uom1
    deleteUom 1

    add name3 uom3 country3
    list
    get 4
    update 4 name4 uom4 country4
    delete 4
