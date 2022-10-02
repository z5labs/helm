package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/z5labs/helm/tools/ksplit/kubeyaml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logLevel zapcore.Level

func (l logLevel) String() string {
	return (zapcore.Level)(l).String()
}

func (l *logLevel) Set(s string) error {
	return (*zapcore.Level)(l).Set(s)
}

func (l logLevel) Type() string {
	return "Level"
}

var rootCmd = &cobra.Command{
	Use:   "ksplit FILE|-",
	Short: "Split joing K8 yamls into multiple files",
	Args:  cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var lvl zapcore.Level
		lvlStr := cmd.Flags().Lookup("log-level").Value.String()
		err := lvl.UnmarshalText([]byte(lvlStr))
		if err != nil {
			panic(err)
		}

		cfg := zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		cfg.OutputPaths = []string{viper.GetString("log-file")}
		l, err := cfg.Build(zap.IncreaseLevel(lvl))
		if err != nil {
			panic(err)
		}

		zap.ReplaceGlobals(l)
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		outDir := viper.GetString("output-dir")
		outDir, err = filepath.Abs(outDir)
		if err != nil {
			zap.L().Fatal("unable to get absolute path for output directory", zap.Error(err))
		}

		b, err := readFromFileOrStdin(cmd, args[0])
		if err != nil {
			zap.L().Fatal("failed to read source file", zap.Error(err))
		}

		err = kubeyaml.Split(
			bytes.NewReader(b),
			afero.NewBasePathFs(afero.NewOsFs(), outDir),
		)
		if err != nil {
			zap.L().Fatal("failed to split compound kube yaml file", zap.Error(err))
		}
	},
}

func init() {
	// Persistent flags
	lvl := logLevel(zapcore.ErrorLevel)
	rootCmd.PersistentFlags().Var(&lvl, "log-level", "Specify log level")
	rootCmd.PersistentFlags().String("log-file", "stderr", "Specify log file")

	rootCmd.Flags().StringP("output-dir", "o", ".", "Specify output directory")

	viper.BindPFlag("log-file", rootCmd.PersistentFlags().Lookup("log-file"))
	viper.BindPFlag("output-dir", rootCmd.Flags().Lookup("output-dir"))
}

func readFromFileOrStdin(cmd *cobra.Command, filename string) ([]byte, error) {
	if filename == "-" {
		return io.ReadAll(cmd.InOrStdin())
	}

	absFileName, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(absFileName)
}
