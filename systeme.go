// Réalisé par Léo et Nacer

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// Nous avons utilisé une structure pour permettre de stocker nos processus
// Cette structure nous permet d'optimiser le code
type Process struct {
	Pid string
	Cwd string
	Exe string
}

func main() {
	fmt.Println(pid())
}

//Fonction qui permet de récup l'ensemble de nos PID
//LA fonction vas ouvre le dossier /proc afin de récuperer l'ensemble des dossier correspondant a des PID
//Pour chaque sous dossier on appel nos deux fonction inspectPIDCWD et inspectPidEXE
//Puis nous affichons le tout à l'aide de notre fonction display()
func pid() string {
	var process []Process
	f, err := os.Open("/proc")

	//On gere le cas ou le dossier /proc n'est pas present
	if err != nil {
		return "Error 500 , you dont have directory /proc in your linux architecture"
	}

	//Récupère les sous dossier du dossier ./proc
	fileInfo, err := f.Readdir(-1)
	f.Close()
	//On gere le cas ou t'utilisateur n'as pas la permission
	if err != nil {
		return "Error 500, permission denied, cant open /proc"
	}

	//On boucle sur l'ensemble de nos sous-dossier
	for _, file := range fileInfo {

		//On verifie qu'il s'agit bien de chiffre (chiffre = pid)
		//Si oui on appel nos fonction
		if _, err := strconv.Atoi(file.Name()); err == nil {
			var p Process
			p.Pid = file.Name()
			p.Cwd = inspectPidCwd(file.Name())
			p.Exe = inspectPidExe(file.Name())
			process = append(process, p)

		}
	}

	return display(process)

}

// inspecte le dossier lié au processus pour chercher si le dossier cwd
// existe et si oui retourne son path (chemin absolue)
func inspectPidCwd(id string) string {

	//On utilise la fonction ReadLink de la librairie OS afin de de retourner lien symbolique
	file, err := os.Readlink("/proc/" + id + "/cwd")

	if err == nil {
		//On utilise la fonction abs() de la librairie filepath afin de récuperer le chemin absolue
		path, _ := filepath.Abs(file)
		return path + "/proc/" + id + "/cwd"
	} else {
		return "no cwd"
	}

}

// inspecte le dossier lié au processus pour chercher si le fichier exe
// existe et si oui retourne son path
func inspectPidExe(id string) string {

	//On utilise la fonction ReadLink de la librairie OS afin de de retourner lien symbolique
	file, err := os.Readlink("/proc/" + id + "/exe")
	if err == nil {

		//On utilise la fonction abs() de la librairie filepath afin de récuperer le chemin absolue
		path, _ := filepath.Abs(file)
		return path + "/proc/" + id + "/exe\n"
	} else {
		return "no exe\n"
	}
}

// fonction d'affichage  avec pour chaque processus , son pid , le chemin s'il existe de cwd et exe
func display(process []Process) string {
	var show string
	for i := 0; i < len(process); i++ {
		pid := "PID : " + process[i].Pid
		cwd := "CWD : " + process[i].Cwd
		exe := "EXE : " + process[i].Exe
		show = show + pid + " |  " + cwd + "  | " + exe + "\n"

	}
	return show

}
