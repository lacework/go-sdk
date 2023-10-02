package components

// type DevComponent struct {
// 	Name          string
// 	Size          int64
// 	ComponentType Type
// }

// // Execute implements CDKComponent
// func (c *DevComponent) Execute(args []string, envs ...string) (stdout string, stderr string, err error) {
// 	if c.ComponentType == BinaryType || c.ComponentType == CommandType {
// 		err = errors.Errorf("component '%s' is not a binary", c.Name)
// 		return
// 	}

// 	// verify component

// 	path, err := executablePath(c.Name)
// 	if err != nil {
// 		return
// 	}

// 	return execute(path, args, envs...)
// }

// // InstallMessage implements CDKComponent
// func (c *DevComponent) InstallMessage() string {
// 	return ""
// }

// // PrintSummary implements CDKComponent
// func (c *DevComponent) PrintSummary() []string {
// 	colorize := color.New(color.FgBlue, color.Bold)

// 	version, err := c.Version()
// 	if err != nil {
// 		panic(err)
// 	}

// 	return []string{
// 		colorize.Sprintf("Development"),
// 		c.Name,
// 		version.String(),
// 		"",
// 	}
// }

// // Version implements CDKComponent
// func (c *DevComponent) Version() (*semver.Version, error) {
// 	return localVersion(c.Name)
// }

// func NewDevComponent(name string) (*DevComponent, error) {
// 	dir, err := componentDirectory(name)
// 	if err != nil {
// 		return nil, errors.Wrap(err, fmt.Sprintf("development component '%s'", name))
// 	}

// 	stats, err := os.Stat(filepath.Join(dir, name))
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Component binary in KB
// 	size := stats.Size() / 1024

// 	return &DevComponent{Name: name, Size: size}, nil
// }
