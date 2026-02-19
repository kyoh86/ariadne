package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage ariadne configuration",
	}

	cmd.AddCommand(newConfigShowCommand())
	cmd.AddCommand(newConfigSetGameDirRootCommand())
	cmd.AddCommand(newConfigSetLauncherDirCommand())

	return cmd
}

func newConfigShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			fmt.Printf("game-dir-root: %s\n", cfg.GameDirRoot)
			fmt.Printf("launcher-dir: %s\n", cfg.LauncherDir)
			fmt.Printf("active-profile: %s\n", cfg.ActiveProfile)
			return nil
		},
	}
}

func newConfigSetGameDirRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set-game-dir-root <path>",
		Short: "Set default root directory for game-dir",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("resolve path: %w", err)
			}
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			cfg.GameDirRoot = root
			if err := saveConfig(cfg); err != nil {
				return err
			}
			fmt.Printf("set game-dir-root: %s\n", root)
			return nil
		},
	}
}

func newConfigSetLauncherDirCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set-launcher-dir <path>",
		Short: "Set Minecraft launcher directory path",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("resolve path: %w", err)
			}
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			cfg.LauncherDir = root
			if err := saveConfig(cfg); err != nil {
				return err
			}
			fmt.Printf("set launcher-dir: %s\n", root)
			return nil
		},
	}
}
