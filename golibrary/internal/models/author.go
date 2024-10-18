package models

type Author struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Times_ordered int    `json:"times_ordered,omitempty"`
	Books         []Book `json:"books,omitempty"`
}

// func (a *Author) MarshalJSON() ([]byte, error) {
//     type Alias Author
//     alias := (*Alias)(a)
//     return json.Marshal(&struct {
//         ID        int64    `json:"id"`
//         Name      string `json:"name"`
//         TimesOrdered int `json:"times_ordered"`
//         Books     []Book `json:"books,omitempty"`
//     }{
//         ID:        alias.ID,
//         Name:      alias.Name,
//         TimesOrdered: alias.TimesOrdered,
//         Books:     alias.Books,
//     })
// }
