package applemusic

type PlaylistResponse struct {
	Data []struct {
		Attributes struct {
			PlayParams struct {
				ReportingID string `json:"reportingId"`
			} `json:"playParams"`
		} `json:"attributes"`
	} `json:"data"`
}
