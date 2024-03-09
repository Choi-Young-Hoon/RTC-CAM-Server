package rtccamweb

func NewPageData(page, imageServerUrl string) PageData {
	return PageData{
		Page:           page,
		ImageServerUrl: imageServerUrl,
	}
}

type PageData struct {
	Page string

	RoomRequestType string
	RequestId       string

	AuthToken string

	ImageServerUrl string
}
