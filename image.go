package goss

type (
	ImageResponse struct {
		ID           string `json:"id,omitempty"`
		LocationID   string `json:"location_id,omitempty"`
		Type         string `json:"type,omitempty"`
		OSVersion    string `json:"os_version,omitempty"`
		Architecture string `json:"architecture,omitempty"`
		AllowSSHKeys bool   `json:"allow_ssh_keys,omitempty"`
	}

	imageListResponseWrap struct {
		Images []*ImageResponse `json:"locations,omitempty"`
	}
)

func (c *SSClient) GetImageList() ([]*ImageResponse, error) {
	url := getImageBaseURL()
	resp, err := makeRequest(c.client, url, methodGet, nil, &imageListResponseWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*imageListResponseWrap).Images, nil
}

func getImageBaseURL() string {
	return "images"
}
