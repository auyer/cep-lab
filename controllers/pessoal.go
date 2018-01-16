package controllers

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"pessoalAPI-gingonic/db"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	error_reason string `json:"error_reason"`
}

type Servidor struct {
	ID                 int    `db:"id, primarykey, autoincrement" json:"id"`
	Siape              int    `db:"siape" json:"siape"`
	Id_pessoa          int    `db:"id_pessoa" json:"id_pessoa"`
	Nome               string `db:"nome" json:"nome"`
	Matricula_interna  int    `db:"matricula_interna" json:"matricula_interna"`
	Nome_identificacao string `db:"nome_identificacao" json:"nome_identificacao"`
	Data_nascimento    string `db:"data_nascimento" json:"data_nascimento"`
	Sexo               string `db:"sexo" json:"sexo"`
}

type ServidorTextual struct {
	ID                 string `json:"id"`
	Siape              string `json:"siape"`
	Id_pessoa          string `json:"id_pessoa"`
	Nome               string `json:"nome"`
	Matricula_interna  string `json:"matricula_interna"`
	Nome_identificacao string `json:"nome_identificacao"`
	Data_nascimento    string `json:"data_nascimento"`
	Sexo               string `json:"sexo"`
}

type ServidorController struct{} // THis is used to make functions callable from ServidorCOntroller

func (ctrl ServidorController) GetServidores(c *gin.Context) {
	q := `select s.id_servidor, s.siape, s.id_pessoa, s.matricula_interna, s.nome_identificacao,
		p.nome, p.data_nascimento, p.sexo from rh.servidor s
	inner join comum.pessoa p on (s.id_pessoa = p.id_pessoa)` //Manual query

	rows, err := db.GetDB().Query(q) //Get Database cursor from DB module
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //will close DB after function GetServidor is over.

	var servidores []Servidor

	var id, id_pessoa, siape, matricula_interna int
	var nome, nome_identificacao, data_nascimento, sexo string
	for rows.Next() {
		err := rows.Scan(&id, &siape, &id_pessoa, &matricula_interna, &nome_identificacao, &nome, &data_nascimento, &sexo)
		if err != nil {
			log.Fatal(err)
		}
		// log.Println(id)
		servidores = append(servidores, Servidor{
			ID:                 id,
			Siape:              siape,
			Id_pessoa:          id_pessoa,
			Nome:               nome,
			Matricula_interna:  matricula_interna,
			Nome_identificacao: nome_identificacao,
			Data_nascimento:    data_nascimento,
			Sexo:               sexo,
		})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	c.JSON(200, servidores)

	if err != nil {
		return
	}

	return
}

func (ctrl ServidorController) GetServidorMat(c *gin.Context) {
	mat := c.Param("matricula") // URL parameter
	// Data security checking to be insterted here
	r, _ := regexp.Compile(`\b[0-9]+\b`)
	if !r.MatchString(mat) {
		log.Println("REGEX CAUGHT AN ERROR")
		c.JSON(404, nil)
		return
	}

	q := fmt.Sprintf(`select s.id_servidor, s.siape, s.id_pessoa, s.matricula_interna, s.nome_identificacao,
		p.nome, p.data_nascimento, p.sexo from rh.servidor s
	inner join comum.pessoa p on (s.id_pessoa = p.id_pessoa) where s.matricula_interna = %s`, mat) //String formating

	rows, err := db.GetDB().Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var servidores []Servidor

	var id, id_pessoa, siape, matricula_interna int
	var nome, nome_identificacao, data_nascimento, sexo string
	for rows.Next() {
		err := rows.Scan(&id, &siape, &id_pessoa, &matricula_interna, &nome_identificacao, &nome, &data_nascimento, &sexo)
		if err != nil {
			log.Fatal(err)
		}

		servidores = append(servidores, Servidor{
			ID:                 id,
			Siape:              siape,
			Id_pessoa:          id_pessoa,
			Nome:               nome,
			Matricula_interna:  matricula_interna,
			Nome_identificacao: nome_identificacao,
			Data_nascimento:    data_nascimento,
			Sexo:               sexo,
		})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	c.JSON(200, servidores)

	if err != nil {
		return
	}

	return
}

func (ctrl ServidorController) PostServidor(c *gin.Context) {
	var serTex ServidorTextual
	regexcheck := false
	var reasons []ErrorBody
	err := c.ShouldBindJSON(&serTex)
	if err != nil {
		log.Panic("|TEXTUAL BINDING|")
		return
	}

	timestamp := time.Now().Unix()
	rand.Seed(timestamp)
	log.Print(sha256.Sum256([]byte(fmt.Sprintf(string(serTex.Nome), string(timestamp)))))
	log.Print(timestamp)

	// REGEX CHEKING PHASE
	r, _ := regexp.Compile(`^(19[0-9]{2}|2[0-9]{3})-(0[1-9]|1[012])-([123]0|[012][1-9]|31)T([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])Z$`)
	if !r.MatchString(serTex.Data_nascimento) {
		regexcheck = true
		reasons = append(reasons, ErrorBody{
			error_reason: "[data_nascimento] failed to match standards. It should look like this: 1969-02-12T00:00:00Z",
		})
	}
	r, _ = regexp.Compile(`^([A-Z][a-z]+([ ]?[a-z]?['-]?[A-Z][a-z]+)*)$`)
	if !r.MatchString(serTex.Nome) {
		regexcheck = true
		reasons = append(reasons, ErrorBody{
			error_reason: "[nome] failed to match standards. It should look like this: Firstname Middlename(optional) Lastname",
		})
	}
	r, _ = regexp.Compile(`^([A-Z][a-z]+([ ]?[a-z]?['-]?[A-Z][a-z]+)*)$`)
	if !r.MatchString(serTex.Nome_identificacao) {
		regexcheck = true
		reasons = append(reasons, ErrorBody{
			error_reason: "[nome_identificacao] failed to match standards. It should look like this: Firstname Middlename(optional) Lastname",
		})
	}
	r, _ = regexp.Compile(`\b[MF]{1}\b`)
	if !r.MatchString(serTex.Sexo) {
		regexcheck = true
		reasons = append(reasons, ErrorBody{
			error_reason: "[sexo] failed to match standards. It should look like this: M for male, F for female",
		})
	}
	r, _ = regexp.Compile(`\b[0-9]+\b`)
	if !r.MatchString(serTex.Siape) {
		regexcheck = true
		reasons = append(reasons, ErrorBody{
			error_reason: "[siape] failed to match standards. It should be only numeric.",
		})
	}
	r, _ = regexp.Compile(`\b[0-9]+\b`)
	if !r.MatchString(serTex.Id_pessoa) {
		regexcheck = true
		reasons = append(reasons, ErrorBody{
			error_reason: "[id_pessoa] failed to match standards. It should be only numeric.",
		})
	}
	if regexcheck {
		c.JSON(400, reasons)
		return
	}
	// END OF REGEX CHEKING PHASE

	// ser.Data_nascimento = serTex.Data_nascimento
	// ser.ID, _ = strconv.Atoi(serTex.ID)

	q := fmt.Sprintf(`
		INSERT INTO rh.servidor_tmp(
			nome, nome_identificacao, id_servidor, siape, id_pessoa, matricula_interna, id_foto,
			data_nascimento, sexo)
			VALUES ('%s', '%s', %d, %d, %d, %d, null, '%s', '%s');
			`, serTex.Nome, serTex.Nome_identificacao, serTex.ID, serTex.Siape, serTex.Id_pessoa, serTex.Matricula_interna,
		serTex.Data_nascimento, serTex.Sexo) //String formating

	rows, err := db.GetDB().Query(q)
	if err != nil {
		log.Panic("|DATABASE ERROR|")
		log.Print(err)
	}

	defer rows.Close()

	var pessoas []Servidor

	c.JSON(200, pessoas)

	// if err != nil {
	// 	return
	// }

	return
}
