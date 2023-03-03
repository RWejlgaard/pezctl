package cmd

import (
    "fmt"
    "pezctl/util"
    "regexp"
    "strings"
    "sync"

    "github.com/spf13/cobra"
)

var Sequential bool
var Filter string

func init() {
    rootCmd.AddCommand(bulkCmd)
    bulkCmd.Flags().BoolVarP(&Sequential, "sequential", "s", false, "Runs commands one after another instead of all at once")
    bulkCmd.Flags().StringVarP(&Filter, "filter", "f", "pez-.*", "Applies regex filter to kube context selection")
}

var bulkCmd = &cobra.Command{
    Use:   "bulk",
    Short: "Runs commands on all clusters at once",
    Long: `bulk runs any kind of kubectl command on all or some contexts available in the default kubeconfig.
The targeting can be changed by using the -f or --filter flag to select with regex`,

    Run: func(cmd *cobra.Command, args []string) {
        contexts := filterArray(getKubernetesContexts(), Filter)
        if Sequential {
            for _, context := range contexts {
                output, err := util.RunCommand(fmt.Sprintf("kubectl --context=%s %s", context, strings.Join(args, " ")), "")
                util.HandleErr(err, true, Logger)
                fmt.Println(output)
            }
        } else {
            var wg sync.WaitGroup
            for _, context := range contexts {
                wg.Add(1)

                context := context

                go func() {
                    defer wg.Done()
                    output, err := util.RunCommand(fmt.Sprintf("kubectl --context=%s %s", context, strings.Join(args, " ")), "")
                    util.HandleErr(err, true, Logger)
                    fmt.Println(">>>", context)
                    fmt.Println(output)
                }()

                wg.Wait()
            }
        }
    },
}

func getKubernetesContexts() []string {
    command := "kubectl config get-contexts | tail -n +2 | awk '{ print $2 }'"
    if Verbose {
        Logger.Infow("running command", "command", command)
    }
    output, err := util.RunCommand(command, "")
    if err != nil {
        Logger.Fatalw("failed to list contexts", "error", err.Error())
    }
    contexts := strings.Split(output, "\n")
    return contexts[:len(contexts)-1]

}

func filterArray(arr []string, pattern string) []string {
    filtered := []string{}
    regex, err := regexp.Compile(pattern)
    util.HandleErr(err, true, Logger)

    for _, str := range arr {
        if regex.MatchString(str) {
            filtered = append(filtered, str)
        }
    }

    return filtered
}
