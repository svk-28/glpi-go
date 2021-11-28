GLPI
====

This Go library implements the GLPI HTTP REST API. This module provides working with all the assets and some GLPI items.

Getting started
===============

## Setting the variables

Setting a three variables for connecting a GLPI API:

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

Where:

- GLPI_USER_TOKEN - the API access token for each user (setting in a user profile)
- GLPI_APP_TOKEN - the application token (setting in the "Setup"->"General"->"API" section GLPI UI)
- GLPI_SERVER_API - API URL (like as https://glpi/apirest.php)

## Connection to the API and initial session

``` go
glpiSession, err := glpi.NewSession(glpiServerAPI, glpiUserToken, glpiAppToken)
if err != nil {
    fmt.Println(err)
    return
}
_, err = glpiSession.InitSession()
if err != nil {
    fmt.Println(err)
    return
}
```
glpiSession.InitSession() - initializes the session, i.e. when this function is called, the module is connected to the GLPI and the session access code (token) is returned. Its value is stored in the field of the "glpi Session.session Token" structure. This code can be viewed by calling the function "glpi Session.Get Session Token ()", i.e.

``` go
fmt.Println(glpiSession.GetSessionToken())
```

Work with some items
====================

## Adding the Manufacturer

``` go
manufacturersData := glpi.GlpiManufacturer{
    Name:    "Yet Another",
    Comment: "Yet Another inc.",
}
manufacturersID = glpiSession.ItemOperation("add", "Manufacturer", manufacturersData)
```

If it's all right, the new manufacturers ID will be written to the "manufacturersID" variable.

## Adding the printers model

For simple model types, you should use the structure GlpiItemModel and ItemModel() function:

``` go
type GlpiSimpleItem struct {
	Id      int    `json:"id"`
	Comment string `json:"comment"`
	Name    string `json:"name"`
}
```

Adding the new printer model:

``` go
printerModelData := glpi.GlpiItemModel{
    Name: "Sys 1200",
}
printerModelID = glpiSession.ItemOperation("add", "PrinterModel", printerModelData)
```

## Adding the printer

``` go
// Add new printer
printerData := glpi.GlpiPrinter{
    // Id:   "11",
    Name:             "First printer",
    Have_usb:         1,
    Manufacturers_id: manufacturersID,
    Printermodels_id: printerModelID,
}

printerID := glpiSession.ItemOperation("add", "Printer", printerData)
fmt.Println("Printer was added, ID:", printerID)
```

## Adding a cartridge model and linking them to the printer model

``` go
// ADD cartridge
requestDataCartridgeItem := glpi.GlpiCartridgeItem{
    // Id:   "11",
    Name: "С4092B",
    Manufacturers_id: manufacturersID,
    Alarm_threshold:  5,
}
cartridgeItemsID := glpiSession.ItemOperation("add", "CartridgeItem", requestDataCartridgeItem)
fmt.Println("Add a cartridge item:", cartridgeItemsID)

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

## Adding a cartridge and linking them to the printer

``` go
requestDataCartridge := glpi.GlpiCartridge{
    // Id:   "11",
    Cartridgeitems_id: cartridgeItemsID,
}

cartridgeID := glpiSession.ItemOperation("add", "Cartridge", requestDataCartridge)
fmt.Println("Added a cartridge:", cartridgeID)

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

## Search the items 

To get the element ID, the module provides the ```SearchItem()``` function:

``` func (glpi *Session) SearchItem(itemType string, itemName string) int {}```

If an element with the specified name was found, the function returns its ID.

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

Getting ended
=============

## Kill the session

After all operations are completed, the session should be closed

``` go
_, err = glpiSession.KillSession()
if err != nil {
    fmt.Println(err)
    return
}
```

License
=======

Released under the GNU GPL License