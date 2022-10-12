package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = API("objects", func() {
	Title("Near earth objects")
	Description("Near earth objects - CRUD")
	Server("objects", func() {
		Host("localhost", func() {
			URI("http://localhost:8080")
		})
	})
})

var _ = Service("objects", func() {
	Description("API Server")

	Method("list", func() {
		Payload(func() {
			Attribute("ObservationDate", String)
			Attribute("Offset", Int)
			Required("ObservationDate", "Offset")
		})

		Result(CollectionOf(nearEarthObjectManagement))
		Error("wrong_date_format")
		Error("wrong_offset_format")
		Error("internal_server")

		HTTP(func() {
			GET("/objects")
			Param("ObservationDate")
			Param("Offset")
			Response("wrong_date_format", StatusBadRequest)
			Response("wrong_offset_format", StatusBadRequest)
			Response("internal_server", StatusInternalServerError)
		})
	})

	Method("find", func() {
		Payload(func() {
			Attribute("Id", Int)
			Required("Id")
		})

		Result(nearEarthObjectManagement)
		Error("not_found")
		Error("internal_server")

		HTTP(func() {
			GET("/objects/{Id}")
			Response("not_found", StatusNotFound)
			Response("internal_server", StatusInternalServerError)
		})
	})

	Method("add", func() {
		// Payload(nearEarthObject)
		Payload(func() {
			Attribute("Id", Int)
			Attribute("Name", String)
			Attribute("ObservationTime", Int64)
			Required("Id", "Name", "ObservationTime")
		})

		Result(Boolean)
		Error("wrong_id")
		Error("wrong_time")
		Error("already_exists")
		Error("internal_server")

		HTTP(func() {
			POST("/objects")
			Body(func() {
				Attribute("Id")
				Attribute("Name")
				Attribute("ObservationTime")
			})
			Response("wrong_id", StatusBadRequest)
			Response("wrong_time", StatusBadRequest)
			Response("already_exists", StatusBadRequest)
			Response("internal_server", StatusInternalServerError)
		})
	})

	Method("update", func() {
		Payload(func() {
			Attribute("Id", Int)
			Attribute("Name", String)
			Attribute("ObservationTime", Int64)
			Required("Id")
		})

		Result(Boolean)
		Error("wrong_id")
		Error("wrong_time")
		Error("nothing_to_update")
		Error("internal_server")

		HTTP(func() {
			PATCH("/objects/{Id}")
			Param("Name")
			Param("ObservationTime")
			Response("wrong_id", StatusBadRequest)
			Response("wrong_time", StatusBadRequest)
			Response("nothing_to_update", StatusBadRequest)
			Response("internal_server", StatusInternalServerError)
		})
	})

	Method("delete", func() {
		Payload(func() {
			Attribute("Id", Int)
			Required("Id")
		})

		Result(Boolean)
		Error("not_found")
		Error("internal_server")

		HTTP(func() {
			DELETE("objects/{Id}")
			Response("not_found", StatusNotFound)
			Response("internal_server", StatusInternalServerError)
		})
	})

	Files("/openapi.json", "./gen/http/openapi.json")
})

var nearEarthObjectManagement = ResultType("application/vnd.object", func() {
	TypeName("nearEarthObjectManagement")
	Attribute("Id", Int, "object ID")
	Attribute("Name", String, "object Name")
	Attribute("ObservationTime", Int64, "object observation time")
	Attribute("ObservationDate", String, "object observation date")
	Required("Id", "Name", "ObservationTime", "ObservationDate")
})
