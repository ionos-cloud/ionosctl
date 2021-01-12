package commands

import "github.com/spf13/cobra"

func list() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List command for resources",
	}

	listCmd.AddCommand(
		listDataCenter(),
	)

	return listCmd
}

func create() *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create command for resources",
	}

	createCmd.AddCommand(
		createDataCenter(),
	)

	return createCmd
}

func update() *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update command for resources",
	}

	updateCmd.AddCommand(
		updateDataCenter(),
	)

	return updateCmd
}

func delete() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete command for resources",
	}

	deleteCmd.AddCommand(
		deleteDataCenter(),
	)

	return deleteCmd
}
