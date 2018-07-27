package plugin

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plugin"

	"event_processor"
	"event_processor/collectors/aws"
)

type Plugins []event_processor.Plugin

func (plgs Plugins) GetPluginList(pathS, ext string) []string {
	var files []string
	filepath.Walk(pathS, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			if filepath.Ext(path) == ext {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	return files
}

func (plugins Plugins) LoadPlugins() {

	plugins = make([]event_processor.Plugin, 0)
	pluginNames := plugins.GetPluginList("../plugins/", ".so")

	for _, pluginName := range pluginNames {
		plug, err := plugin.Open(fmt.Sprintf("%s%s", "../plugins/", pluginName))
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		processor, err := plug.Lookup("MessageProc")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		messageProc, ok := processor.(event_processor.MessageProc)
		if !ok {
			log.Println("Unexpected type from module symbols")
			os.Exit(1)
		}

		vals := messageProc.Register()
		log.Println(vals)
		plugins = append(plugins, event_processor.Plugin{
			MessageProc: messageProc,
			Type:        vals["Type"],
			Provider:    vals["Provider"],
			QueueName:   vals["QueueName"],
		})
	}
}

func (plugins Plugins) SetupPlugins() {
	for _, plugin := range plugins {
		switch plugin.Type {
		case "QueueTask":
			switch plugin.Provider {
			case "AWS":
				collector, err := aws.CreateCollector(plugin)
				if err != nil {
					panic(err)
				}
				plugin.Collector = collector
				break
			default:
				log.Printf("Plugin Queue Provider %s not implemented yet!\n", plugin.Provider)
				break
			}
			break
		case "TimerTask":
			switch plugin.Provider {
			case "AWS":
				//TODO Setup Timer for Task - by using internal queue
				break
			default:
				log.Printf("Plugin Queue Provider %s not implemented yet\n", plugin.Provider)
				break
			}
			break
		default:
			log.Printf("Plugin type %s not implemented yet!!!\n", plugin.Type)
			break
		}
	}
}
