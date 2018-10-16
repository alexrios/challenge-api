package starwarsapi

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
)

type Planet struct {
	Name           string        `json:"name"`
	RotationPeriod string        `json:"rotation_period"`
	OrbitalPeriod  string        `json:"orbital_period"`
	Diameter       string        `json:"diameter"`
	Climate        string        `json:"climate"`
	Gravity        string        `json:"gravity"`
	Terrain        string        `json:"terrain"`
	SurfaceWater   string        `json:"surface_water"`
	Population     string        `json:"population"`
	ResidentURLs   []residentURL `json:"residents"`
	FilmURLs       []filmURL     `json:"films"`
	Created        string        `json:"created"`
	Edited         string        `json:"edited"`
	URL            string        `json:"url"`
}

type PlanetSearchResult struct {
	Count   int      `json:"count"`
	Results []Planet `json:"results"`
}

type filmURL string
type residentURL string

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

//retorna o numero de aparicoes do planeta (receiver)
func (p *Planet) Appearances() int {
	return len(p.FilmURLs)
}

func FetchPlanet(name string) (Planet, error) {
	planet, err := GetPlanet(fmt.Sprintf("https://swapi.co/api/planets/?format=json&search=%v", name))
	if err != nil {
		return Planet{}, err
	}
	if len(planet.Results) > 0 {
		return planet.Results[0], nil;
	} else {
		return Planet{}, errors.New("No planet in response")
	}

}
