package types

// info about the class that Meeting struct is pointing to
type MeetingInfo struct {
	ID        string `json:"id" firestore:"id"`               // CRN as the ID ex) "10492"
	CourseID  string `json:"course_id" firestore:"course_id"` // "CS110"
	Section   string `json:"section" firestore:"section"`     // "001"
	Professor string `json:"professor" firestore:"professor"`
}

// meeting time for a class in a specific room
type Meeting struct {
	Day       int           `json:"day" firestore:"day"`
	StartTime int           `json:"start_time" firestore:"start_time"`
	EndTime   int           `json:"end_time" firestore:"end_time"`
	Location  string        `json:"location" firestore:"location"`
	Label     []MeetingInfo `json:"label,omitempty" firestore:"label,omitempty"`
}

// a classroom with its aggregated schedule
type Room struct {
	ID       string    `firestore:"id"`       // ex) "HORIZN_2014"
	Building string    `firestore:"building"` // " Horizon Hall"
	Number   string    `firestore:"number"`   // "2014"
	Schedule []Meeting `firestore:"schedule"`
}
