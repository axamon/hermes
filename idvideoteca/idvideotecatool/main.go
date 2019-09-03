package main

import (
	"strings"
	"fmt"
	"github.com/axamon/hermes/idvideoteca"
)

func main() {

	var str = `http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2019/01/50734100/SS/11483192/11483192_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2014/09/50434361/SS/20086428/20086428_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2014/09/50434361/SS/20086428/20086428_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2014/09/50434361/SS/20086428/20086428_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000243/SS/20089777/20089777_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000241/SS/20089779/20089779_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2019/01/50734100/SS/11483192/11483192_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2019/01/50734100/SS/11483192/11483192_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000242/SS/20089785/20089785_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000242/SS/20089785/20089785_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000242/SS/20089785/20089785_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000241/SS/20089779/20089779_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2014/09/50434361/SS/20086428/20086428_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000243/SS/20089777/20089777_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2019/01/50734100/SS/11483192/11483192_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000242/SS/20089785/20089785_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000242/SS/20089785/20089785_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000241/SS/20089779/20089779_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2019/01/50734100/SS/11483192/11483192_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2019/01/50734100/SS/11483192/11483192_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2019/01/50734100/SS/11483192/11483192_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2019/01/50734100/SS/11483192/11483192_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2019/01/50734100/SS/11483192/11483192_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000241/SS/20089779/20089779_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000243/SS/20089777/20089777_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000241/SS/20089779/20089779_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2019/01/50734100/SS/11483192/11483192_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000243/SS/20089777/20089777_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2019/01/50734100/SS/11483192/11483192_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2019/01/50734100/SS/11483192/11483192_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2019/01/50734100/SS/11483192/11483192_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2019/01/50734100/SS/11483192/11483192_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2014/09/50434361/SS/20086428/20086428_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2019/01/50734100/SS/11483192/11483192_HD.ism/Manifest
	http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2019/01/50734100/SS/11483192/11483192_HD.ism/Manifest`
	

	elements := strings.Split(str,"\n")


	for _, element := range elements  {
		//fmt.Println(element)
		idv, err := idvideoteca.Find(element)
		if err != nil {
			idv = "NON DISPOBINILE"
		}
		fmt.Println(idv)
	}
}