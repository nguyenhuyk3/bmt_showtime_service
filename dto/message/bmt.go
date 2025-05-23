package message

type BMTPublicOutboxesMsg struct {
	// Before interface{}  `json:"before"`
	After AfterPayload `json:"after"`
	// Source      SourceInfo   `json:"source"`
	// Op          string       `json:"op"`
	// TsMs        int64        `json:"ts_ms"`
	// TsUs        int64        `json:"ts_us"`
	// TsNs        int64        `json:"ts_ns"`
	// Transaction interface{}  `json:"transaction"`
}

type AfterPayload struct {
	ID             string `json:"id"`
	AggregatedType string `json:"aggregated_type"`
	AggregatedId   int32  `json:"aggregated_id"`
	EventType      string `json:"event_type"`
	Payload        string `json:"payload"`
	CreatedAt      int64  `json:"created_at"`
}

// type SourceInfo struct {
// 	Version   string      `json:"version"`
// 	Connector string      `json:"connector"`
// 	Name      string      `json:"name"`
// 	TsMs      int64       `json:"ts_ms"`
// 	Snapshot  string      `json:"snapshot"`
// 	DB        string      `json:"db"`
// 	Sequence  []string    `json:"sequence"`
// 	TsUs      int64       `json:"ts_us"`
// 	TsNs      int64       `json:"ts_ns"`
// 	Schema    string      `json:"schema"`
// 	Table     string      `json:"table"`
// 	TxID      int         `json:"txId"`
// 	LSN       int64       `json:"lsn"`
// 	XMin      interface{} `json:"xmin"`
// }
