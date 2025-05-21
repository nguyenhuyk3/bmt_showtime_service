package message

type BMTOrderPublicOutboxesMsg struct {
	// Before interface{}  `json:"before"`
	After AfterPayload `json:"after"`
	// Source      SourceInfo   `json:"source"`
	// Op          string       `json:"op"`
	// TsMs        int64        `json:"ts_ms"`
	// TsUs        int64        `json:"ts_us"`
	// TsNs        int64        `json:"ts_ns"`
	// Transaction interface{}  `json:"transaction"`
}

type PayloadData struct {
	Fab        []FabItem  `json:"fab"`
	Note       string     `json:"note"`
	Seats      []SeatItem `json:"seats"`
	OrderedBy  string     `json:"OrderedBy"`
	ShowDate   string     `json:"show_date"`
	ShowtimeId int32      `json:"showtime_id"`
}

type FabItem struct {
	FabID    int32 `json:"fab_id"`
	Quantity int32 `json:"quantity"`
}

type SeatItem struct {
	SeatID int32 `json:"seat_id"`
}

type AfterPayload struct {
	ID             string `json:"id"`
	AggregatedType string `json:"aggregated_type"`
	AggregatedID   int32  `json:"aggregated_id"`
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
