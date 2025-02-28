package root

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/cmd/config"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/cmd/config/open"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/cmd/create"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/cmd/export"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/cmd/list"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/cmd/palette"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/cmd/scripts"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/cmd/show"
	"github.com/spinozanilast/aseprite-assets-cli/pkg/environment"
)

func NewRootCmd(env *environment.Environment) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "aseprite-assets",
		Short: "A CLI tool to manage aseprite assets",
		Long: `A CLI tool to manage aseprite assets. 
			This tool allows you to manage aseprite assets by listing, adding, removing and renaming them.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(`
--______-------------------------------------------__----__-----------------------______-----------------------------------__---------------
-/------\-----------------------------------------/--|--/--|---------------------/------\---------------------------------/--|--------------
/$$$$$$--|--_______---______----______----______--$$/--_$$-|_-----______--------/$$$$$$--|--_______---_______---______---_$$-|_-----_______-
$$-|__$$-|-/-------|-/------\--/------\--/------\-/--|/-$$---|---/------\-------$$-|__$$-|-/-------|-/-------|-/------\-/-$$---|---/-------|
$$----$$-|/$$$$$$$/-/$$$$$$--|/$$$$$$--|/$$$$$$--|$$-|$$$$$$/---/$$$$$$--|------$$----$$-|/$$$$$$$/-/$$$$$$$/-/$$$$$$--|$$$$$$/---/$$$$$$$/-
$$$$$$$$-|$$------\-$$----$$-|$$-|--$$-|$$-|--$$/-$$-|--$$-|-__-$$----$$-|------$$$$$$$$-|$$------\-$$------\-$$----$$-|--$$-|-__-$$------\-
$$-|--$$-|-$$$$$$--|$$$$$$$$/-$$-|__$$-|$$-|------$$-|--$$-|/--|$$$$$$$$/-------$$-|--$$-|-$$$$$$--|-$$$$$$--|$$$$$$$$/---$$-|/--|-$$$$$$--|
$$-|--$$-|/-----$$/-$$-------|$$----$$/-$$-|------$$-|--$$--$$/-$$-------|------$$-|--$$-|/-----$$/-/-----$$/-$$-------|--$$--$$/-/-----$$/-
$$/---$$/-$$$$$$$/---$$$$$$$/-$$$$$$$/--$$/-------$$/----$$$$/---$$$$$$$/-------$$/---$$/-$$$$$$$/--$$$$$$$/---$$$$$$$/----$$$$/--$$$$$$$/--
------------------------------$$-|----------------------------------------------------------------------------------------------------------
------------------------------$$-|----------------------------------------------------------------------------------------------------------
------------------------------$$/-----------------------------------------------------------------------------------------------------------

CLI interface for aseprite assets interaction. For in-terminal opening of aseprite files and more.
`)
		},
	}

	cmd.AddCommand(
		config.NewConfigCmd(env),
		create.NewCraeteCmd(env),
		export.NewExportCmd(env),
		list.NewListCmd(env),
		open.NewConfigOpenCmd(env),
		palette.NewPaletteCmd(env),
		scripts.NewScriptsCmd(env),
		show.NewShowCmd(env),
	)

	return cmd
}
