package banner

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/utils/colors"
)

func Print() {
	fmt.Printf(colors.Blue + `
     \                  |                  
    _ \    __ \    __|  __ \    _ \    __| 
   ___ \   |   |  (     | | |  (   |  |    
 _/    _\ _|  _| \___| _| |_| \___/  _|

		`)
}
