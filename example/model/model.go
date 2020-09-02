package model

import (
	"time"

	"github.com/bold-commerce/protoc-gen-struct-transformer/example/nulls"
)

type (
	Product struct {
		ID          int    `db:"id" json:"id"`
		Name        string `db:"name" json:"name"`
		One         string `db:"one" json:"one"`
		SecondID    string `db:"two" json:"two"`
		CustomField string `db:"custom_field" json:"custom_field"`
	}

	Order struct {
		ID       int
		FirstID  string
		SecondID string
		ThirdURL string
	}

	Address struct {
		ID   int
		Type string
	}

	Customer struct {
		ID             int
		Name           string
		Addresses      []Address
		DefaultAddress *Address
		BillingAddress Address
		MapField1      string
		MapField2      string
	}

	MyLineItem struct {
		ID   int
		Type string
		URL  string
		SKU  int
	}

	MyLineItemUsage struct {
		List []MyLineItem
		Item *MyLineItem
	}

	Value2Pointer struct {
		AddressNil *Address
	}

	Pointer2Value struct {
		AddressNotNil Address
	}

	// TimeModel is used for testing time-related transformations.
	TimeModel struct {
		TimeTime     time.Time
		PtrTimeTime  *time.Time
		NullsTime    nulls.Time
		PtrNullsTime *nulls.Time

		NullsTime2    nulls.Time
		PtrNullsTime2 *nulls.Time
	}

	IntsModel struct {
		IntFor64Value int
		IntFor32Value int
		Int32Value    int32
		Int64Value    int64
		StringValue   string
	}
)
