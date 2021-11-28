module zabbix2glpi

go 1.16

replace glpi => ../glpi

replace zabbix => ../../zabbix/zabbix

require (
	glpi v0.0.0-00010101000000-000000000000
	zabbix v0.0.0-00010101000000-000000000000
)
