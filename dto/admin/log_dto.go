package admin

type ActivityLogListItem struct {
	ID          string
	ActorName   string
	Action      string
	TargetID    string
	Description string
	IPAddress   string
	CreatedAt   string
}

type ActivityLogListResult struct {
	Items       []ActivityLogListItem
	CurrentPage int
	TotalPages  int
	TotalItems  int64
	Action      string
	Keyword     string
	HasPrev     bool
	HasNext     bool
}
