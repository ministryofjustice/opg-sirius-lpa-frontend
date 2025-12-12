package shared

type LpaStoreChange struct {
	Key      string `json:"key"`
	New      any    `json:"new"`
	Old      any    `json:"old"`
	Template string
	Readable string
}
