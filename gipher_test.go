package gipher

import (
	"fmt"
	"testing"
)

func TestCreateTimeStampGIF(t *testing.T) {
	const fontFilename = "./test_data/Raleway-Black.ttf"

	backgroundGIF := []string{"1e.gif", "2e.gif", "5e.gif", "10e.gif"}

	for _, file := range backgroundGIF {
		inputFilename := fmt.Sprintf("./test_data/%s", file)
		outputFilename := fmt.Sprintf("./test_data/output_%s", file)
		err := CreateTimeStampGIF(inputFilename, outputFilename, fontFilename)
		if err != nil {
			t.Error(err)
		}
	}

}
