
package main

import (
	"flag"
)

//Mudar isto para não usar o visualization, visto que isto poderia ser usado sem visualização
//Mas ao mesmo tempo, num contexto distribuido podemos não saber todos os Nodes.
func main(){
	vis_address := flag.String("address", "", "V")
	flag.Parse()

	if *vis_address == "" {
		flag.PrintDefaults()
	}

}

