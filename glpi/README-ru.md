Данный модуль для языка GO реализует механизм работы с HTTP API системы GLPI. Описание АПИ доступно по адресу https://github.com/glpi-project/glpi/blob/9.5/bugfixes/apirest.md, и также в установленной системе по ссылке https://glpi/apirest.php .

Ниже представлен вариант использования на примере работы с принтером. Т.е. в примере будет выполнено:
- Добавление производителя
- Добавление модели принтера
- Добавление домена
- Добавление сети (виртуальная сущность связывающаяя активы на сетевом уровне)
- Добавление принтера
- Добавление сетевого порта и привязка его к принтеру
- Добавление IP-адреса и привязка его к порту
- Добавление модели картриджа
- Добавление картриджа как еденицы учета
- "Установка" картриджа (привязка его к соответствующему принтеру)

Алгоритм работы с активами и другими сущностями GLPI ничем не отличается от данного примера.

# Краткое описание API

## Инициализация сессии

 Работа с API в GLPI производится в рамках сессии, которая определяется на основе кода доступа (токена). Данный код генерится и возвращается после положительной авторизации пользователя. Для этого требуется выполнить HTTP запрос к API методу "initSession" и передать в нём Application токен ("Настройки" -> "Общие" -> "API") и токен пользователя ("Администрирование" -> "Пользователи" -> Пользователь -> "Ключи Удаленного доступа") либо имя и пароль.

 ```curl -s -X GET -H "Content-Type: application/json" -H "Authorization: user_token ${GLPI_USER_TOKEN}" -H "App-Token: ${GLPI_APP_TOKEN}" "https://glpi/apirest.php/initSession"```

При положительном ответе система вернет:

```
{
  "session_token": "6ud8p3llbvl4qkf56t7klvhvr"
}
```

В дальнейшем в запросах передается токен приложения и токен сессии:

```curl -s -X GET -H "Content-Type: application/json" -H "Session-Token: ${GLPI_SESSION_TOKEN}" -H "App-Token: ${GLPI_APP_TOKEN}" "https://glpi/apirest.php/passivedcequipment"```

## Получение данных из GLPI

Для получения данных необходимо выполнить http-запрос к определенному методу API. Например для вывода информации о всех компьютерах ссылка будет выглядеть вот так:

```"https://glpi/apirest.php/computer"```

для получения данных по конкретному компьютеру:

```"https://glpi/apirest.php/computer/123"```


где "computer" - это тип актива (сущности, еденицы учета и т.п.), "123" - это уникальный идентификатор записи.

# Пример использования go-модуля glpi

## Импорт пакетов

``` go
package main
import (
	"encoding/json"
	"fmt"
	"glpi"
	"os"
)
```

## Установка переменных
Объявляем и инициализируем переменные:

``` go
var (
    glpiServerAPI string
    glpiUserToken string
    glpiAppToken  string
)
glpiServerAPI = os.Getenv("GLPI_SERVER_API")
glpiUserToken = os.Getenv("GLPI_USER_TOKEN")
glpiAppToken = os.Getenv("GLPI_APP_TOKEN")
```

* glpiServerAPI - адрес вызова API, в нашем случае берется из переменной окружения и равен https://glpi/apirest.php
* glpiUserToken - код доступа к API GLPI для конкретного пользователя
* glpiAppToken - код доступа к API GLPI

Для удобства работы, значения для данных переменных, беруться из переменных окружения.

Так как для начала работы необходимо получить код доступа для конкретной сессии, то произведем инициализацию сессиии:
``` go
glpiSession, err := glpi.NewSession(glpiServerAPI, glpiUserToken, glpiAppToken)
if err != nil {
    fmt.Println(err)
    return
}
```
glpi.NewSession - создает переменную типа Session и возвращает указатель на нее.

``` go
_, err = glpiSession.InitSession()
if err != nil {
    fmt.Println(err)
    return
}
```
glpiSession.InitSession() - инициализирует сессию, т.е. при вызове данной функции, происходит подключение к GLPI API и возвращается значение кода доступа (токена)сессии. Его значение сохраняется в поле структуры "glpiSession.sessionToken". Данный код можно посмотреть вызвав функцию "glpiSession.GetSessionToken()", т.е.

```fmt.Println(glpiSession.GetSessionToken())```

Установим значения некоторых переменных:

``` go
itemType := "Printer"
printerName := "acc-prn1"
printerDomain := "domain.local"
```

В GLPI API для поиска данных реализован метод search, вызываемый по ссылке http://glpi/apirest.php/search/itemtype/, где itemtype это тип еденицы учета либо другой сущности внутри GLPI (computer, monitor, printer, line и так далее). 

В модуле glpi для поиска используется функция "SearchItem(itemType string, itemName string)" вызываемая с параметрами:

- itemType - вышеозначенный тип
- itemName - имя.

В случае положительного результата функция возвращает цифровой уникальный идентификатор записи "ID".

Тут стоит остановиться подробнее на описании структуры и типах. В GO-модуле glpi каждая из сущностей представлена ввиде структуры с соответсвующим именем. Внтури структуры описано соответствие её полей и полей в итоговом JSON запросе. Например для производителей:

``` go
type GlpiManufacturer struct {
	Id      int    `json:"id"`
	Comment string `json:"comment"`
	Name    string `json:"name"`
}
```

Имена и типы полей go-структуры соответсвуют именам и типам таковых в JSON, кроме первой буквы (в json она маленькая). 

Для операций с записями в модуле реализованы функции с именами соответствующими glpi-типам (кроме простых типов с одинаковой струкрутрой). Функция на вход принимает тип запроса ("add" - добавление записи, "update" - редактирование, "delete" -удаление) и структуру соответствующего типа.

Т.е. для производителей это будет glpiSession.Manufacturer("add", manufacturersData) для принтеров glpiSession.Printer("add", printerData) и так далее.

Для простых типов введена структура:

``` go
type GlpiSimpleItem struct {
	Id      int    `json:"id"`
	Comment string `json:"comment"`
	Name    string `json:"name"`
}
```

Для моделей оборудования различного типа:

``` go
type GlpiItemModel struct {
	Id             int    `json:"id"`
	Comment        string `json:"comment"`
	Name           string `json:"name"`
	Product_number string `json:"product_number"`
}
```

И соответсвенно им и функции SimpleItem() и ItemModel().

В данном примере произведем поиск производителя по наименованию, и если записи с таким именем не найдено то добавим её.

``` go
// Get manufacturer ID or append new records/
manufacturersID := glpiSession.SearchItem("Manufacturer", "Systers")
if manufacturersID == 0 {
    manufacturersData := glpi.GlpiManufacturer{
        Name:    "Systers",
        Comment: "Holly Systers inc.",
    }
    manufacturersID = glpiSession.ItemOperation("add", "Manufacturer", manufacturersData)
}
```

Т.е. если manufacturersID будет равен 0 то инициализируем переменную типа GlpiManufacturer, заполним поля требуемыми значениями и вызовем нужную функцию. На выходе получим ID производителя для дальнейшего использования.

Найдем или добавим сеть (в данном контексте "сеть" - используется для объединения активов и построения схемы сети):

``` go
// Get network ID or append new records/
networksID := glpiSession.SearchItem("Network", "Head office")
if networksID == 0 {
    networksData := glpi.GlpiSimpleItem{
        Name:    "Head office",
        Comment: "Head office network",
    }
    networksID = glpiSession.ItemOperation("add", "Network", networksData)
}
```

Модель принтера:

``` go
// Search a used printer model, or append if don't find it.
printerModelID := glpiSession.SearchItem("PrinterModel", "Sys 1200")
if printerModelID == 0 {
    printerModelData := glpi.GlpiItemModel{
        Name: "Sys 1200",
    }
    printerModelID = glpiSession.ItemOperation("add", "PrinterModel", printerModelData)
}
````


Аналогичным образом поступаем и с остальными данными. Добавляем принтер с данными указанными в нижеприведенной структуре. В значениях полей указаны идентификаторы уже созданных записей, плюс переменные полученные выше (в принципе, для добавления принтера достаточно одного имени, остальное можно добавить позже):

``` go
printerData := glpi.GlpiPrinter{
    Name:             printerName,
    Locations_id:     1,
    Users_id_tech:    3,
    Users_id:         2,
    Groups_id_tech:   1,
    Groups_id:        2,
    Have_usb:         1,
    Networks_id:      networksID,
    Manufacturers_id: manufacturersID,
    Printermodels_id: printerModelID,
}

printerID := glpiSession.ItemOperation("add", "Printer", printerData)
fmt.Println("Printer was added, ID:", printerID)
```

Так как принтер подразумевается сетевой то добавим к нему сетевой порт и привяжем к нему IP-адрес. Данная функциональность подобным образом обрабатывается и для других активов GLPI, т.е. алгоритм работы тот-же меняются только itemtype.

Добавляем порт:

``` go
// Added Ethernet network port
networkPortData := glpi.GlpiNetworkPort{
    Name:               "Eth1",
    Items_id:           printerID,
    Itemtype:           itemType,
    Logical_number:     1,
    Instantiation_type: "NetworkPortEthernet",
}
glpiPortID := glpiSession.ItemOperation("add", "NetworkPort", networkPortData)
fmt.Println("Network port was added, ID: ", glpiPortID)
````

Добавляем домен:

``` go
// Getting domain ID, or add it
glpiFqdnID := glpiSession.SearchItem("FQDN", printerDomain)
if glpiFqdnID == 0 {
    fqdnData := glpi.GlpiFQDN{
        Name:    printerDomain,
        FQDN:    printerDomain,
        Comment: "Our local domain",
    }
    glpiFqdnID = glpiSession.ItemOperation("add", "fqdn", fqdnData)
}
fmt.Println("Domain - ", printerDomain, glpiFqdnID)
```

Связываем воедино сетевой порт, сетевое имя принтера и домен:

``` go
	networkData := glpi.GlpiNetworkName{
		Name:     printerName,
		Itemtype: "NetworkPort",
		Items_id: glpiPortID,
		Fqdns_id: glpiFqdnID,
		Comment:  "Network port",
	}
	glpiNetworkID := glpiSession.ItemOperation("add", "NetworkName", networkData)
	fmt.Println("Added network name for", printerName, "with ID:", glpiNetworkID)
```

Добавляем IP-адрес и привязываем его к принтеру и сетевому имени (порту):

``` go
ipData := glpi.GlpiIPAddress{
    Name:         "172.30.30.123",
    Itemtype:     "NetworkName",
    Items_id:     glpiNetworkID,
    Mainitems_id: printerID,
    Mainitemtype: itemType,
    Comment:      "Network port",
}
glpiIpAddressID := glpiSession.ItemOperation("add", "IPAddress", ipData)
fmt.Println("Added IP address:", glpiIpAddressID)
```

Т.е. схема добавления "Железка" -> "Порт" -> "Сетевое имя (сеть)" -> "IP-адрес" будет общая для всех активов или устройств имеющих сетевые порты.

На этом с принтером всё. Теперь можно добавить картриджей. Сперва добавляем модель картриджа

``` go
// ADD cartridge
requestDataCartridgeItem := glpi.GlpiCartridgeItem{
    // Id:   "11",
    Name: "С4092B",
    // Users_id_tech:  3,
    // Users_id:       2,
    // Groups_id_tech: 1,
    // Groups_id:      2,
    // Cartridgeitemtypes_id: 0,
    Manufacturers_id: manufacturersID,
    Alarm_threshold:  5,
}
cartridgeItemsID := glpiSession.ItemOperation("add", "CartridgeItem", requestDataCartridgeItem)
fmt.Println("Add a cartridge item:", cartridgeItemsID)
```

Затем указываем соответсвие (добавляем запись) между моделью принтера и моделью картриджа

``` go
// Linked printer model and cartridge
// printerModelID := glpiSession.SearchItem("PrinterModel", "Sys 1200")
requestData := glpi.GlpiCartridgeItemPrinterModel{
    // Id:   "11",
    Cartridgeitems_id: cartridgeItemsID,
    Printermodels_id:  printerModelID,
}

cartridgeItemsPrinterModel := glpiSession.ItemOperation("add", "cartridgeitem_printermodel", requestData)
fmt.Println("Added a link between Cartridge model and Printer model:", cartridgeItemsPrinterModel)
```

Затем добавляем картридж уже как еденицу хранения, привязывая к модели:

``` go
// Get the ID of an existing printer, or add it
// printerID = glpiSession.SearchItem("Printer", "acc-prn1")

requestDataCartridge := glpi.GlpiCartridge{
    // Id:   "11",
    Cartridgeitems_id: cartridgeItemsID,
}

cartridgeID := glpiSession.ItemOperation("add", "Cartridge", requestDataCartridge)
fmt.Println("Added a cartridge:", cartridgeID)
```

"Установим" картридж в принтер, т.е. привяжем конкретный картридж к конкретному принтеру:

``` go
// Getting a new cartridge info
res := glpiSession.GetItem("Cartridge", cartridgeID, "")
fmt.Println("Cartridge:", string(res))

var cartridgeInfo glpi.GlpiCartridge
err = json.Unmarshal(res, &cartridgeInfo)

if err != nil {
    fmt.Println(err)
}

// Installing the cartridge in the printer
cartridgeInfo.Printers_id = printerID
cartridgeInfo.Date_use = "2021-08-20"

res1 := glpiSession.ItemOperation("update", "Cartridge", cartridgeInfo)
fmt.Println("Add a link between cartridge and printer:", res1)
```

Перед обновлением конкретной записи требуется выбрать все данные и сохранить в переменную соответствующего типа, если этого не сделать то при выполнении операции "update" не инициализированные поля будут пустыми (т.е. данные затрутся).

Для получения информации о конкретном оборудовании в модуле есть функция ```GetItem(itemType string, itemID int, otherParam string)```. Она принимает на вход тип оборудования и его ID, возвращает JSON ввиде массива байт ([]byte).  Так как функция GetItem() возвращает JSON, то его можно увидеть просто преобразовав байты в строку при помощи string()

Получим информацию о добавленном картридже и выведем в косоль:

``` go
// Getting a new cartridge info
res := glpiSession.GetItem("Cartridge", cartridgeID, "")
fmt.Println("Cartridge:", string(res))
```

Для преобразования данных из json в структуру требуемого нам типа применим функцию json.Unmarshal():

``` go
var cartridgeInfo glpi.GlpiCartridge
err = json.Unmarshal(res, &cartridgeInfo)
if err != nil {
    fmt.Println(err)
}
fmt.Println("Cartridge:", cartridgeInfo)
```

Доступ к полям структуры можно получить штатным образом:

- cartridgeInfo.Id - уникальный идентификатор
- cartridgeInfo.Printers_id - идентификатор принтера
- cartridgeInfo.Date_in - дата установкии
- и т.д.

После выполнения всех операция удаляем нашу сессию. GLPI периодически сессии очищает сам, но лучше удалить:

``` go
_, err = glpiSession.KillSession()
if err != nil {
    fmt.Println(err)
    return
}
```