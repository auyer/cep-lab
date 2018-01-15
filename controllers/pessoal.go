package controllers

import (
	"fmt"
	"log"
	"pessoalAPI-gingonic/db"
	"regexp"

	"github.com/gin-gonic/gin"
)

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
	var ser Servidor
	err := c.BindJSON(&ser)
	if err != nil {
		log.Fatal(err)
		return
	}
	q := fmt.Sprintf(`
		INSERT INTO rh.servidor_tmp(
			nome, nome_identificacao, id_servidor, siape, id_pessoa, matricula_interna, id_foto, 
			data_nascimento, sexo)
			VALUES ('%s', '%s', %d, %d, %d, %d, null, '%s', '%s'); 
			`, ser.Nome, ser.Nome_identificacao, ser.ID, ser.Siape, ser.Id_pessoa, ser.Matricula_interna,
		ser.Data_nascimento, ser.Sexo) //String formating

	rows, err := db.GetDB().Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var pessoas []Servidor

	c.JSON(200, pessoas)

	// if err != nil {
	// 	return
	// }

	return
}
