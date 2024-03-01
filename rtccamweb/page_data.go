package rtccamweb

func NewPageData(page string) PageData {
	return PageData{
		Page: page,
	}
}

type PageData struct {
	Page string

	RoomRequestType string
	RequestId       string

	AuthToken string
}
