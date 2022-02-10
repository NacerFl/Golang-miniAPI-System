// Réalisé par Léo et Nacer

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
)

// structure Attraction qui vas nous permettre de stocker nos objets
type Attraction struct {
	Id           int    `json:"Id"`
	Name         string `json:"Name"`
	InPark       string `json:"InPark"`
	Place        string `json:"Place"`
	Manufacturer string `json:"Manufacturer"`
}

// Permet d'exporter en json l'ensemble de nos attractions
// La fonction boucle sur l'ensemble de nos attractions,
//afin de les stocker dans un tableau puis les sauvgarder.
func saveInfile(attr []Attraction) {

	var data []Attraction

	// boucle sur les différents élement du tableau
	for i := 0; i < len(attr); i++ {

		//On ajoute au tableau nos Attraction
		data = append(data, Attraction{

			Id:           attr[i].Id,
			Name:         attr[i].Name,
			InPark:       attr[i].InPark,
			Place:        attr[i].Place,
			Manufacturer: attr[i].Manufacturer,
		})

	}

	//On utilise MarshalIndent du package Json afin de les mettre au bon format
	file, _ := json.MarshalIndent(data, "", " ")
	//On utlise WriteFile du package ioutl afin d'exporter le tout dans un fichier json
	_ = ioutil.WriteFile("test.json", file, 0644)
}

//Permet de supprimer une attraction
//Prend en parametre un id de l'attraction à delete
//Prend en parametre un tableau d'attraction
// La fonction boucle sur l'ensemble de nos attraction et supprime l'attraction avec l'id correspondant
func delete(id string, attr []Attraction) []Attraction {

	//on convert notre id en int
	var int_id, _ = strconv.Atoi(id)
	fmt.Println("Tableau de base :", attr)

	// boucle sur les id du tableau
	for i := 0; i < len(attr); i++ {
		if attr[i].Id == int_id {
			//Supprime notre attraction en concatenant nos tableau
			attr = append(attr[:i], attr[(i+1):]...)

		}

	}

	fmt.Println("Tableau delete : ", attr)
	return attr

}

//Permet de modifier une attraction
//Prend en parametre l'ensmble des parametre qui compose une attraction
//Prend en parametre un tableau  de structure (d'attraction)
// La fonction boucle sur les element du tableau et update les different element si il ne sont pas vide
func patch(id string, attr []Attraction, name string, inPark string, place string, manufacturer string) []Attraction {

	//On convert notre id en Int
	var int_id, _ = strconv.Atoi(id)

	for i := 0; i < len(attr); i++ {
		fmt.Println(int_id)
		if attr[i].Id == int_id {

			//Pour chaque élement on verifie si le champs n'est pas vide si oui on le modifie
			if len(name) != 0 {
				attr[i].Name = name
			}
			if len(inPark) != 0 {
				attr[i].InPark = inPark
			}
			if len(manufacturer) != 0 {
				attr[i].Manufacturer = manufacturer
			}
			if len(place) != 0 {
				attr[i].Place = place
			}
			fmt.Println("Attraction updated : ", attr[i])

			break
		}

	}

	return attr
}

//Permet de verifier la presence d'un ID
//La fonction boucle sur les element du tableau et verifie si l'id est present
//Si oui la fonction nous renvois True
func exist(attr []Attraction, id int) bool {
	for i := 0; i < len(attr); i++ {

		if attr[i].Id == id {

			return true

		}

	}
	return false
}

//Permet de crée une attraction
//Prend en parametre l'ensmble des parametre qui compose une attraction
//Prend en parametre un tableau  de structure (d'attraction)
// La fonction boucle sur les element du tableau et update les different element si il ne sont pas vide
func create(attr []Attraction, name string, inPark string, place string, manufacturer string) []Attraction {

	var id int

	//On verifie la presence de l'id à l'aide d'une boucle for
	for {
		id = rand.Intn(100000000000)
		//Notre fonction exist permet de verifier ça presence
		if exist(attr, id) {
			continue

		}
		break
	}

	//On crée notre attraction à l'aide de notre Structure
	var attraction = Attraction{id, name, inPark, place, manufacturer}
	//On encode notre attraction à l'aide de notre fonction encode()
	json := encode(attraction)
	//Puis on decode notre attraction à l'aide de notre fonction decode()
	attraction = decode(json)
	//On ajoute ensuite notre attraction à la liste des attractions
	attr = append(attr, attraction)

	return attr

}

//Permet d'encoder notre attraction en Json à l'aide de la fonction json.Marshal
//Prend en param une Atraction
func encode(s1 Attraction) []byte {

	//On utilise la fonction json.Marshal du package Json afin d'encoder notre attraction
	bytes_s1, err := json.Marshal(s1)
	//gestion des erreur
	if err != nil {
		panic(err)
	}
	return bytes_s1
}

//Permet de decoder une attraction au format Json à l'aide de la fonction json.Unmarshal
//Prend en param un Json
//Renvois une Attraction
func decode(s1 []byte) Attraction {

	var stu = &Attraction{}

	//On utilise la fonction json.Unmarshal du package Json afin de decoder notre json
	var err = json.Unmarshal(s1, stu)
	//gestion des erreur
	if err != nil {
		panic(err)
	}

	fmt.Printf("Unmarshal: Name: %s, InPark: %s ,Place : %s  ,  Manufacturer : %s \n", stu.Name, stu.InPark, stu.Place, stu.Manufacturer)
	return *stu
}

//Fonction qui permet:
//-De gerer les url (la redirection)
//-De crée notre tabeau d'attraction dans lequel on stock nos structure
//-De lancer notre serveur
//Pour chaque routes on appel une fonction spécifique cree préalablement
func router() {

	var size = 0
	//On crée notre tableau à taille variable
	attractions := make([]Attraction, size)

	//Route qui permet l'afficahge de notre html Index
	http.Handle("/", http.FileServer(http.Dir("../static")))

	//Route qui permet de récuperer la liste des attractions
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		//Permet d'autoriser le cross domaine (Ajouté afin de gerer certaine erreur du coter front)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(attractions)

	})

	//Route qui permet de supprimer , récup l'id du formulaire pui appel notre fonction delete()
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		var id = r.FormValue("Id")
		attractions = delete(id, attractions)
	})

	//Route qui permet de crée une attraction , récup les info du formulaire puis appel notre fonction create()
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("creation en cours ... : ")
		name := r.FormValue("Name")
		place := r.FormValue("Place")
		manufacturer := r.FormValue("Manufacturer")
		inPark := r.FormValue("InPark")

		attractions = create(attractions, name, inPark, place, manufacturer)
		fmt.Println(attractions)

	})

	//Route qui permet d'exporter nos données , appel notre fonction saveInFIle()
	http.HandleFunc("/saveInfile", func(w http.ResponseWriter, r *http.Request) {
		saveInfile(attractions)

	})

	//Route qui permet de modifier une attraction , récup les info du formulaire puis appel notre fonction patch()
	http.HandleFunc("/patch", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("update en cours ... : ")
		id := r.FormValue("Id")
		name := r.FormValue("Name")
		place := r.FormValue("Place")
		manufacturer := r.FormValue("Manufacturer")
		inPark := r.FormValue("InPark")

		patch(id, attractions, name, inPark, place, manufacturer)
		fmt.Println(attractions)

		//	patch(id, attractions, name, inPark, manufacturer)

	})
	http.ListenAndServe(":8001", nil)
}

func main() {

	router()

}
