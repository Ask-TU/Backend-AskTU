package models

type AllClass struct {
	Subject_name string     `json:"subject_name"`
	Class_owner  string     `json:"class_owner"`
	Question     []Question `json:"question"`
	Member       []member   `json:"member"`
}

type Question struct {
	Question string `json:"question"`
}

type member struct {
	Member string `json:"member"`
}
