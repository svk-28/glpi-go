package glpi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	url          string
	userToken    string
	appToken     string
	sessionToken string
}

type GlpiComputer struct {
	Id                   int    `json:"id"`
	Autoupdatesystems_id int    `json:"autoupdatesystems_id"`
	Computermodels_id    int    `json:"computermodels_id"`
	Computertypes_id     int    `json:"computertypes_id"`
	Comment              string `json:"comment"`
	Contact              string `json:"contact"`
	Contact_num          string `json:"contact_num"`
	Entities_id          int    `json:"entities_id"`
	Groups_id            int    `json:"groups_id"`
	Groups_id_tech       int    `json:"groups_id_tech"`
	Is_template          int    `json:"is_template"`
	Is_deleted           int    `json:"is_deleted"`
	Locations_id         int    `json:"locations_id"`
	Manufacturers_id     int    `json:"manufacturers_id"`
	Name                 string `json:"name"`
	Networks_id          int    `json:"networks_id"`
	Otherserial          string `json:"otherserial"`
	Serial               string `json:"serial"`
	States_id            int    `json:"states_id"`
	Users_id             int    `json:"users_id"`
	Users_id_tech        int    `json:"users_id_tech"`
	UUID                 string `json:"uuid"`
	Template_name        string `json:"template_name"`
}

type GlpiMonitor struct {
	Id               int    `json:"id"`
	Monitormodels_id int    `json:"monitormodels_id"`
	Monitortypes_id  int    `json:"monitortypes_id"`
	Comment          string `json:"comment"`
	Contact          string `json:"contact"`
	Contact_num      string `json:"contact_num"`
	Entities_id      int    `json:"entities_id"`
	Groups_id        int    `json:"groups_id"`
	Groups_id_tech   int    `json:"groups_id_tech"`
	Is_global        int    `json:"is_global"`
	Is_template      int    `json:"is_template"`
	Is_deleted       int    `json:"is_deleted"`
	Have_micro       int    `json:"have_micro"`
	Have_speaker     int    `json:"have_speaker"`
	Have_subd        int    `json:"have_subd"`
	Have_bnc         int    `json:"have_bnc"`
	Have_dvi         int    `json:"have_dvi"`
	Have_pivot       int    `json:"have_pivot"`
	Have_hdmi        int    `json:"have_hdmi"`
	Have_displayport int    `json:"have_displayport"`
	Locations_id     int    `json:"locations_id"`
	Manufacturers_id int    `json:"manufacturers_id"`
	Name             string `json:"name"`
	Otherserial      string `json:"otherserial"`
	Serial           string `json:"serial"`
	Size             string `json:"size"`
	States_id        int    `json:"states_id"`
	Template_name    string `json:"template_name"`
	Users_id         int    `json:"users_id"`
	Users_id_tech    int    `json:"users_id_tech"`
}

type GlpiSoftware struct {
	Id                    int    `json:"id"`
	Comment               string `json:"comment"`
	Entities_id           int    `json:"entities_id"`
	Is_helpdesk_visible   int    `json:"is_helpdesk_visible"`
	Is_update             int    `json:"is_update"`
	Is_template           int    `json:"is_template"`
	Is_valid              int    `json:"is_valid"`
	Is_deleted            int    `json:"is_deleted"`
	Groups_id             int    `json:"groups_id"`
	Groups_id_tech        int    `json:"groups_id_tech"`
	Locations_id          int    `json:"locations_id"`
	Manufacturers_id      int    `json:"manufacturers_id"`
	Name                  string `json:"name"`
	Softwares_id          int    `json:"softwares_id"`
	Softwarecategories_id int    `json:"softwarecategories_id"`
	Template_name         string `json:"template_name"`
	Users_id              int    `json:"users_id"`
	Users_id_tech         int    `json:"users_id_tech"`
}

type GlpiOperatingSystem struct {
	Id      int    `json:"id"`
	Comment string `json:"comment"`
	Name    string `json:"name"`
}

type GlpiManufacturer struct {
	Id      int    `json:"id"`
	Comment string `json:"comment"`
	Name    string `json:"name"`
}

type GlpiComputerModel struct {
	Id                int    `json:"id"`
	Comment           string `json:"comment"`
	Name              string `json:"name"`
	Product_number    string `json:"product_number"`
	Weight            int    `json:"weight"`
	Required_units    int    `json:"required_units"`
	Depth             int    `json:"depth"`
	Power_connections int    `json:"power_connections"`
	Power_consumption int    `json:"power_consumption"`
	Is_half_rack      int    `json:"is_half_rack"`
	Picture_front     string `json:"picture_front"`
	Picture_rear      string `json:"picture_rear"`
}

type GlpiNetworkEquipment struct {
	Id                        int    `json:"id"`
	Networkequipmentmodels_id int    `json:"networkequipmentmodels_id"`
	Networkequipmenttypes_id  int    `json:"networkequipmenttypes_id"`
	Comment                   string `json:"comment"`
	Contact                   string `json:"contact"`
	Contact_num               string `json:"contact_num"`
	Entities_id               int    `json:"entities_id"`
	Groups_id                 int    `json:"groups_id"`
	Groups_id_tech            int    `json:"groups_id_tech"`
	Is_template               int    `json:"is_template"`
	Is_deleted                int    `json:"is_deleted"`
	Locations_id              int    `json:"locations_id"`
	Manufacturers_id          int    `json:"manufacturers_id"`
	Name                      string `json:"name"`
	Networks_id               int    `json:"networks_id"`
	Otherserial               string `json:"otherserial"`
	Ram                       string `json:"ram"`
	Serial                    string `json:"serial"`
	States_id                 int    `json:"states_id"`
	Template_name             string `json:"template_name"`
	Users_id                  int    `json:"users_id"`
	Users_id_tech             int    `json:"users_id_tech"`
}

type GlpiInfocom struct {
	Id                     int    `json:"id"`
	Bill                   string `json:"bill"`
	Budgets_id             int    `json:"budgets_id"`
	Businesscriticities_id int    `json:"businesscriticities_id"`
	Buy_date               string `json:"buy_date"`
	Comment                string `json:"comment"`
	Decommission_date      string `json:"decommission_date"`
	Delivery_date          string `json:"delivery_date"`
	Delivery_number        string `json:"delivery_number"`
	Entities_id            int    `json:"entities_id"`
	Immo_number            string `json:"immo_number"`
	Inventory_date         string `json:"inventory_date"`
	Items_id               int    `json:"items_id"`
	Itemtype               string `json:"itemtype"`
	Order_date             string `json:"order_date"`
	Order_number           string `json:"order_number"`
	Sink_coeff             int    `json:"sink_coeff"`
	Sink_time              int    `json:"sink_time"`
	Sink_type              int    `json:"sink_type"`
	Suppliers_id           int    `json:"suppliers_id"`
	Use_date               string `json:"use_date"`
	Value                  string `json:"value"`
	Warranty_date          string `json:"warranty_date"`
	Warranty_duration      int    `json:"warranty_duration"`
	Warranty_info          string `json:"warranty_info"`
	Warranty_value         string `json:"warranty_value"`
}

type GlpiTicket struct {
	Id                         int    `json:"id"`
	Entities_id                int    `json:"entities_id"`
	Name                       string `json:"name"`
	Date                       string `json:"date"`
	Closedate                  string `json:"closedate"`
	Solvedate                  string `json:"solvedate"`
	Status                     int    `json:"status"`
	Users_id_lastupdater       int    `json:"users_id_lastupdater"`
	Users_id_recipient         int    `json:"users_id_recipient"`
	Requesttypes_id            int    `json:"requesttypes_id"`
	Content                    string `json:"content"`
	Urgency                    int    `json:"urgency"`
	Impact                     int    `json:"impact"`
	Priority                   int    `json:"priority"`
	Itilcategories_id          int    `json:"itilcategories_id"`
	Ticket_type                int    `json:"type"`
	Global_validation          int    `json:"global_validation"`
	Slas_id_ttr                int    `json:"slas_id_ttr"`
	Slas_id_tto                int    `json:"slas_id_tto"`
	Slalevels_id_ttr           int    `json:"slalevels_id_ttr"`
	Time_to_resolve            string `json:"time_to_resolve"`
	Time_to_own                string `json:"time_to_own"`
	Begin_waiting_date         string `json:"begin_waiting_date"`
	Sla_waiting_durationrity   int    `json:"sla_waiting_duration"`
	Ola_waiting_duration       int    `json:"ola_waiting_duration"`
	Olas_id_tto                int    `json:"olas_id_tto"`
	Olas_id_ttr                int    `json:"olas_id_ttr"`
	Olalevels_id_ttr           int    `json:"olalevels_id_ttr"`
	Ola_ttr_begin_date         string `json:"ola_ttr_begin_date"`
	Internal_time_to_resolve   string `json:"internal_time_to_resolve"`
	Internal_time_to_own       string `json:"internal_time_to_own"`
	Waiting_duration           int    `json:"waiting_duration"`
	Close_delay_stat           int    `json:"close_delay_stat"`
	Solve_delay_stat           int    `json:"solve_delay_stat"`
	Takeintoaccount_delay_stat int    `json:"takeintoaccount_delay_stat"`
	Actiontime                 int    `json:"actiontime"`
	Is_deleted                 int    `json:"is_deleted"`
	Locations_id               int    `json:"locations_id"`
	Validation_percent         int    `json:"validation_percent"`
}

type GlpiOSItem struct {
	Id                               int    `json:"id"`
	Items_id                         int    `json:"items_id"`
	Itemtype                         string `json:"itemtype"`
	Entities_id                      int    `json:"entities_id"`
	Operatingsystems_id              int    `json:"operatingsystems_id"`
	Operatingsystemversions_id       int    `json:"operatingsystemversions_id"`
	Operatingsystemservicepacks_id   int    `json:"operatingsystemservicepacks_id"`
	Operatingsystemarchitectures_id  int    `json:"operatingsystemarchitectures_id"`
	Operatingsystemkernelversions_id int    `json:"operatingsystemkernelversions_id"`
	Operatingsystemeditions_id       int    `json:"operatingsystemeditions_id"`
	License_number                   string `json:"license_number"`
	Licenseid                        string `json:"licenseid"`
}

type GlpiPeripheral struct {
	Id                  int    `json:"id"`
	Peripheralmodels_id int    `json:"peripheralmodels_id"`
	Peripheraltypes_id  int    `json:"peripheraltypes_id"`
	Comment             string `json:"comment"`
	Contact             string `json:"contact"`
	Contact_num         string `json:"contact_num"`
	Entities_id         int    `json:"entities_id"`
	Groups_id           int    `json:"groups_id"`
	Groups_id_tech      int    `json:"groups_id_tech"`
	Is_global           int    `json:"is_global"`
	Is_template         int    `json:"is_template"`
	Is_deleted          int    `json:"is_deleted"`
	Locations_id        int    `json:"locations_id"`
	Manufacturers_id    int    `json:"manufacturers_id"`
	Name                string `json:"name"`
	Otherserial         string `json:"otherserial"`
	Brand               string `json:"brand"`
	Serial              string `json:"serial"`
	States_id           int    `json:"states_id"`
	Template_name       string `json:"template_name"`
	Users_id            int    `json:"users_id"`
	Users_id_tech       int    `json:"users_id_tech"`
}

type GlpiPrinter struct {
	Id                 int    `json:"id"`
	Printermodels_id   int    `json:"printermodels_id"`
	Printertypes_id    int    `json:"printertypes_id"`
	Comment            string `json:"comment"`
	Contact            string `json:"contact"`
	Contact_num        string `json:"contact_num"`
	Is_global          int    `json:"is_global"`
	Groups_id          int    `json:"groups_id"`
	Groups_id_tech     int    `json:"groups_id_tech"`
	Have_serial        int    `json:"have_serial"`
	Have_parallel      int    `json:"have_parallel"`
	Have_usb           int    `json:"have_usb"`
	Have_wifi          int    `json:"have_wifi"`
	Have_ethernet      int    `json:"have_ethernet"`
	Memory_size        string `json:"memory_size"`
	Init_pages_counter int    `json:"init_pages_counter"`
	Last_pages_counter int    `json:"last_pages_counter"`
	Locations_id       int    `json:"locations_id"`
	Manufacturers_id   int    `json:"manufacturers_id"`
	Name               string `json:"name"`
	Networks_id        int    `json:"networks_id"`
	Otherserial        string `json:"otherserial"`
	Serial             string `json:"serial"`
	States_id          int    `json:"states_id"`
	Users_id           int    `json:"users_id"`
	Users_id_tech      int    `json:"users_id_tech"`
	Is_template        int    `json:"is_template"`
	Template_name      string `json:"template_name"`
	Is_deleted         int    `json:"is_deleted"`
}

// https://glpi/apirest.php/cartridge
type GlpiCartridge struct {
	Id                int    `json:"id"`
	Cartridgeitems_id int    `json:"cartridgeitems_id"`
	Entities_id       int    `json:"entities_id"`
	Printers_id       int    `json:"printers_id"`
	Date_in           string `json:"date_in"`
	Date_use          string `json:"date_use"`
	Date_out          string `json:"date_out"`
	Pages             int    `json:"pages"`
}

// https://glpi/apirest.php/cartridgeitem
type GlpiCartridgeItem struct {
	Id                    int    `json:"id"`
	Entities_id           int    `json:"entities_id"`
	Locations_id          int    `json:"locations_id"`
	Name                  string `json:"name"`
	Ref                   string `json:"ref"`
	Cartridgeitemtypes_id int    `json:"cartridgeitemtypes_id"`
	Manufacturers_id      int    `json:"manufacturers_id"`
	Users_id_tech         int    `json:"users_id_tech"`
	Groups_id_tech        int    `json:"groups_id_tech"`
	Comment               string `json:"comment"`
	Alarm_threshold       int    `json:"alarm_threshold"`
	Is_deleted            int    `json:"is_deleted"`
}

// https://glpi/apirest.php/cartridgeitem_printermodel
type GlpiCartridgeItemPrinterModel struct {
	Id                int `json:"id"`
	Cartridgeitems_id int `json:"cartridgeitems_id"`
	Printermodels_id  int `json:"printermodels_id"`
}

// https://glpi/apirest.php/Consumable
type GlpiConsumable struct {
	Id                 int    `json:"id"`
	Consumableitems_id int    `json:"consumableitems_id"`
	Entities_id        int    `json:"entities_id"`
	Date_in            string `json:"date_in"`
	Date_out           string `json:"date_out"`
	Itemtype           string `json:"itemtype"`
	Items_id           int    `json:"items_id"`
}

//  https://glpi/apirest.php/consumableitem
type GlpiConsumableItem struct {
	Id                     int    `json:"id"`
	Entities_id            int    `json:"entities_id"`
	Locations_id           int    `json:"locations_id"`
	Name                   string `json:"name"`
	Ref                    string `json:"ref"`
	Сonsumableitemtypes_id int    `json:"consumableitemtypes_id"`
	Manufacturers_id       int    `json:"manufacturers_id"`
	Users_id_tech          int    `json:"users_id_tech"`
	Groups_id_tech         int    `json:"groups_id_tech"`
	Comment                string `json:"comment"`
	Alarm_threshold        int    `json:"alarm_threshold"`
	Otherserial            string `json:"otherserial"`
	Is_deleted             int    `json:"is_deleted"`
}

// http://glpi/apirest.php/phone
type GlpiPhone struct {
	Id                    int    `json:"id"`
	Brand                 string `json:"brand"`
	Phonemodels_id        int    `json:"phonemodels_id"`
	Phonetypes_id         int    `json:"phonetypes_id"`
	Comment               string `json:"comment"`
	Contact               string `json:"contact"`
	Contact_num           string `json:"contact_num"`
	Entities_id           int    `json:"entities_id"`
	Groups_id             int    `json:"groups_id"`
	Groups_id_tech        int    `json:"groups_id_tech"`
	Is_template           int    `json:"is_template"`
	Is_global             int    `json:"is_global"`
	Is_deleted            int    `json:"is_deleted"`
	Locations_id          int    `json:"locations_id"`
	Manufacturers_id      int    `json:"manufacturers_id"`
	Name                  string `json:"name"`
	Phonepowersupplies_id int    `json:"phonepowersupplies_id"`
	Otherserial           string `json:"otherserial"`
	Serial                string `json:"serial"`
	States_id             int    `json:"states_id"`
	Users_id              int    `json:"users_id"`
	Users_id_tech         int    `json:"users_id_tech"`
	Template_name         string `json:"template_name"`
	Have_headset          int    `json:"have_headset"`
	Have_hp               int    `json:"have_hp"`
	Number_line           string `json:"number_line"`
}

// http://glpi/apirest.php/rack
type GlpiRack struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	Comment          string `json:"comment"`
	Entities_id      int    `json:"entities_id"`
	Locations_id     int    `json:"locations_id"`
	Serial           string `json:"serial"`
	Otherserial      string `json:"otherserial"`
	Rackmodels_id    int    `json:"phonemodels_id"`
	Manufacturers_id int    `json:"manufacturers_id"`
	Racktypes_id     int    `json:"phonetypes_id"`
	States_id        int    `json:"states_id"`
	Users_id_tech    int    `json:"users_id_tech"`
	Groups_id_tech   int    `json:"groups_id_tech"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
	Depth            int    `json:"depth"`
	Number_units     int    `json:"number_units"`
	Is_template      int    `json:"is_template"`
	Template_name    string `json:"template_name"`
	Dcrooms_id       int    `json:"dcrooms_id"`
	Room_orientation int    `json:"room_orientation"`
	Position         string `json:"position"`
	Bgcolor          string `json:"bgcolor"`
	Max_power        int    `json:"max_power"`
	Mesured_power    int    `json:"mesured_power"`
	Max_weight       int    `json:"max_weight"`
	Is_deleted       int    `json:"is_deleted"`
}

// http://glpi/apirest.php/item_rack
type GlpiRackItem struct {
	Id          int    `json:"id"`
	Racks_id    int    `json:"racks_id"`
	Itemtype    string `json:"itemtype"`
	Items_id    int    `json:"items_id"`
	Position    int    `json:"position"`
	Orientation int    `json:"orientation"`
	Bgcolor     string `json:"bgcolor"`
	Hpos        int    `json:"hpos"`
	Is_reserved int    `json:"is_reserved"`
}

// http://glpi/apirest.php/enclosure
type GlpiEnclosure struct {
	Id                 int    `json:"id"`
	Name               string `json:"name"`
	Entities_id        int    `json:"entities_id"`
	Locations_id       int    `json:"locations_id"`
	Serial             string `json:"serial"`
	Otherserial        string `json:"otherserial"`
	Enclosuremodels_id int    `json:"enclosuremodels_id"`
	Users_id_tech      int    `json:"users_id_tech"`
	Groups_id_tech     int    `json:"groups_id_tech"`
	Is_template        int    `json:"is_template"`
	Template_name      string `json:"template_name"`
	Orientation        int    `json:"orientation"`
	Power_supplies     int    `json:"power_supplies"`
	States_id          int    `json:"states_id"`
	Comment            string `json:"comment"`
	Manufacturers_id   int    `json:"manufacturers_id"`
	Is_deleted         int    `json:"is_deleted"`
}

// http://glpi/apirest.php/enclosuremodel
type GlpiEnclosureModel struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	Comment           string `json:"comment"`
	Product_number    int    `json:"product_number"`
	Weight            int    `json:"weight"`
	Required_units    string `json:"required_units"`
	Depth             string `json:"depth"`
	Power_connections int    `json:"power_connections"`
	Power_consumption int    `json:"power_consumption"`
	Is_half_rack      int    `json:"is_half_rack"`
	Picture_front     string `json:"picture_front"`
	Picture_rear      string `json:"picture_rear"`
}

type GlpiEnclosureItem struct {
	Id            int    `json:"id"`
	Enclosures_id int    `json:"enclosures_id"`
	Itemtype      string `json:"itemtype"`
	Items_id      int    `json:"items_id"`
	Position      int    `json:"position"`
}

// http://glpi/apirest.php/pdu
type GlpiPDU struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	Entities_id      int    `json:"entities_id"`
	Locations_id     int    `json:"locations_id"`
	Serial           string `json:"serial"`
	Otherserial      string `json:"otherserial"`
	Pdumodels_id     int    `json:"pdumodels_id"`
	Users_id_tech    int    `json:"users_id_tech"`
	Groups_id_tech   int    `json:"groups_id_tech"`
	Is_template      int    `json:"is_template"`
	Template_name    string `json:"template_name"`
	States_id        int    `json:"states_id"`
	Comment          string `json:"comment"`
	Manufacturers_id int    `json:"manufacturers_id"`
	Pdutypes_id      int    `json:"pdutypes_id"`
	Is_deleted       int    `json:"is_deleted"`
}

// http://glpi/apirest.php/pduemodel
type GlpiPDUModel struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	Comment           string `json:"comment"`
	Product_number    int    `json:"product_number"`
	Weight            int    `json:"weight"`
	Required_units    string `json:"required_units"`
	Depth             string `json:"depth"`
	Power_connections int    `json:"power_connections"`
	Max_power         int    `json:"max_power"`
	Is_half_rack      int    `json:"is_half_rack"`
	Picture_front     string `json:"picture_front"`
	Picture_rear      string `json:"picture_rear"`
	Is_rackable       int    `json:"is_rackable"`
}

// http://glpi/apirest.php/pdu_plug
type GlpiPDUPlug struct {
	Id           int `json:"id"`
	Plugs_id     int `json:"plugs_id"`
	Pdus_id      int `json:"pdus_id"`
	Number_plugs int `json:"number_plugs"`
}

// http://glpi/apirest.php/passivedcequipment
type GlpiPassivedcEquipment struct {
	Id                          int    `json:"id"`
	Name                        string `json:"name"`
	Entities_id                 int    `json:"entities_id"`
	Locations_id                int    `json:"locations_id"`
	Serial                      string `json:"serial"`
	Otherserial                 string `json:"otherserial"`
	Passivedcequipmentmodels_id int    `json:"passivedcequipmentmodels_id"`
	Passivedcequipmenttypes_id  int    `json:"passivedcequipmenttypes_id"`
	Users_id_tech               int    `json:"users_id_tech"`
	Groups_id_tech              int    `json:"groups_id_tech"`
	Is_template                 int    `json:"is_template"`
	Template_name               string `json:"template_name"`
	States_id                   int    `json:"states_id"`
	Comment                     string `json:"comment"`
	Manufacturers_id            int    `json:"manufacturers_id"`
	Is_deleted                  int    `json:"is_deleted"`
}

// http://glpi/apirest.php/PassivedcEquipmentModel
type GlpiPassivedcEquipmentModel struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	Comment           string `json:"comment"`
	Product_number    int    `json:"product_number"`
	Weight            int    `json:"weight"`
	Required_units    string `json:"required_units"`
	Depth             string `json:"depth"`
	Power_connections int    `json:"power_connections"`
	Power_consumption int    `json:"power_consumption"`
	Is_half_rack      int    `json:"is_half_rack"`
	Picture_front     string `json:"picture_front"`
	Picture_rear      string `json:"picture_rear"`
}

// http://glpi/apirest.php/DeviceSimcard
type GlpiDeviceSimcard struct {
	Id                    int    `json:"id"`
	Designation           string `json:"name"`
	Comment               string `json:"comment"`
	Entities_id           int    `json:"entities_id"`
	Manufacturers_id      int    `json:"manufacturers_id"`
	Voltage               int    `json:"voltage"`
	Devicesimcardtypes_id int    `json:"devicesimcardtypes_id"`
	Allow_voip            int    `json:"allow_voip"`
}

// http://glpi/apirest.php/item_DeviceSimcard
type GlpiDeviceSimcardItem struct {
	Id                int    `json:"id"`
	Items_id          int    `json:"items_id"`
	Itemtype          string `json:"itemtype"`
	Devicesimcards_id int    `json:"devicesimcards_id"`
	Entities_id       int    `json:"entities_id"`
	Serial            string `json:"serial"`
	Otherserial       string `json:"otherserial"`
	States_id         int    `json:"states_id"`
	Locations_id      int    `json:"locations_id"`
	Lines_id          int    `json:"lines_id"`
	Users_id          int    `json:"users_id"`
	Groups_id         int    `json:"groups_id"`
	Pin               string `json:"pin"`
	Pin2              string `json:"pin2"`
	Puk               string `json:"puk"`
	Puk2              string `json:"puk2"`
	Msin              string `json:"msin"`
	Is_deleted        int    `json:"is_deleted"`
}

// http://glpi/apirest.php/line
type GlpiLine struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	Entities_id      int    `json:"entities_id"`
	Caller_num       string `json:"caller_num"`
	Caller_name      string `json:"caller_name"`
	Users_id         int    `json:"users_id"`
	Groups_id        int    `json:"groups_id"`
	Lineoperators_id int    `json:"lineoperators_id"`
	Locations_id     int    `json:"locations_id"`
	States_id        int    `json:"states_id"`
	Linetypes_id     int    `json:"linetypes_id"`
	Comment          string `json:"comment"`
	Is_deleted       int    `json:"is_deleted"`
}

// http://glpi/apirest.php/lineoperator
type GlpiLineOperator struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Comment     string `json:"comment"`
	Mcc         int    `json:"mcc"`
	Mnc         int    `json:"mnc"`
	Entities_id int    `json:"entities_id"`
}

// http://glpi/apirest.php/NetworkPort
type GlpiNetworkPort struct {
	Id                 int    `json:"id"`
	Items_id           int    `json:"items_id"`
	Itemtype           string `json:"itemtype"`
	Entities_id        int    `json:"entities_id"`
	Logical_number     int    `json:"logical_number"`
	Name               string `json:"name"`
	Instantiation_type string `json:"instantiation_type"`
	Mac                string `json:"mac"`
	Comment            string `json:"comment"`
}

// http://glpi/apirest.php/NetworkName
type GlpiNetworkName struct {
	Id          int    `json:"id"`
	Entities_id int    `json:"entities_id"`
	Items_id    int    `json:"items_id"`
	Itemtype    string `json:"itemtype"`
	Name        string `json:"name"`
	Comment     string `json:"comment"`
	Fqdns_id    int    `json:"fqdns_id"`
}

// http://glpi/apirest.php/ipaddress
type GlpiIPAddress struct {
	Id           int    `json:"id"`
	Entities_id  int    `json:"entities_id"`
	Items_id     int    `json:"items_id"`
	Itemtype     string `json:"itemtype"`
	Version      int    `json:"version"`
	Name         string `json:"name"`
	Comment      string `json:"comment"`
	Binary_0     int    `json:"binary_0"`
	Binary_1     int    `json:"binary_1"`
	Binary_2     int    `json:"binary_2"`
	Binary_3     int    `json:"binary_3"`
	Mainitems_id int    `json:"mainitems_id"`
	Mainitemtype string `json:"mainitemtype"`
}

// http://glpi/apirest.php/fqdn
type GlpiFQDN struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Comment     string `json:"comment"`
	FQDN        string `json:"fqdn"`
	Entities_id int    `json:"entities_id"`
}

// http://glpi/apirest.php/problem
type GlpiProblem struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name"`
	Entities_id          int    `json:"entities_id"`
	Is_deleted           int    `json:"is_deleted"`
	Status               int    `json:"status"`
	Content              string `json:"content"`
	Date                 string `json:"date"`
	Solvedate            string `json:"solvedate"`
	Closedate            string `json:"closedate"`
	Time_to_resolve      string `json:"time_to_resolve"`
	Users_id_recipient   int    `json:"users_id_recipient"`
	Users_id_lastupdater int    `json:"users_id_lastupdater"`
	Urgency              int    `json:"urgency"`
	Impact               int    `json:"impact"`
	Priority             int    `json:"priority"`
	Itilcategories_id    int    `json:"itilcategories_id"`
	Impactcontent        string `json:"impactcontent"`
	Causecontent         string `json:"causecontent"`
	Symptomcontent       string `json:"symptomcontent"`
	Actiontime           int    `json:"actiontime"`
	Begin_waiting_date   string `json:"begin_waiting_date"`
	Waiting_duration     int    `json:"waiting_duration"`
	Close_delay_stat     int    `json:"close_delay_stat"`
	Solve_delay_stat     int    `json:"solve_delay_stat"`
}

// http://glpi/apirest.php/change
type GlpiChange struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name"`
	Entities_id          int    `json:"entities_id"`
	Is_deleted           int    `json:"is_deleted"`
	Status               int    `json:"status"`
	Content              string `json:"content"`
	Date                 string `json:"date"`
	Solvedate            string `json:"solvedate"`
	Closedate            string `json:"closedate"`
	Time_to_resolve      string `json:"time_to_resolve"`
	Users_id_recipient   int    `json:"users_id_recipient"`
	Users_id_lastupdater int    `json:"users_id_lastupdater"`
	Urgency              int    `json:"urgency"`
	Impact               int    `json:"impact"`
	Priority             int    `json:"priority"`
	Itilcategories_id    int    `json:"itilcategories_id"`
	Impactcontent        string `json:"impactcontent"`
	Controlistcontent    string `json:"controlistcontent"`
	Rolloutplancontent   string `json:"rolloutplancontent"`
	Backoutplancontent   string `json:"backoutplancontent"`
	Checklistcontent     string `json:"checklistcontent"`
	Global_validation    int    `json:"global_validation"`
	Validation_percent   int    `json:"validation_percent"`
	Actiontime           int    `json:"actiontime"`
	Begin_waiting_date   string `json:"begin_waiting_date"`
	Waiting_duration     int    `json:"waiting_duration"`
	Close_delay_stat     int    `json:"close_delay_stat"`
	Solve_delay_stat     int    `json:"solve_delay_stat"`
}

// http://glpi/apirest.php/ticketrecurrent
type GlpiTicketrecurrent struct {
	Id                 int    `json:"id"`
	Name               string `json:"name"`
	Entities_id        int    `json:"entities_id"`
	Is_active          int    `json:"is_active"`
	Tickettemplates_id int    `json:"tickettemplates_id"`
	Begin_date         string `json:"begin_date"`
	Periodicity        string `json:"periodicity"`
	Create_before      int    `json:"create_before"`
	Next_creation_date string `json:"next_creation_date"`
	Calendars_id       int    `json:"calendars_id"`
	End_date           string `json:"end_date"`
}

// This is a common data type for simple GLPI entities, like as:
// computertype, monitortype, networkequipmenttype, consumableitemtype, phonetype
// phonepowersupply, racktype, pdu, plug, passivedcequipmenttype, linetype
type GlpiSimpleItem struct {
	Id      int    `json:"id"`
	Comment string `json:"comment"`
	Name    string `json:"name"`
}

// This is a common data type for simple GLPI entities, like as:
// PrinterModel, PhoneModel, RackModel, DeviceCaseModel, DeviceDriveModel, DeviceGenericModel,
// DeviceGraphicCardModel, DeviceHardDriveModel, DeviceMemoryModel, DeviceMotherBoardModel,
// DeviceNetworkCardModel, DevicePciModel, DevicePowerSupplyModel, DeviceProcessorModel,
// DeviceSoundCardModel, DeviceSensorModel
type GlpiItemModel struct {
	Id             int    `json:"id"`
	Comment        string `json:"comment"`
	Name           string `json:"name"`
	Product_number string `json:"product_number"`
}

// Http client setting
var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: 60 * time.Second,
		DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
		TLSClientConfig: &tls.Config{
			MaxVersion:         tls.VersionTLS11,
			InsecureSkipVerify: true,
		},
	},
}

func NewSession(server, apiUserToken, apiAppToken string) (*Session, error) {
	return &Session{server, apiUserToken, apiAppToken, ""}, nil
}

func (glpi *Session) GetSessionToken() string {
	return glpi.sessionToken
}

// Init GLPI-session.
// Return glpi session token
func (glpi *Session) InitSession() (string, error) {
	request, err := http.NewRequest("GET", glpi.url+"/initSession", nil)
	if err != nil {
		fmt.Println(request, err)
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "user_token "+glpi.userToken)
	request.Header.Add("App-Token", glpi.appToken)
	// fmt.Println(request)

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	// fmt.Println(result["sessionToken"])

	glpi.sessionToken = result["session_token"].(string)

	return result["session_token"].(string), nil
}

// Kill GLPI-session.
// Return response code and error
func (glpi *Session) KillSession() (int, error) {
	request, err := http.NewRequest("GET", glpi.url+"/killSession", nil)
	if err != nil {
		fmt.Println(request, err)
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Session-Token", glpi.sessionToken)
	request.Header.Add("App-Token", glpi.appToken)

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	return resp.StatusCode, err
}

// Making http(s) request for metthod PUT and POST
func (glpi *Session) MakeRequest(method string, requestURL string, data map[string]interface{}) (map[string]interface{}, error) {
	encodedData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(encodedData, err)
	}
	request, err := http.NewRequest(method, glpi.url+"/"+requestURL+"/", bytes.NewBuffer(encodedData))

	if err != nil {
		fmt.Println(request, err)
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Session-Token", glpi.sessionToken)
	request.Header.Add("App-Token", glpi.appToken)

	response, err := client.Do(request)

	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()

	var result map[string]interface{}
	// This is a magic. Because glpi PUT(DELETE) return an array [], but POST return {}
	switch method {
	case "PUT":
		var r []map[string]interface{}
		json.NewDecoder(response.Body).Decode(&r)
		var res map[string]interface{}
		// fmt.Println("length", len(r[0]))
		if len(r[0]) > 0 {
			res = map[string]interface{}{
				"status": 1,
			}
		} else {
			res = map[string]interface{}{
				"status": 0,
			}
		}
		result = res
		// fmt.Println(result["status"])
	case "DELETE":
		var r []map[string]interface{}
		json.NewDecoder(response.Body).Decode(&r)
		var res map[string]interface{}
		if len(r[0]) > 0 {
			res = map[string]interface{}{
				"status": 1,
			}
		} else {
			res = map[string]interface{}{
				"status": 0,
			}
		}
		result = res
		// fmt.Println(result["status"])
	case "POST":
		// var result map[string]interface{}
		json.NewDecoder(response.Body).Decode(&result)
	}
	// fmt.Println(result)
	return result, nil

}

// Making http(s) request for method GET
func (glpi *Session) MakeRequestGET(requestURL string) (bytes.Buffer, error) {
	request, err := http.NewRequest("GET", glpi.url+"/"+requestURL, nil)

	if err != nil {
		fmt.Println(request, err)
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Session-Token", glpi.sessionToken)
	request.Header.Add("App-Token", glpi.appToken)

	response, err := client.Do(request)

	if err != nil {
		fmt.Println(err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, response.Body)
	if err != nil {
		return buf, err
	}
	response.Body.Close()
	return buf, nil
}

// Add GLPI item
// Return new item id
func (glpi *Session) AddItem(itemType string, data map[string]interface{}) int {
	response, err := glpi.MakeRequest("POST", itemType, data)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(response)
	return int(response["id"].(float64))
}

// Update GLPI item
// Return result
func (glpi *Session) UpdateItem(itemType string, data map[string]interface{}) int {
	response, err := glpi.MakeRequest("PUT", itemType+"/", data)
	if err != nil {
		fmt.Println(err)
	}

	res, _ := strconv.Atoi(fmt.Sprint(response["status"]))
	return res

}

// Delete GLPI item
// Return result
func (glpi *Session) DeleteItem(itemType string, data map[string]interface{}) int {
	response, err := glpi.MakeRequest("DELETE", itemType+"/", data)
	if err != nil {
		fmt.Println(err)
	}
	res, _ := strconv.Atoi(fmt.Sprint(response["status"]))
	return res

}

// Get GLPI config
func (glpi *Session) GetConfig() map[string]map[string]interface{} {
	response, err := glpi.MakeRequestGET("getGlpiConfig")
	if err != nil {
		fmt.Println(err)
	}
	var result map[string]map[string]interface{}
	json.Unmarshal(response.Bytes(), &result)
	return result
}

// Get GLPI version
func (glpi *Session) Version() string {
	result := glpi.GetConfig()
	return fmt.Sprint(result["cfg_glpi"]["version"])
}

// Get one Item.
func (glpi *Session) GetItem(itemType string, itemID int, otherParam string) []byte {
	response, err := glpi.MakeRequestGET(itemType + "/" + strconv.Itoa(itemID) + string('?') + otherParam)
	if err != nil {
		fmt.Println(err)
	}
	return response.Bytes()
}

// Search Item. Return item id if object was found
func (glpi *Session) SearchItem(itemType string, itemName string) int {
	response, err := glpi.MakeRequestGET(itemType + "/?range=0-10000")
	if err != nil {
		fmt.Println(err)
	}
	var result []map[string]interface{}
	json.Unmarshal(response.Bytes(), &result)
	var found_id int
	for _, item := range result {
		if strings.EqualFold(itemName, fmt.Sprint(item["name"])) {
			found_id, _ = strconv.Atoi(fmt.Sprint(item["id"]))
		}
	}
	return found_id
}

// Update item status
// http://glpi/apirest.php/_glpi_itemType_
func (glpi *Session) UpdateItemStatus(itemType string, itemID int, itemStatusID int) {
	message := map[string]interface{}{
		"input": map[string]int{
			"id":        itemID,
			"states_id": itemStatusID,
		},
	}
	fmt.Println(glpi.UpdateItem(itemType, message))
}

//-------------------------------------------------------------------------------------------------
func (glpi *Session) ItemOperation(operation string, glpiItemType string, data interface{}) int {
	message := map[string]interface{}{
		"input": data,
	}
	var response int
	switch operation {
	case "add":
		response = glpi.AddItem(glpiItemType, message)
	case "update":
		response = glpi.UpdateItem(glpiItemType, message)
	case "delete":
		response = glpi.DeleteItem(glpiItemType, message)
	}
	return response
}
