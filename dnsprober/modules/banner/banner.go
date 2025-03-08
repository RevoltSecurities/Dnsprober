package banner

import (
	"fmt"
	"math/rand"

	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/logger"
	"github.com/common-nighthawk/go-figure"
	"github.com/logrusorgru/aurora"
)

var loggers = logger.New(true)

func Randomchoice(choices []string) string {
	choosed := rand.Intn(len(choices))
	return choices[choosed]
}

func BannerGenerator(banner_name string) aurora.Value {
	choices := []string{"big", "ogre", "shadow", "script", "graffiti", "slant"}
	colors := []string{"green", "cyan", "blue", "white", "magenta"}
	banner := figure.NewFigure(banner_name, Randomchoice(choices), true)
	banners := loggers.Colorizer(fmt.Sprintf(`%s`, banner.String()), Randomchoice(colors))
	return banners
}
