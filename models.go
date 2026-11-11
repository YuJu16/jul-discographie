package main

// // Titre représente un morceau affiché sur le site
type Titre struct {
	ID         string
	Nom        string
	Artist     []string
	DateSortie string // Format : YYYY-MM-DD
	DureeMS    int    // Durée en millisecondes (comme spotify)
	Type       string // Type de la musique (single, album, EP, etc.)
	Pochette   string // URL de la pochette (lien de l'image)
}
