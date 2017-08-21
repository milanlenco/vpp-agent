# Flavors

A flavor is a reusable collection of plugins. It allows inspecting
capabilities available to an agent. By convention, a Flavor should
provide:
 * a method called `Inject()`, which should wire the flavor's plugins
   together and to their external dependencies (typcially via dependency
   injection. 
 * a method called `Plugins()`, which should return all plugins contained
   in the Flavor, including plugins contained in [reused Flavors][1].

Example:
```go
type MyFlavor struct {
	injected     bool
	Aplugin      a.Plugin
	Bplugin      b.Plugin
	Cplugin      c.Plugin
}

func (f *MyFlavor) Inject() error {
	if f.injected {
		return nil
	}
	f.injected = true
	f.Bplugin.A = &f.Aplugin
	f.Cplugin.B = &f.Bplugin
	return nil
}

func (f *MyFlavor) Plugins() []*core.NamedPlugin {
	f.Inject()
	return core.ListPluginsInFlavor(f)
}
```

[1]: https://github.com/ligato/cn-infra/blob/master/docs/guidelines/PLUGIN_FLAVORS.md