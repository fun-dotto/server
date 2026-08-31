package domain

type Calendar struct {
	ID        string
	ServiceID string
	Monday    int
	Tuesday   int
	Wednesday int
	Thursday  int
	Friday    int
	Saturday  int
	Sunday    int
	StartDate string
	EndDate   string
}
