package main

// from https://github.com/oladush/json-to-mmdb/blob/main/main.go

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime/debug"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

const (
	FAIL = "\033[91m"
	ENDC = "\033[0m"
)

// returns length of file
func lineCount(filename string) (int64, error) {
	lc := int64(0)
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		lc++
	}
	return lc, s.Err()
}

// returns number is integer
func isInteger(val float64) bool {
	return val == float64(int(val))
}

// save mmdb
func writeMMDB(filename string, tree *mmdbwriter.Tree) {
	fh, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = tree.WriteTo(fh)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// returns prefix from json (prefix is key)
func getPrefix(record *map[string]interface{}) string {
	for k := range *record {
		return k
	}
	return ""
}

// convert json to special mmdb map
func parseToMap(record *map[string]interface{}) mmdbtype.Map {
	root_map := mmdbtype.Map{}
	var obj interface{}
	for k := range *record {
		// check type
		obj = (*record)[k]

		// if string
		if val, ok := obj.(string); ok {
			root_map[mmdbtype.String(k)] = mmdbtype.String(val)

			// if number
		} else if val, ok := obj.(float64); ok {
			if isInteger(val) {
				switch k {
				case "geoname_id":
					root_map[mmdbtype.String(k)] = mmdbtype.Uint32(val)
				case "confidence":
					root_map[mmdbtype.String(k)] = mmdbtype.Uint16(val)
				case "accuracy_radius":
					root_map[mmdbtype.String(k)] = mmdbtype.Uint16(val)
				case "metro_code":
					root_map[mmdbtype.String(k)] = mmdbtype.Uint16(val)
				default:
					root_map[mmdbtype.String(k)] = mmdbtype.Uint32(val)
				}
			} else {
				root_map[mmdbtype.String(k)] = mmdbtype.Float64(val)
			}
			// if bool
		} else if val, ok := obj.(bool); ok {
			root_map[mmdbtype.String(k)] = mmdbtype.Bool(val)
			// if struct
		} else if val, ok := obj.(map[string]interface{}); ok {
			root_map[mmdbtype.String(k)] = parseToMap(&val)

			// if array
		} else if val, ok := obj.([]interface{}); ok {
			arr := mmdbtype.Slice{}
			for _, v := range val {
				if value, ok := v.(string); ok {
					arr = append(arr, mmdbtype.String(value))
				} else if value, ok := v.(float64); ok {
					if isInteger(value) {
						arr = append(arr, mmdbtype.Uint32(value))
					} else {
						arr = append(arr, mmdbtype.Float64(value))
					}
				} else if value, ok := v.(map[string]interface{}); ok {
					arr = append(arr, parseToMap(&value))
				}
				root_map[mmdbtype.String(k)] = arr
			}
		} else {
			// fmt.Printf("Type is %T\n", obj)
		}
	}
	return root_map
}

// create record for prefix and add to tree
func addRecord(prefix string, record *map[string]interface{}, tree *mmdbwriter.Tree) {
	// get prefix like string
	// prefix := getPrefix(record)

	// try parse prefix
	_, prefix_p, err := net.ParseCIDR(prefix)
	if err != nil {
		log.Fatal(err)
	}

	rec_ := record // (*record)[prefix].(map[string]interface{})

	var record_map mmdbtype.Map = parseToMap(rec_)

	tree.InsertFunc(
		prefix_p, inserter.TopLevelMergeWith(record_map))

	//for k := range record_map {
	//	delete(record_map, k)
	//}
}

func main() {
	var input_json string
	var output_mmdb string
	var mmdb_type string

	flag.StringVar(&input_json, "i", "", "input list of json records")
	flag.StringVar(&output_mmdb, "o", "", "output mmdb")
	flag.StringVar(&mmdb_type, "t", "", "mmdb type")

	flag.Usage = func() {
		fmt.Println("┌───────────────────────────────────────────────────────────────────────────────────┐")
		fmt.Println("│   /\\__/\\    jsonToMMDB                                                            │")
		fmt.Println("│  ( ⊙ ‿ ⊙)    simple script for convert json to mmdb                               │")
		fmt.Println("│   (｡  ｡)     usage: ./jsonToMMDB -i input.json -o output.mmdb -t GeoLite2-Country │")
		fmt.Println("└───────────────────────────────────────────────────────────────────────────────────┘")
		fmt.Println("────────────────────────────────────┐")
		flag.PrintDefaults()
		fmt.Println("──────────────────────────────────────────────────────────────────────────────────────────────┐")
		fmt.Println("Input example:")
		fmt.Println("{")
		fmt.Println("\t\"10.0.0.0/16\": { \"country\": {\"geoname_id\": 666, \"names\": {\"en\": \"Zalibobastan\"}}},")
		fmt.Println("\t\"10.0.1.0/24\": { \"country\": {\"geoname_id\": 333, \"names\": {\"en\": \"Kyrgyzstan\"}}}")
		fmt.Println("}")
	}

	flag.Parse()

	if input_json == "" || output_mmdb == "" || mmdb_type == "" {
		flag.Usage()
		return
	}

	file, err := os.Open(input_json)
	if err != nil {
		log.Fatalf("Error to open file: %v", err)
	}
	defer file.Close() // Certifique-se de fechar o arquivo após a leitura

	// Mapa para armazenar o objeto JSON
	var data map[string]map[string]interface{}

	// Decodifica o JSON do arquivo
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatalf("Erro ao decodificar JSON: %v", err)
	}

	// mmdb tree
	writer, _ := mmdbwriter.New(mmdbwriter.Options{
		DatabaseType: mmdb_type,
	})

	// Itera sobre cada chave no objeto
	for key, subObject := range data {
		fmt.Printf("Inserting: %s\n", key)
		addRecord(key, &subObject, writer)
		delete(data, key)
	}

	debug.FreeOSMemory()
	fmt.Println("Saving...")
	writeMMDB(output_mmdb, writer)
}
