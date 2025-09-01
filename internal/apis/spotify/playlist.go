package spotify

type PlaylistResponse struct {
	Tracks struct {
		Items []struct {
			Track struct {
				ExternalIDs struct {
					ISRC string `json:"isrc"`
				} `json:"external_ids"`
			} `json:"track"`
		} `json:"items"`
	} `json:"tracks"`
}
