package ui

import (
	"os"
	"path"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/pug"
)

func GetMain(ctx *fiber.Ctx) error {
	return ctx.Render("index", nil)
}

func StartUp(bind string) error {

	// As the pug engine resolves paths relative to the workdir, but paths inside the pug files are relative to the
	// files themselves, I'mma change into that directory. This means, consecutively all paths should be handled
	// relative to the original workdir, which we're gonna store as origWd.
	origWd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir("./app"); err != nil {
		return err
	}

	// Get pug templating engine and assign it as the fiber view engine
	engine := pug.New(".", ".pug")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Here we compile all the routing for the app
	AttachStaticRoutes(app, origWd)
	AttachRoutes(app)

	// Let's go then... Start the App listener.
	if err := app.Listen(bind); err != nil {
		return err
	}

	return nil
}

// AttachRoutes defines all UI paths for this App
func AttachRoutes(app *fiber.App) {
	app.Get("/", GetMain)
}

// AttachStaticRoutes defines all resource or static paths for this App
func AttachStaticRoutes(app *fiber.App, origWd string) {
	app.Static("/css", path.Join(origWd, "app/css"))
	app.Static("/img", path.Join(origWd, "app/img"))
	app.Static("/js", path.Join(origWd, "app/js"))
}
