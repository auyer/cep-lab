package controllers

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/latitude-RESTsec-lab/api-gingonic/db"

	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Reason string `json:"reason"`
}

type Servidor struct {
	ID                int    `db:"id, primarykey, autoincrement" json:"id"`
	Siape             int    `db:"siape" json:"siape"`
	Idpessoa          int    `db:"id_pessoa" json:"id_pessoa"`
	Nome              string `db:"nome" json:"nome"`
	Matriculainterna  int    `db:"matricula_interna" json:"matricula_interna"`
	Nomeidentificacao string `db:"nome_identificacao" json:"nome_identificacao"`
	Datanascimento    string `db:"data_nascimento" json:"data_nascimento"`
	Sexo              string `db:"sexo" json:"sexo"`
}

type ServidorController struct{} // THis is used to make functions callable from ServidorCOntroller

func (ctrl ServidorController) GetServidores(c *gin.Context) {
	q := `select s.id_servidor, s.siape, s.id_pessoa, s.matricula_interna, s.nome_identificacao,
		p.nome, p.data_nascimento, p.sexo from rh.servidor s
	inner join comum.pessoa p on (s.id_pessoa = p.id_pessoa)` //Manual query

	rows, err := db.GetDB().Query(q) //Get Database cursor from DB module
	if err != nil {
		log.Println(err)
		c.JSON(400, ErrorBody{
			Reason: err.Error(),
		})
		return
	}
	defer rows.Close() //will close DB after function GetServidor is over.

	var servidores []Servidor

	var id, idpessoa, siape, matriculainterna int
	var nome, nomeidentificacao, datanascimento, sexo string
	for rows.Next() {
		err := rows.Scan(&id, &siape, &idpessoa, &matriculainterna, &nomeidentificacao, &nome, &datanascimento, &sexo)
		if err != nil {
			log.Println(err)
			c.JSON(400, ErrorBody{
				Reason: err.Error(),
			})
			return
		}
		// log.Println(id)
		date, _ := time.Parse("1969-02-12", datanascimento)
		servidores = append(servidores, Servidor{
			ID:                id,
			Siape:             siape,
			Idpessoa:          idpessoa,
			Nome:              nome,
			Matriculainterna:  matriculainterna,
			Nomeidentificacao: nomeidentificacao,
			Datanascimento:    date.Format("1969-02-12"),
			Sexo:              sexo,
		})
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		c.JSON(400, ErrorBody{
			Reason: err.Error(),
		})
		return
	}
	c.JSON(200, servidores)

	if err != nil {
		log.Println(err)
		c.JSON(400, ErrorBody{
			Reason: err.Error(),
		})
		return
	}

	return
}

func (ctrl ServidorController) GetServidorMat(c *gin.Context) {
	mat := c.Param("matricula") // URL parameter
	// Data security checking to be insterted here
	r, err := regexp.Compile(`\b[0-9]+\b`)
	if !r.MatchString(mat) {
		log.Println(err)
		c.JSON(404, ErrorBody{
			Reason: err.Error(),
		})
		return
	}

	q := fmt.Sprintf(`select s.id_servidor, s.siape, s.id_pessoa, s.matricula_interna, s.nome_identificacao,
		p.nome, p.data_nascimento, p.sexo from rh.servidor s
	inner join comum.pessoa p on (s.id_pessoa = p.id_pessoa) where s.matricula_interna = %s`, mat) //String formating

	rows, err := db.GetDB().Query(q)
	if err != nil {
		log.Println(err)
		c.JSON(400, ErrorBody{
			Reason: err.Error(),
		})
		return
	}
	defer rows.Close()

	var servidores []Servidor

	var id, idpessoa, siape, matriculainterna int
	var nome, nomeidentificacao, datanascimento, sexo string
	for rows.Next() {
		err := rows.Scan(&id, &siape, &idpessoa, &matriculainterna, &nomeidentificacao, &nome, &datanascimento, &sexo)
		if err != nil {
			log.Println(err)
			c.JSON(400, ErrorBody{
				Reason: err.Error(),
			})
			return
		}

		date, _ := time.Parse("1969-02-12", datanascimento)
		servidores = append(servidores, Servidor{
			ID:                id,
			Siape:             siape,
			Idpessoa:          idpessoa,
			Nome:              nome,
			Matriculainterna:  matriculainterna,
			Nomeidentificacao: nomeidentificacao,
			Datanascimento:    date.Format("1969-02-12"),
			Sexo:              sexo,
		})
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		c.JSON(400, ErrorBody{
			Reason: err.Error(),
		})
		return
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
		log.Println(err)
		c.JSON(400, ErrorBody{
			Reason: err.Error(),
		})
		return
	}

	// REGEX CHEKING PHASE

	r, _ := regexp.Compile(`^(19[0-9]{2}|2[0-9]{3})-(0[1-9]|1[012])-([123]0|[012][1-9]|31)$`)
	if (!r.MatchString(ser.Datanascimento)){
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[data_nascimento] failed to match API requirements. It should look like this: 1969-02-12",
		})
	}else{
		now := time.Now()	
		now.Format(time.RFC3339)
		time, err := time.Parse("2006-01-02 15:04:05 -0200",fmt.Sprint(ser.Datanascimento," 00:00:00 -0200"))
		if err != nil {
			log.Println(err)
			c.JSON(500, ErrorBody{
				Reason: err.Error(),
			})
		}
		if(!now.After(time)){
			regexcheck = true
			Reasons = append(Reasons, ErrorBody{
				Reason: "[data_nascimento] failed to match API requirements. It should not be in future",
		})
		}
	}
	r, _ = regexp.Compile(`^([A-Z][a-z]+([ ]?[a-z]?['-]?[A-Z][a-z]+)*)$`)
	if !r.MatchString(ser.Nome) {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[nome] failed to match API requirements. It should look like this: Firstname Middlename(optional) Lastname",
		})
	}
	r, _ = regexp.Compile(`^([A-Z][a-z]+([ ]?[a-z]?['-]?[A-Z][a-z]+)*)$`)
	if !r.MatchString(ser.Nomeidentificacao) {
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
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05-0700")
	b := md5.Sum([]byte(fmt.Sprintf(string(ser.Nome), string(timestamp))))
	bid := binary.BigEndian.Uint64(b[:])

	q := fmt.Sprintf(`
		INSERT INTO rh.servidor_tmp(
			nome, nome_identificacao, siape, id_pessoa, matricula_interna, id_foto,
			data_nascimento, sexo)
			VALUES ('%s', '%s', %d, %d, %d, null, '%s', '%s');
			`, ser.Nome, ser.Nomeidentificacao, ser.Siape, ser.Idpessoa, bid%99999,
		ser.Datanascimento, ser.Sexo) //String formating

	rows, err := db.GetDB().Query(q)
	if err != nil {
		log.Println(err)
		c.JSON(400, ErrorBody{
			Reason: err.Error(),
		})
		return
	}

	defer rows.Close()

	var pessoas []Servidor

	c.JSON(201, pessoas)

	// if err != nil {
	// 	return
	// }

	return
}
