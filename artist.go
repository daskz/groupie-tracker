package main

type artist struct {
	ID           int
	Name         string
	CreationDate int
	FirstAlbum   string
	Image        string
	Members      []string
}

type relation struct {
	Index []artistRelation
}

type artistRelation struct {
	ID             int
	DatesLocations map[string]interface{}
}

type artistData struct {
	Artist         artist
	DatesLocations map[string]interface{}
}
