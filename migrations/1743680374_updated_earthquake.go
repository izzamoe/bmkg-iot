package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_428798634")
		if err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(1, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text512383393",
			"max": 0,
			"min": 0,
			"name": "Coordinates",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(2, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text2464052163",
			"max": 0,
			"min": 0,
			"name": "Lintang",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(3, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text301170639",
			"max": 0,
			"min": 0,
			"name": "Bujur",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(4, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text3078029218",
			"max": 0,
			"min": 0,
			"name": "Magnitude",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(5, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text3610236738",
			"max": 0,
			"min": 0,
			"name": "Kedalaman",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(6, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text4010982572",
			"max": 0,
			"min": 0,
			"name": "Wilayah",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(7, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text897900443",
			"max": 0,
			"min": 0,
			"name": "Potensi",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(8, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text3171604681",
			"max": 0,
			"min": 0,
			"name": "Dirasakan",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(9, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text1267517429",
			"max": 0,
			"min": 0,
			"name": "Shakemap",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(10, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text3450443271",
			"max": 0,
			"min": 0,
			"name": "Jam",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(11, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text2094490243",
			"max": 0,
			"min": 0,
			"name": "Tanggal",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(12, []byte(`{
			"hidden": false,
			"id": "date3400993698",
			"max": "",
			"min": "",
			"name": "DateTime",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "date"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_428798634")
		if err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("text512383393")

		// remove field
		collection.Fields.RemoveById("text2464052163")

		// remove field
		collection.Fields.RemoveById("text301170639")

		// remove field
		collection.Fields.RemoveById("text3078029218")

		// remove field
		collection.Fields.RemoveById("text3610236738")

		// remove field
		collection.Fields.RemoveById("text4010982572")

		// remove field
		collection.Fields.RemoveById("text897900443")

		// remove field
		collection.Fields.RemoveById("text3171604681")

		// remove field
		collection.Fields.RemoveById("text1267517429")

		// remove field
		collection.Fields.RemoveById("text3450443271")

		// remove field
		collection.Fields.RemoveById("text2094490243")

		// remove field
		collection.Fields.RemoveById("date3400993698")

		return app.Save(collection)
	})
}
