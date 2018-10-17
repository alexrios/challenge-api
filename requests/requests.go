package requests

/*
* Mapeamento dos request recebidos pela App
 */

type NewPlanetRequest struct {
	Climate string `json:"climate" valid:"alphanum"`
	Name    string `json:"name" valid:"alphanum"`
	Terrain string `json:"terrain" valid:"alphanum"`
}

type FindByNameRequest struct {
	Name string `json:"name" valid:"alphanum"`
}

type FindByIDRequest struct {
	Id string `json:"id" valid:"alphanum"`
}
