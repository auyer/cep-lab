package controllers

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"log"
	"pessoalAPI-gingonic/db"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Reason string `json:"reason"`
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
	regexcheck := false
	var ser Servidor
	var Reasons []ErrorBody
	err := c.ShouldBindJSON(&ser)
	if err != nil {
		log.Println("BINDING ERROR")
		c.JSON(400, ErrorBody{
			Reason: "Wrong Datatype",
		})
		return
	}

	// REGEX CHEKING PHASE
	r, _ := regexp.Compile(`^(19[0-9]{2}|2[0-9]{3})-(0[1-9]|1[012])-([123]0|[012][1-9]|31)T([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])Z$`)
	if !r.MatchString(ser.Data_nascimento) {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[data_nascimento] failed to match API requirements. It should look like this: 1969-02-12T00:00:00Z",
		})
	}
	r, _ = regexp.Compile(`^([A-Z][a-z]+([ ]?[a-z]?['-]?[A-Z][a-z]+)*)$`)
	if !r.MatchString(ser.Nome) {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[nome] failed to match API requirements. It should look like this: Firstname Middlename(optional) Lastname",
		})
	}
	r, _ = regexp.Compile(`^([A-Z][a-z]+([ ]?[a-z]?['-]?[A-Z][a-z]+)*)$`)
	if !r.MatchString(ser.Nome_identificacao) {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: string("[nome_identificacao] failed to match API requirements. It should look like this: Firstname Middlename(optional) Lastname"),
		})
	}
	r, _ = regexp.Compile(`\b[MF]{1}\b`)
	if !r.MatchString(ser.Sexo) {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[sexo] failed to match API requirements. It should look like this: M for male, F for female",
		})
	}
	// r, _ = regexp.Compile(`\b[0-9]+\b`)
	// if !r.MatchString(strconv.Itoa(ser.Siape)) {
	// 	regexcheck = true
	// 	Reasons = append(Reasons, ErrorBody{
	// 		error_reason: "[siape] failed to match API requirements. It should be only numeric.",
	// 	})
	// }
	// r, _ = regexp.Compile(`\b[0-9]+\b`)
	// if !r.MatchString(strconv.Itoa(ser.Id_pessoa)) {
	// 	regexcheck = true
	// 	Reasons = append(Reasons, ErrorBody{
	// 		error_reason: "[id_pessoa] failed to match API requirements. It should be only numeric.",
	// 	})
	// }
	if regexcheck {
		c.JSON(400, Reasons)
		return
	}
	// END OF REGEX CHEKING PHASE

	timestamp := time.Now().Unix()

	b := md5.Sum([]byte(fmt.Sprintf(string(ser.Nome), string(timestamp))))
	bid := binary.BigEndian.Uint64(b[:])
	// log.Println(strconv.Atoi(string(b[:])))
	// ser.Data_nascimento = serTex.Data_nascimento
	// ser.ID, _ = strconv.Atoi(serTex.ID)

	q := fmt.Sprintf(`
		INSERT INTO rh.servidor_tmp(
			nome, nome_identificacao, siape, id_pessoa, matricula_interna, id_foto,
			data_nascimento, sexo)
			VALUES ('%s', '%s', %d, %d, %d, null, '%s', '%s');
			`, ser.Nome, ser.Nome_identificacao, ser.Siape, ser.Id_pessoa, bid%9999,
		ser.Data_nascimento, ser.Sexo) //String formating

	rows, err := db.GetDB().Query(q)
	if err != nil {
		log.Print("|DATABASE ERROR|")
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
