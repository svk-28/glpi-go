package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"glpi"
	"os"
	"strings"
	"zabbix"
)

// Find and return a single host object by name
func GetHost(zabbixAPI *zabbix.API, host string) (zabbix.ZabbixHost, error) {
	params := make(map[string]interface{}, 0)
	filter := make(map[string]string, 0)
	filter["host"] = host
	params["filter"] = filter
	params["output"] = "extend"
	params["select_groups"] = "extend"
	params["templated_hosts"] = 1
	params["selectInventory"] = "extend"
	params["selectInterfaces"] = "extend"
	ret, err := zabbixAPI.Host("get", params)

	// This happens if there was an RPC error
	if err != nil {
		return nil, err
	}

	// If our call was successful
	if len(ret) > 0 {
		return ret[0], err
	}

	// This will be the case if the RPC call was successful, but
	// Zabbix had an issue with the data we passed.
	return nil, &zabbix.ZabbixError{0, "", "Host not found"}
}

func GetHostGroup(zabbixAPI *zabbix.API, hostgroup string, select_type string) (zabbix.ZabbixHostGroup, error) {
	params := make(map[string]interface{}, 0)
	filter := make(map[string]string, 0)
	filter["name"] = hostgroup
	params["filter"] = filter
	params["output"] = "extend"
	if select_type == "all_hosts" {
		params["selectHosts"] = "extend"
	}
	// params["select_groups"] = "extend"
	//params["templated_hosts"] = 1
	ret, err := zabbixAPI.HostGroup("get", params)

	if err != nil {
		return nil, err
	}

	// If our call was successful
	if len(ret) > 0 {
		return ret[0], err
	}

	// This will be the case if the RPC call was successful, but
	// Zabbix had an issue with the data we passed.
	return nil, &zabbix.ZabbixError{0, "", "HostGroup not found"}
}

func MigrateHost(glpiAPI *glpi.Session, zabbixAPI *zabbix.API, item interface{}, glpiLocationID int) {
	var (
		zbxHostOS            string
		zbxHostIP            string
		zbxHostFQDN          string
		zbxHostSerial        string
		zbxHostOtherSerial   string
		zbxHostUseDate       string
		zbxHostManufacturer  string
		zbxHostComputerModel string
		glpiItemStatusID     int
	)
	itemType := "Computer"
	typeID := glpiAPI.SearchItem("ComputerType", "Сервер")
	// location_id := glpiAPI.SearchItem("Location", "ГО")
	fmt.Println(typeID)
	// Получаем имя узла, id и флаг инвентаризации
	zbxHostName := fmt.Sprint(item.(map[string]interface{})["host"])

	zbxHostStatus := fmt.Sprint(item.(map[string]interface{})["status"])
	// // если "inventory_mode": "-1" то инвентаризация отключена
	zbxHostHasInventory := fmt.Sprint(item.(map[string]interface{})["inventory_mode"])
	// Если имя узла в заббикс как FQDN то обрезаем до первой точки
	// zbxHostName_short := strings.Split(zbxHostName, ".")[0]

	glpiCompID := glpiAPI.SearchItem("Computer", strings.ToLower(zbxHostName))

	fmt.Println(">>>>>>> ", zbxHostName)

	if glpiCompID != 0 {
		fmt.Println("Узел ", zbxHostName, " существует в glpi c id:", glpiCompID)
		return
	}
	// По имени узла запрашиваем полную информацию из заббикса и выдергиваем данные
	zbxHostInfo, _ := GetHost(zabbixAPI, fmt.Sprint(zbxHostName))
	for _, i := range zbxHostInfo["interfaces"].([]interface{}) {
		zbxHostIP = fmt.Sprint(i.(map[string]interface{})["ip"])
		zbxHostFQDN = strings.ToLower(fmt.Sprint(i.(map[string]interface{})["dns"]))
		fmt.Println(zbxHostIP, zbxHostFQDN)
	}
	// Проверяем если флаг инвентаризации не равен -1 то получаем инвентарные данные
	if zbxHostHasInventory != "-1" {
		zbxHostOS = fmt.Sprint(zbxHostInfo["inventory"].(map[string]interface{})["os"])
		zbxHostSerial = fmt.Sprint(zbxHostInfo["inventory"].(map[string]interface{})["serialno_a"])
		zbxHostOtherSerial = fmt.Sprint(zbxHostInfo["inventory"].(map[string]interface{})["serialno_b"])
		zbxHostUseDate = fmt.Sprint(zbxHostInfo["inventory"].(map[string]interface{})["date_hw_install"])
		zbxHostManufacturer = fmt.Sprint(zbxHostInfo["inventory"].(map[string]interface{})["vendor"])
		zbxHostComputerModel = fmt.Sprint(zbxHostInfo["inventory"].(map[string]interface{})["model"])

		fmt.Println(zbxHostOS, zbxHostSerial, zbxHostOtherSerial, zbxHostUseDate, zbxHostManufacturer, zbxHostComputerModel)

	}

	// получаем id модели оборудования. Если не найдено добавляем в glpi
	glpiModelID := glpiAPI.SearchItem("ComputerModel", zbxHostComputerModel)
	if glpiModelID == 0 {
		requestData := glpi.GlpiComputerModel{
			Name: zbxHostComputerModel,
		}
		glpiModelID = glpiAPI.ItemOperation("add", "ComputerModel", requestData)
	}

	// получаем id производителя из glpiAPI. если нету добавляем
	glpiManufacturerID := glpiAPI.SearchItem("Manufacturer", zbxHostManufacturer)
	if glpiManufacturerID == 0 {
		requestData := glpi.GlpiManufacturer{
			Name: zbxHostManufacturer,
		}
		glpiManufacturerID = glpiAPI.ItemOperation("add", "Manufacturer", requestData)
	}

	// добавляем железку
	// проверяем статус
	if zbxHostStatus == "0" {
		glpiItemStatusID = glpiAPI.SearchItem("State", "В работе")
	}
	if zbxHostStatus == "1" {
		glpiItemStatusID = glpiAPI.SearchItem("State", "Выключено")
	}

	requestDataComp := glpi.GlpiComputer{
		Name:              strings.ToLower(zbxHostName),
		Serial:            zbxHostSerial,
		Otherserial:       zbxHostOtherSerial,
		Locations_id:      glpiLocationID,
		Computermodels_id: glpiModelID,
		Computertypes_id:  typeID,
		Manufacturers_id:  glpiManufacturerID,
		States_id:         glpiItemStatusID,
	}
	glpiCompID = glpiAPI.ItemOperation("add", "Computer", requestDataComp)
	fmt.Println("Добавлен хост", glpiCompID)

	// Если выставлена дата установки то включаем фин. информаци.
	// к году прибвавляем месяц и день
	if len(zbxHostUseDate) <= 4 {
		zbxHostUseDate = zbxHostUseDate + "-01-01"
	}
	if zbxHostUseDate != "" {
		requestDataInfocom := glpi.GlpiInfocom{
			Items_id: glpiCompID,
			Itemtype: itemType,
			Use_date: zbxHostUseDate,
		}
		// fmt.Println("--->", requestDataInfocom)
		glpiAPI.ItemOperation("add", "InfoCom", requestDataInfocom)
		fmt.Println("Добавлена фин. информация ", glpiCompID)
	}

	// Добавляем сетевой порт, адрес, FQDN, e.t.c
	if zbxHostIP != "" {
		glpiNetworkPortData := glpi.GlpiNetworkPort{
			Items_id:           glpiCompID,
			Itemtype:           itemType,
			Name:               "Eth",
			Instantiation_type: "NetworkPortEthernet",
		}
		glpiPortID := glpiAPI.ItemOperation("add", "NetworkPort", glpiNetworkPortData)
		fmt.Println("Добавлен port", glpiPortID)

		// Отделяем имя узла от домена
		zbxHostFQDN_short := strings.Split(zbxHostFQDN, ".")[0]
		prefix := zbxHostFQDN_short + "."
		zbxHostDomain := strings.TrimLeft(zbxHostFQDN, prefix)
		fmt.Println(zbxHostFQDN, zbxHostFQDN_short, zbxHostDomain)
		// Определяем ID домена в GLPI если он там есть
		glpiFqdnID := glpiAPI.SearchItem("FQDN", zbxHostDomain)
		fmt.Println("Домен - ", zbxHostDomain, glpiFqdnID)

		glpiNetworkData := glpi.GlpiNetworkName{
			Items_id: glpiPortID,
			Itemtype: "NetworkPort",
			Name:     zbxHostFQDN_short,
			Fqdns_id: glpiFqdnID,
		}
		glpiNetworkID := glpiAPI.ItemOperation("add", "NetworkName", glpiNetworkData)
		fmt.Println("Добавлена сеть", glpiNetworkID)

		glpiIPAdressData := glpi.GlpiIPAddress{
			Items_id: glpiNetworkID,
			Itemtype: "NetworkName",
			Name:     zbxHostIP,
		}
		glpiIpAddressID := glpiAPI.ItemOperation("add", "IPAddress", glpiIPAdressData)
		fmt.Println("Добавлен IP", glpiIpAddressID)
	}
	// Проверяем что в заббикс для узла есть запись об ОС
	if zbxHostOS != "" {
		// получаем id ОС из glpiAPI. если нет добавляем
		glpiOSID := glpiAPI.SearchItem("OperatingSystem", zbxHostOS)
		fmt.Println("Добавлена ОС", glpiOSID)
		if glpiOSID == 0 {
			// request_data_os := glpi.GlpiOperatingSystem{
			// 	Name: zbxHostOS,
			// }
			glpiOSID = glpiAPI.ItemOperation("add", "OperatingSystem", glpi.GlpiOperatingSystem{Name: zbxHostOS})
			fmt.Println("Добавлена ОС", glpiOSID)
		}
		// Привязываем ОС к узлу в GLPI
		requestDataOSItem := glpi.GlpiOSItem{
			Items_id:            glpiCompID,
			Itemtype:            itemType,
			Operatingsystems_id: glpiOSID,
		}
		glpiAPI.ItemOperation("add", "Item_OperatingSystem", requestDataOSItem)
	}

}

func UpdateHostStatus(glpiAPI *glpi.Session, zabbixAPI *zabbix.API, item interface{}) {
	// Получаем имя узла, id и флаг инвентаризации
	itemType := "Computer"
	// typeID := glpiAPI.SearchItem("ComputerType", "Сервер")

	zbxHostName := fmt.Sprint(item.(map[string]interface{})["host"])

	zbxHostStatus := fmt.Sprint(item.(map[string]interface{})["status"])
	glpiItemID := glpiAPI.SearchItem(itemType, strings.ToLower(zbxHostName))
	res := glpiAPI.GetItem(itemType, glpiItemID, "")

	var statusInfo map[string]interface{}
	err := json.Unmarshal(res, &statusInfo)

	if err != nil {
		fmt.Println(err)
	}

	glpiItemStatusID := statusInfo["states_id"]

	fmt.Println(">>>>>>> ", zbxHostName, zbxHostStatus, "glp id:", glpiItemID, "glpi status", glpiItemStatusID)

	var glpiStatusID int

	if zbxHostStatus == "0" {
		glpiStatusID = glpiAPI.SearchItem("State", "В работе")
	}
	if zbxHostStatus == "1" {
		glpiStatusID = glpiAPI.SearchItem("State", "Выключено")
	}
	fmt.Println(glpiItemStatusID, glpiStatusID)
	if glpiItemStatusID != glpiStatusID {
		fmt.Println("Update")
		glpiAPI.UpdateItemStatus(itemType, glpiItemID, glpiStatusID)
	}
}

func main() {
	var (
		zabbixServerAPI  string
		zabbixUser       string
		zabbixPassword   string
		zabbixHost       string
		zabbixHostGroups string
		operation        string
		outputFormat     string
		glpiServerAPI    string
		glpiUserToken    string
		glpiAppToken     string
		glpiLocationID   int
	)

	flag.StringVar(&zabbixServerAPI, "zabbix-server-api", "", "zabbix server API URL. ZABBIX_SERVER_API")
	flag.StringVar(&zabbixUser, "zabbix-user", "", "Zabbix user name. ZABBIX_USER")
	flag.StringVar(&zabbixPassword, "zabbix-password", "", "Zabbix password. ZABBIX_PASSWORD")
	flag.StringVar(&zabbixHost, "zabbix-host", "", "Zabbix host name")
	flag.StringVar(&zabbixHostGroups, "zabbix-host-groups", "", "Zabbix host groups, can be comma-separated list")

	flag.StringVar(&glpiServerAPI, "glpi-server-api", "", "Glpi instance API URL. GLPI_SERVER_URL")
	flag.StringVar(&glpiUserToken, "glpi-user-token", "", "Glpi user API token. GLPI_USER_TOKEN")
	flag.StringVar(&glpiAppToken, "glpi-app-token", "", "Glpi aplication API token. GLPI_APP_TOKEN")
	flag.IntVar(&glpiLocationID, "glpi-location-id", 0, "Glpi location id for migration hosts")

	flag.StringVar(&operation, "operation", "zabbix-version", "Opertation type, must be:\n\tzabbix-version\n\tglpi-version\n\tget-host - getting zabbix host information\n\tget-hostgroup - getting zabbix hosts group\n\tget-hosts-from-group - getting zabbix group members hosts\n\tmigrate-host - migrate host from Zabbix to GLPI\n\tmigrate-hosts-from-group - migrate all zabbix groups members hosts to GLPI\n\tupdate-host-status - get host status from zabbix and update them into zabbix\n\tupdate-hosts-status-from-group - get group members hosts status from zabbix and update into glpi\n\t")

	flag.StringVar(&outputFormat, "output-format", "go", "Output format: 'go' - go map, 'json' - json string")

	flag.Parse()

	if zabbixServerAPI == "" && os.Getenv("ZABBIX_SERVER_API") == "" {
		fmt.Println("Send error: make sure environment variables `ZABBIX_SERVER_API`, or used with '-zabbix-server-api' argument")
		os.Exit(1)
	} else if zabbixServerAPI == "" && os.Getenv("ZABBIX_SERVER_API") != "" {
		zabbixServerAPI = os.Getenv("ZABBIX_SERVER_API")
	}
	// fmt.Println(zabbixServerAPI)

	if zabbixUser == "" && os.Getenv("ZABBIX_USER") == "" {
		fmt.Println("Send error: make sure environment variables `ZABBIX_USER`, or used with '-zabbix-user' argument")
		os.Exit(1)
	} else if zabbixUser == "" && os.Getenv("ZABBIX_USER") != "" {
		zabbixUser = os.Getenv("ZABBIX_USER")
	}
	// fmt.Println(zabbixUser	)

	if zabbixPassword == "" && os.Getenv("ZABBIX_PASSWORD") == "" {
		fmt.Println("Send error: make sure environment variables `ZABBIX_PASSWORD`, or used with '-zabbix-password' argument")
		os.Exit(1)
	} else if zabbixPassword == "" && os.Getenv("ZABBIX_PASSWORD") != "" {
		zabbixPassword = os.Getenv("ZABBIX_PASSWORD")
	}
	// fmt.Println(zabbixPassword)

	//-----------------------------------------------------------
	if glpiServerAPI == "" && os.Getenv("GLPI_SERVER_API") == "" {
		fmt.Println("Send error: make sure environment variables `GLPI_SERVER_API`, or used with '-glpi-server-api' argument")
		os.Exit(1)
	} else if glpiServerAPI == "" && os.Getenv("GLPI_SERVER_API") != "" {
		glpiServerAPI = os.Getenv("GLPI_SERVER_API")
	}

	if glpiUserToken == "" && os.Getenv("GLPI_USER_TOKEN") == "" {
		fmt.Println("Send error: make sure environment variables `GLPI_USER_TOKEN`, or used with '-glpi-user-token' argument")
		os.Exit(1)
	} else if glpiUserToken == "" && os.Getenv("GLPI_USER_TOKEN") != "" {
		glpiUserToken = os.Getenv("GLPI_USER_TOKEN")
	}

	if glpiAppToken == "" && os.Getenv("GLPI_APP_TOKEN") == "" {
		fmt.Println("Send error: make sure environment variables `GLPI_APP_TOKEN`, or used with '-glpi-app-token' argument")
		os.Exit(1)
	} else if glpiAppToken == "" && os.Getenv("GLPI_APP_TOKEN") != "" {
		glpiAppToken = os.Getenv("GLPI_APP_TOKEN")
	}

	// -------  Connected to zabbix API ----------
	zabbixAPI, err := zabbix.NewAPI(zabbixServerAPI, zabbixUser, zabbixPassword)
	if err != nil {
		fmt.Println(err)
		return
	}

	if operation != "zabbix-version" {
		_, err = zabbixAPI.Login()
		if err != nil {
			fmt.Println(err)
			return
		}
		// fmt.Println("Connected to API: OK")
	}

	// ------- Connected to GLPI API -----------
	glpiAPI, err := glpi.NewSession(glpiServerAPI, glpiUserToken, glpiAppToken)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = glpiAPI.InitSession()
	if err != nil {
		fmt.Println(err)
		return
	}

	switch operation {
	case "zabbix-version":
		zabbixVersion, err := zabbixAPI.Version()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Zabbix version:", zabbixVersion)
	case "glpi-version":
		glpiVersion := glpiAPI.Version()
		fmt.Println("GLPI version:", glpiVersion)
	case "get-host":
		if zabbixHost == "" {
			fmt.Println("Zabbix host not setting, using '-zabbix-host' command line argument")
			os.Exit(1)
		} else {
			res, err := GetHost(zabbixAPI, zabbixHost)
			if err != nil {
				fmt.Println(err)
			}
			if outputFormat == "json" {
				result, err := json.Marshal(res)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(string(result))
			} else {
				fmt.Println(res)
			}
		}
	case "get-hostgroup":
		if zabbixHostGroups == "" {
			fmt.Println("Zabbix host group not setting, using '-zabbix-host-groups' command line argument")
			os.Exit(1)
		} else {
			groupList := strings.Split(zabbixHostGroups, ",")
			for _, item := range groupList {
				fmt.Println(GetHostGroup(zabbixAPI, strings.TrimSpace(item), "no_host"))
			}
		}
	case "get-hosts-from-group":
		if zabbixHostGroups == "" {
			fmt.Println("Zabbix host group not setting, using '-zabbix-host-groups' command line argument and comma separated list of groups into qoute")
			os.Exit(1)
		} else {
			groupList := strings.Split(zabbixHostGroups, ",")
			for _, item := range groupList {
				res, err1 := GetHostGroup(zabbixAPI, strings.TrimSpace(item), "all_hosts")
				// println(strings.TrimSpace(item))
				if err1 != nil {
					fmt.Println(err1)
				}
				if outputFormat == "json" {
					result, err := json.Marshal(res)
					if err != nil {
						fmt.Println(err)
					}
					// если "inventory_mode": "-1" то инвентаризация отключена
					fmt.Println(string(result))
				} else {
					fmt.Println(res)
				}
			}
		}
	case "migrate-host":
		fmt.Println("Function not implemented yet")
	case "migrate-hosts-from-group":
		if zabbixHostGroups == "" {
			fmt.Println("Zabbix host group not setting, using '-zabbix-host-groups' command line argument and comma separated list of groups into qoute")
			os.Exit(1)
		} else {
			groupList := strings.Split(zabbixHostGroups, ",")
			for _, item := range groupList {
				res, err1 := GetHostGroup(zabbixAPI, strings.TrimSpace(item), "all_hosts")
				println(strings.TrimSpace(item))
				if err1 != nil {
					fmt.Println(err1)
				}
				// Пробегаемся по массиву узлов из забикса и выдергиваем интерфейсы и инвентарные данные
				for _, item := range res["hosts"].([]interface{}) {
					MigrateHost(glpiAPI, zabbixAPI, item, glpiLocationID)
				}
			}
		}
	case "update-host-status":
		fmt.Println("Function not implemented yet")
	case "update-hosts-status-from-group":
		if zabbixHostGroups == "" {
			fmt.Println("Zabbix host group not setting, using '-zabbix-host-groups' command line argument and comma separated list of groups into qoute")
			os.Exit(1)
		} else {
			groupList := strings.Split(zabbixHostGroups, ",")
			for _, item := range groupList {
				res, err1 := GetHostGroup(zabbixAPI, strings.TrimSpace(item), "all_hosts")
				println(strings.TrimSpace(item))
				if err1 != nil {
					fmt.Println(err1)
				}
				// Пробегаемся по массиву узлов из забикса и выдергиваем статус
				for _, item := range res["hosts"].([]interface{}) {
					fmt.Println(item)
					UpdateHostStatus(glpiAPI, zabbixAPI, item)
				}
			}
		}

	}

}
