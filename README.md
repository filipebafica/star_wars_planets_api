# star_wars_planets_api
An API to handle Star Wars planets

## 🗂 Table of Contents
* [About](#-about)
* [Getting Started](#-getting-started)
* [How to Use](#-how-to-use)
* [Testing](#-testing)

## 🧐 About
This is an API to handle Star Wars Planets.\
It was used: Golang, Docker and MongoDB.

## 🏁 Getting Started
You need to have Docker and Docker Compose (v 3.9 was used to this projetc, you may change on `docker-compose.yml` first line).

#### ⚙️ Installing
To compile the code, clone the repo, and run the following commands.
```
$ git clone https://github.com/filipebafica/star_wars_planets_api.git
$ cd star_wars_planets_api
$ docker compose build
$ docker compose up -d
```

## 🎈 How to Use
The following endpoints are available on localhost:8000\
To insert a planet: `[POST] /v1/planeta`\
The `nome`, `clima` and `terreno` fields are mandatory. `_id` and `filmes` will be added automatically. 
```
{
  "nome": "Coruscant",
  "clima": "arido",
  "terreno": "deserto"
}
```
To retrieve all planets: `[GET] /v1/planetas`
```
{
  "_id": "63109f997c65950feadbb7c1",
  "nome": "Coruscant",
  "clima": "temperado",
  "terreno": "urbano",
  "filmes": 4
},
{
  "_id": "6310a3d57c65950feadbb7c2",
  "nome": "Tatooine",
  "clima": "arido",
  "terreno": "deserto",
  "filmes": 5
}
```
To query a planet by ID `[GET] /v1/planeta?id=<id>`

To query a planet by name `[GET] /v1/planeta?nome=<name>`

To remove a planet `[DELETE] /v1/planeta?id=<id>`

## 🧪 Testing
To run the unit tests do the following commands.
```
$ docker compose exec app bash
$ cd tests
$ go test -v
```
![image](https://user-images.githubusercontent.com/31427890/188003926-1cae3a94-4690-412b-be7d-dd7bdbc0052f.png)



