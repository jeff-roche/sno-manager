package cmd

import (
	"fmt"
	"os"
	"snoman/internal/vms/network"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the specified resource config",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Fatalf("Error executing the generate command: %v", ErrResourcetTypeNotSpecified)
	},
}

func initGenerateCmd() {
	rootCmd.AddCommand(generateCmd)

	// Subcommands

	// Generate network configuration
	generateCmd.AddCommand(generateNetCmd)
	generateNetCmd.Flags().String("from-xml", "", "Path to the XML file to use for spec generation")
	generateNetCmd.Flags().String("from", "", "Path to the spec file to use for spec generation. This will just sanitize the input.")
	generateNetCmd.Flags().Bool("xml", false, "Generate the libvirt xml config for the network")
}

// Destroy VM Network
var generateNetCmd = &cobra.Command{
	Use:   "network",
	Short: "Generate a libvirt network spec",
	Long: `
	Generate a libvirt network spec

	If --from-xml or --from-yaml are not specified, the default configuration will be used
	`,
	Run: func(cmd *cobra.Command, args []string) {
		outputXML, _ := cmd.Flags().GetBool("xml")
		source, _ := cmd.Flags().GetString("from")
		xmlSource, _ := cmd.Flags().GetString("from-xml")

		var spec *network.VirtualMachineNetworkSpec
		if source != "" {
			data, err := os.ReadFile(source)
			if err != nil {
				logger.Fatalf("unable to read the spec file: %w", err)
			}

			spec = &network.VirtualMachineNetworkSpec{}
			spec.UnmarshalYAML(data)
		} else if xmlSource != "" {
			data, err := os.ReadFile(xmlSource)
			if err != nil {
				logger.Fatalf("unable to read the xml file: %w", err)
			}

			spec = &network.VirtualMachineNetworkSpec{}
			spec.UnmarshalXML(data)
		} else {
			spec = network.GetDefaultVirtualMachineNetworkSpec()
		}

		var output string
		var err error
		if outputXML {
			output, err = spec.MarshalXML()
		} else {
			output, err = spec.MarshalYAML()
		}

		if err != nil {
			logger.Fatalf("unable to generate spec: %w", err)
		}

		fmt.Println(output)
	},
}