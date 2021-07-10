package banner

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/utils/colors"
)

type bannerImpl struct {
	Banner
}

func New() Banner {
	return &bannerImpl{}
}

func (b *bannerImpl) PrintAnchor() {
	fmt.Printf(colors.Blue + `
     \                  |                  
    _ \    __ \    __|  __ \    _ \    __| 
   ___ \   |   |  (     | | |  (   |  |    
 _/    _\ _|  _| \___| _| |_| \___/  _|

		` + colors.Reset)
}
