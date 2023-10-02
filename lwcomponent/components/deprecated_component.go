package components

// type DeprecatedComponent struct {
// 	Name          string
// 	Size          int64
// 	ComponentType Type
// }

// // Execute implements CDKComponent
// func (c *DeprecatedComponent) Execute(args []string, envs ...string) (stdout string, stderr string, err error) {
// 	if c.ComponentType == BinaryType || c.ComponentType == CommandType {
// 		err = errors.Errorf("component '%s' is not a binary", c.Name)
// 		return
// 	}

// 	// TODO verify component

// 	path, err := executablePath(c.Name)
// 	if err != nil {
// 		return
// 	}

// 	return execute(path, args, envs...)
// }

// // InstallMessage implements CDKComponent
// func (c *DeprecatedComponent) InstallMessage() string {
// 	return ""
// }

// // PrintSummary implements CDKComponent
// func (c *DeprecatedComponent) PrintSummary() []string {
// 	colorize := color.New(color.FgRed, color.Bold)

// 	version, err := c.Version()
// 	if err != nil {
// 		panic(err)
// 	}

// 	return []string{
// 		colorize.Sprintf("Deprecated"),
// 		c.Name,
// 		version.String(),
// 		"",
// 	}
// }

// // Version implements CDKComponent
// func (c *DeprecatedComponent) Version() (*semver.Version, error) {
// 	return localVersion(c.Name)
// }

// func NewDeprecatedComponent(name string) (*DeprecatedComponent, error) {
// 	dir, err := componentDirectory(name)
// 	if err != nil {
// 		return nil, errors.Wrap(err, fmt.Sprintf("deprecated component '%s'", name))
// 	}

// 	stats, err := os.Stat(filepath.Join(dir, name))
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Component binary in KB
// 	size := stats.Size() / 1024

// 	// TODO how do we get Type + description?
// 	return &DeprecatedComponent{Name: name, Size: size}, nil
// }
