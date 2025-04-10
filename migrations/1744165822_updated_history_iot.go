package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_320116043")
		if err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(1, []byte(`{
			"cascadeDelete": false,
			"collectionId": "pbc_4135069744",
			"hidden": false,
			"id": "relation154121870",
			"maxSelect": 1,
			"minSelect": 0,
			"name": "device",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(2, []byte(`{
			"hidden": false,
			"id": "bool1260321794",
			"name": "active",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "bool"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(3, []byte(`{
			"hidden": false,
			"id": "date2947204065",
			"max": "",
			"min": "",
			"name": "date_active",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "date"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_320116043")
		if err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("relation154121870")

		// remove field
		collection.Fields.RemoveById("bool1260321794")

		// remove field
		collection.Fields.RemoveById("date2947204065")

		return app.Save(collection)
	})
}
