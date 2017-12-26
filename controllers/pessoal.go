package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"pessoalAPI/db"

	"github.com/gin-gonic/gin"
)

type Pessoal struct {
	ID                 int    `db:"id, primarykey, autoincrement" json:"id"`
	Siape              string `db:"siape" json:"siape"`
	Id_pessoa          int    `db:"id_pessoa" json:"id_pessoa"`
	Nome               string `db:"nome" json:"nome"`
	Matricula_interna  string `db:"matricula_interna" json:"matricula_interna"`
	Nome_identificacao string `db:"nome_identificacao" json:"nome_identificacao"`
	Data_nascimento    string `db:"data_nascimento" json:"data_nascimento"`
	Sexo               string `db:"sexo" json:"sexo"`
}

type PessoalController struct{}

func (ctrl PessoalController) GetPessoal(c *gin.Context) { // Hello
	//func (ctrl PessoalController) getPessoal(c *gin.Context) (pessoal Pessoal, err error) {
	q := `select s.id_servidor, s.siape, s.id_pessoa, s.matricula_interna, s.nome_identificacao,
		p.nome, p.data_nascimento, p.sexo from rh.servidor s
	inner join comum.pessoa p on (s.id_pessoa = p.id_pessoa)`

	// q2 := "select id_servidor from rh.servidor"

	rows, err := db.GetDB().Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var pessoas []Pessoal
	var id, id_pessoa int
	var siape, nome, matricula_interna, nome_identificacao, data_nascimento, sexo string
	for rows.Next() {
		err := rows.Scan(&id, &siape, &id_pessoa, &matricula_interna, &nome_identificacao, &nome, &data_nascimento, &sexo)
		if err != nil {
			log.Fatal(err)
		}
		// log.Println(id)
		pessoas = append(pessoas, Pessoal{
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
	jsonPessoal, errr := json.Marshal(pessoas)
	if errr != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stdout, "%s", jsonPessoal)
	c.JSON(200, gin.H{
		"message": "eh",
	})
	if err != nil {
		return
	}

	return
}
