package shell

import "github.com/gofiber/fiber/v2"

func execScriptHandler(c *fiber.Ctx) error {
	return ExecShell(c, app.Config.Server.ScriptPath)
}
