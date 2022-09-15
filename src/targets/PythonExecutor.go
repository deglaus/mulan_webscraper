package targets

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func getPWD() string {
	// Get current working directory
	pwd, err1 := os.Getwd()
	if err1 != nil {
		fmt.Println(err1)
	}
	// "Move" into the targets folder
	pwd = pwd + "/src/targets"
	return pwd
}

func ExecuteBlocketandTradera(searchString string, wg *sync.WaitGroup) {
	defer wg.Done()

	pwd := getPWD()
	// Execute the scraper function in Blocket.py
	cmdExec := exec.Command("python3", "-c", "import BlocketandTradera; BlocketandTradera.scrapeAndInsert('"+searchString+"')")
	// Inside the targets folder
	cmdExec.Dir = pwd

	out, err := cmdExec.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))

}
