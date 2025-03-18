package trash

type Trash struct {
	TrashID   string `json:"trash_id"`
	TrashName string `json:"trash_name"`
	TrashCode string `json:"trash_code"`
	TrashJSON string `json:"trash_json"` // Можно использовать json.RawMessage для более сложных структур
}
