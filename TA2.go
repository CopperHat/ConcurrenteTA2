package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
)

type Option struct {
	Value, Id, Text string
	Selected        bool
}

const HTML = `
<!DOCTYPE html>
<html lang="en">
     <head>
        <meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">

        <title>selected attribute</title>
<link rel="stylesheet" type="text/css" href="fondo.css">
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
<link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Lobster">
<link rel="stylesheet" href="https://www.w3schools.com/w3css/4/w3.css">
    </head>



    <body style=' margin: 0;
    padding: 0;
	background: url("https://www.entornointeligente.com/wp-content/uploads/2020/02/entornointeligente_roberto_pocaterra_pocaterra_buenos_aires_lo_que_tienes_que_saber_sobre_las_nuevas_cifras_del_coronavirus.jpg");    
	background-size: cover;
	background-position: center;
	background-repeat: no-repeat;
    font-family: sans-serif;'>

<div class="w3-container w3-lobster">                                                 
        <h1 style='font-family: "Lobster", serif;  position: absolute; top: 50px; left: 670px;'>Bienvenido al Analizador Covid</h1>
</div>

<div class="w3-container w3-lobster">                                                 
		<h1 style='font-family: "Lobster", serif; color: #05687; position: absolute; top: 110px; left: 100px;'>Por favor</h1>
		<h1 style='font-family: "Lobster", serif; color: #05687; position: absolute; top: 170px; left: 30px;'>tomate tu tiempo y </h1>
		<h1 style='font-family: "Lobster", serif; color: #05687; position: absolute; top: 230px; left: 100px;'>selecciona </h1>
		<h1 style='font-family: "Lobster", serif; color: #05687; position: absolute; top: 300px; left: 30px;'>la opcion correcta </h1>
		<h1 style='font-family: "Lobster", serif; color: #05687; position: absolute; top: 370px; left: 42px;'>para asegurar </h1>
		<h1 style='font-family: "Lobster", serif; color: #05687; position: absolute; top: 440px; left: 35px;'>un buen resultado </h1>
	
</div>


<div class="flip-card style = left: 30px;'">
  <div class="flip-card-inner">
    <div class="flip-card-front">
      <img src="https://i.ytimg.com/vi/mA1qCnk4Lg4/hqdefault.jpg" alt="" width="330" height="600">
    </div>
  </div>
</div>

        <form method="GET">                                                                                                   
			<label>UBIGEO:</label>
			<select id="UBIGEO" name="UBIGEO">
				{{range .}}
				<option value="{{.Value}}" id="{{.Id}}" {{if .Selected}}selected{{end}}>{{.Text}}</option>
				{{end}}
			</select>
			<label>DEPARTAMENTO:</label>
			<select id="DEPARTAMENTO" name="DEPARTAMENTO">
				{{range .}}
				<option value="{{.Value}}" id="{{.Id}}" {{if .Selected}}selected{{end}}>{{.Text}}</option>
				{{end}}
			</select>
			<label>PROVINCIA:</label>
			<select id="PROVINCIA" name="PROVINCIA">
				{{range .}}
				<option value="{{.Value}}" id="{{.Id}}" {{if .Selected}}selected{{end}}>{{.Text}}</option>
				{{end}}
			</select>
			<label>DISTRITO:</label>
			<select id="DISTRITO" name="DISTRITO">
				{{range .}}
				<option value="{{.Value}}" id="{{.Id}}" {{if .Selected}}selected{{end}}>{{.Text}}</option>
				{{end}}
			</select>
			<label>TIPO_ELECCION:</label>
			<select id="TIPO_ELECCION" name="TIPO_ELECCION">
				{{range .}}
				<option value="{{.Value}}" id="{{.Id}}" {{if .Selected}}selected{{end}}>{{.Text}}</option>
				{{end}}
			</select>
			<label>MESA_DE_VOTACION:</label>
			<select id="MESA_DE_VOTACION" name="MESA_DE_VOTACION">
				{{range .}}
				<option value="{{.Value}}" id="{{.Id}}" {{if .Selected}}selected{{end}}>{{.Text}}</option>
				{{end}}
			</select>
			<label>DESCRIP_ESTADO_ACTA:</label>
			<select id="DESCRIP_ESTADO_ACTA" name="DESCRIP_ESTADO_ACTA">
				{{range .}}
				<option value="{{.Value}}" id="{{.Id}}" {{if .Selected}}selected{{end}}>{{.Text}}</option>
				{{end}}
			</select>
			<label>TIPO_OBSERVACION:</label>
			<select id="TIPO_OBSERVACION" name="TIPO_OBSERVACION">
				{{range .}}
				<option value="{{.Value}}" id="{{.Id}}" {{if .Selected}}selected{{end}}>{{.Text}}</option>
				{{end}}
			</select>
			<label>N_CVAS:</label>
			<select id="N_CVAS" name="N_CVAS">
				{{range .}}
				<option value="{{.Value}}" id="{{.Id}}" {{if .Selected}}selected{{end}}>{{.Text}}</option>
				{{end}}
			</select>
			<label>N_ELEC_HABIL:</label>
			<select id="N_ELEC_HABIL" name="N_ELEC_HABIL">
				{{range .}}
				<option value="{{.Value}}" id="{{.Id}}" {{if .Selected}}selected{{end}}>{{.Text}}</option>
				{{end}}
			</select>
            </div>
           <div style="position: absolute;
  top: 500px;
  left: 850px;
  width: 300px;
  height: 200px;">
            <input style='display: inline-block; padding: 15px 25px; font-weight:  bolder;  font-size: 24px; cursor: pointer; text-align: center; text-decoration: none; outline: none; color: black;
                   background-color: #009C8C; border: none; border-radius: 15px; box-shadow: 0 9px #999;' type="submit" value="Analizar" align="center" name = "submit">
           </div>
        </form>

    </body>
</html>
`

var placesPageTmpl *template.Template = template.Must(template.New("	").Parse(HTML))

const localAddr = "192.168.0.9:8000"

const (
	cnum = iota
	opContagiado
	opNoContagiado
)

var chInfo chan map[string]int

type Registro struct {
	UBIGEO              string
	DEPARTAMENTO        string
	PROVINCIA           string
	DISTRITO            string
	TIPO_ELECCION       string
	MESA_DE_VOTACION    string
	DESCRIP_ESTADO_ACTA string
	TIPO_OBSERVACION    string
	N_CVAS              string
	N_ELEC_HABIL        string
}

type estadoRegistro struct {
	Code int
	Addr string
	Op   int
}

var addrs = []string{
	"192.168.0.27:8000",
	"192.168.0.28:8000"
}

func main() {
	url := "https://raw.githubusercontent.com/mledoze/countries/master/dist/countries.csv"
	data, err := readCSVFromUrl(url)
	if err != nil {
		panic(err)
	}

	for idx, row := range data {
		// skip header
		if idx == 0 {
			continue
		}

		if idx == 6 {
			break
		}

		fmt.Println(row[2])
	}
	fmt.Print(addrs)
	fmt.Println()
	http.HandleFunc("/", name)
	http.ListenAndServe(":8080", nil)
}

func readCSVFromUrl(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ';'
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}


func name(w http.ResponseWriter, r *http.Request) {
	var registro = Registro{}

	registro.UBIGEO = r.FormValue("UBIGEO")
	registro.DEPARTAMENTO = r.FormValue("DEPARTAMENTO")
	registro.PROVINCIA = r.FormValue("PROVINCIA")
	registro.DISTRITO = r.FormValue("DISTRITO")
	registro.TIPO_ELECCION = r.FormValue("TIPO_ELECCION")
	registro.MESA_DE_VOTACION = r.FormValue("MESA_DE_VOTACION")
	registro.DESCRIP_ESTADO_ACTA = r.FormValue("DESCRIP_ESTADO_ACTA")
	registro.TIPO_OBSERVACION = r.FormValue("TIPO_OBSERVACION")
	registro.N_CVAS = r.FormValue("N_CVAS")
	registro.N_ELEC_HABIL = r.FormValue("N_ELEC_HABIL")
}

func handle(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	var msg estadoRegistro
	if err := dec.Decode(&msg); err != nil {
		log.Println("Can't decode from", conn.RemoteAddr())
	} else {
		fmt.Println(msg)
		switch msg.Code {
		case cnum:
			concensus(conn, msg)
		}
	}
}

func concensus(conn net.Conn, msg estadoRegistro) {
	info := <-chInfo
	info[msg.Addr] = msg.Op

	go func() { chInfo <- info }()
}
func send(remoteAddr string, msg estadoRegistro) {
	if conn, err := net.Dial("tcp", remoteAddr); err != nil {
		log.Println("Can't dail", remoteAddr)
	} else {
		defer conn.Close()
		fmt.Println("Sending to ", remoteAddr)
		enc := json.NewEncoder(conn)
		enc.Encode(msg)
	}
}
