package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	years := flag.String("years", "", "years for the graph")
	region := flag.String("region", "", "region for the graph")
	age := flag.String("age", "", "age for the graph")
	flag.Parse()

	fileData := "data/covid-hosp-age.csv"
	fileGraph := "graph/hospitalisations_graph.png"

	// Init Strings
	stringYears := ""
	if *years != "" {
		stringYears = "en " + *years
	}
	titleGraph := "Admissions hospitalières " + stringYears + "\n dans la région " + *region + "\npour les personnes agées de " + *age + " ans"
	labelY := "Nouvelles Admissions Hospitalières"
	labelX := "Années"

	// Open File
	lines, err := openFile(fileData)
	if err != nil {
		log.Fatalf("erreur lors de l'ouverture du fichier: %s", err.Error())
	}

	// Vérifier les arguments
	if *years != "" && len(*years) != 4 {
		log.Fatalf("Veuillez spécifier 'years' au format YYYY (ex : 2020)")
	}
	if *region == "" || len(*region) != 2 {
		log.Fatalf("Veuillez spécifier 'region' au format RR (ex : 01)")
	}
	if *age == "" || len(*age) != 2 {
		log.Fatalf("Veuillez spécifier 'age' au format AA (ex : 09)")
	}

	// Initialiser les slices pour stocker les données
	var semaines []string
	var hospitalisations []int

	// Analyser les lignes pour extraire les données
	for _, line := range lines[1:] { // Ignorer l'en-tête
		value := strings.Split(line[0], ";")
		if strings.Contains(value[0], *region) &&
			strings.Contains(value[2], *age) &&
			(*years == "" || strings.Contains(value[1], *years)) { // Filtre

			semaine := value[1]
			if err != nil {
				log.Fatalf("Erreur string to int: %s", err.Error())
			}
			hospitalisation, err := strconv.Atoi(value[3])
			if err != nil {
				log.Fatalf("Erreur string to int: %s", err.Error())
			}

			// Ajouter les valeurs aux slices correspondantes
			semaines = append(semaines, semaine)
			hospitalisations = append(hospitalisations, hospitalisation)
		}
	}
	if len(semaines) == 0 {
		if *years == "" {
			log.Fatalf("Aucune donnée pour la region %s et l'age %s", *region, *age)
		} else {
			log.Fatalf("Aucune donnée pour l'année %s, la region %s et l'age %s", *years, *region, *age)
		}
	}

	// Créer le deuxième graphique (Semaine vs Nouvelles Admissions Hospitalières)
	err = createPlot(titleGraph, semaines, hospitalisations, labelX, labelY, fileGraph)
	if err != nil {
		log.Fatalf("Erreur lors de la création du graphique des hospitalisations: %s", err.Error())
	}

	fmt.Println("Graphiques créés avec succès.")
}

func openFile(fileName string) ([][]string, error) {
	// Ouvrir le fichier CSV
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Lire les données CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func createPlot(title string, xData []string, yData []int, xLabel string, yLabel string, fileName string) error {
	// Créer un nouveau tracé
	p := plot.New()

	// Définir le titre et les étiquettes des axes
	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel

	// Créer un plotter pour les données
	points := make(plotter.XYs, len(xData))
	for i, x := range xData {
		// Convertir les semaines en numéros de semaine
		weekParts := strings.Split(x, "-S")
		year, _ := strconv.Atoi(weekParts[0])
		weekNumber, _ := strconv.Atoi(weekParts[1])

		points[i].X = float64(year) + float64(weekNumber)/52.0 // Combinez l'année et le numéro de semaine
		points[i].Y = float64(yData[i])
	}

	// Ajouter les données au tracé
	line, _, pointsErr := plotter.NewLinePoints(points)
	if pointsErr != nil {
		return pointsErr
	}
	p.Add(line)

	// Enregistrer le graphique dans un fichier PNG
	return savePlot(p, 6*vg.Inch, 4*vg.Inch, fileName)
}

func savePlot(p *plot.Plot, width, height vg.Length, fileName string) error {
	// Enregistrer le graphique dans un fichier PNG
	if saveErr := p.Save(width, height, fileName); saveErr != nil {
		return saveErr
	}

	return nil
}
