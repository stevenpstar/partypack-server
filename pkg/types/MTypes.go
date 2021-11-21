package types

type MPrompt struct {
	Prompt1 string `json:"prompt1"`
	Prompt2 string `json:"prompt2"`
	Image1  string `json:"image1"`
	Image2  string `json:"image2"`
}

// Player Image struct
type MPlayerImage struct {
  Name string `json:"PlayerName"`
  Image string `json:"image"`
}

