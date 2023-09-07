// Credits to @omissis
// https://github.com/omissis/kube-apiserver-proxy

package cobrax

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitEnvs(envPrefix string) *viper.Viper {
	v := viper.New()

	v.SetEnvPrefix(envPrefix)

	v.AutomaticEnv()

	return v
}

func BindFlags(cmd *cobra.Command, v *viper.Viper, envPrefix string) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if strings.Contains(f.Name, "-") {
			envSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))

			env := envSuffix
			if envPrefix != "" {
				env = fmt.Sprintf("%s_%s", envPrefix, envSuffix)
			}

			if err := v.BindEnv(f.Name, env); err != nil {
				panic(err)
			}
		}

		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)

			if err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val)); err != nil {
				panic(err)
			}
		}
	})
}

func Flag[T bool | int | string | []string](cmd *cobra.Command, name string) any {
	var f T

	if cmd == nil {
		return f
	}

	if cmd.Flag(name) == nil {
		return f
	}

	v := cmd.Flag(name).Value.String()

	if v == "true" {
		return true
	}

	if v == "false" {
		return false
	}

	if vv, err := strconv.Atoi(v); err == nil {
		return vv
	}

	return v
}

func FlagSlice(cmd *cobra.Command, name string) []string {
	if cmd == nil {
		return nil
	}

	if cmd.Flag(name) == nil {
		return nil
	}

	s := cmd.Flag(name).Value.String()

	if !strings.Contains(s, ",") {
		return []string{s}
	}

	res := strings.Split(s, ",")

	return res
}
