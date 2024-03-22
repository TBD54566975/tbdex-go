package tbdex

import _ "embed"

// DefinitionsSchema imports the schema
//go:embed tbdex/hosted/json-schemas/definitions.json
var DefinitionsSchema []byte

// ResourceSchema imports the schema
//go:embed tbdex/hosted/json-schemas/resource.schema.json
var ResourceSchema []byte

// OfferingSchema imports the schema
//go:embed tbdex/hosted/json-schemas/offering.schema.json
var OfferingSchema []byte

// BalanceSchema imports the schema
//go:embed tbdex/hosted/json-schemas/balance.schema.json
var BalanceSchema []byte

// MessageSchema imports the schema
//go:embed tbdex/hosted/json-schemas/message.schema.json
var MessageSchema []byte

// RFQSchema imports the schema
//go:embed tbdex/hosted/json-schemas/rfq.schema.json
var RFQSchema []byte

// QuoteSchema imports the schema
//go:embed tbdex/hosted/json-schemas/quote.schema.json
var QuoteSchema []byte

// OrderSchema imports the schema
//go:embed tbdex/hosted/json-schemas/order.schema.json
var OrderSchema []byte

// CloseSchema imports the schema
//go:embed tbdex/hosted/json-schemas/close.schema.json
var CloseSchema []byte

// OrderStatusSchema imports the schema
//go:embed tbdex/hosted/json-schemas/orderstatus.schema.json
var OrderStatusSchema []byte
