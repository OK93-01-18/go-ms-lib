package views

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type SimpleViewEngine struct {
	*sync.RWMutex
	templatePaths      map[string]string
	templates          map[string]string
	reloadBeforeRender bool
}

func (d *SimpleViewEngine) init() error {

	d.templates = make(map[string]string)

	if d.reloadBeforeRender {
		return nil
	}

	buffer := make(map[string]string)
	for templateName, templatePath := range d.templatePaths {
		file, err := os.OpenFile(filepath.Clean(templatePath), os.O_RDONLY, os.ModePerm)
		if err != nil {
			return err
		}

		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}

		buffer[templateName] = string(bytes)
	}

	d.Lock()
	d.templates = buffer
	d.Unlock()

	return nil
}

func (d *SimpleViewEngine) loadTemplate(path string) (string, error) {
	var res string

	file, err := os.OpenFile(filepath.Clean(path), os.O_RDONLY, os.ModePerm)
	if err != nil {
		return res, err
	}

	bytes, err := ioutil.ReadAll(file)
	return string(bytes), err
}

func (d *SimpleViewEngine) reload(name string) error {
	d.RLock()
	templatePath, ok := d.templatePaths[name]
	d.RUnlock()
	if !ok {
		return errors.New("template " + name + " not found")
	}

	template, err := d.loadTemplate(templatePath)
	if err != nil {
		return err
	}

	d.Lock()
	d.templates[name] = template
	d.Unlock()

	return nil
}

func (d *SimpleViewEngine) Render(name string, params map[string]string) (string, error) {

	var res string

	if d.reloadBeforeRender {
		err := d.reload(name)
		if err != nil {
			return res, err
		}
	}

	d.RLock()
	template, ok := d.templates[name]
	d.RUnlock()

	if !ok {
		return res, errors.New("template " + name + " not found")
	}

	for key, value := range params {
		template = strings.ReplaceAll(template, key, value)
	}

	return template, nil
}

type Config struct {
	Templates          map[string]string
	ReloadBeforeRender bool
}

func NewSimpleViewEngine(conf *Config) (Engine, error) {
	e := &SimpleViewEngine{
		RWMutex:            &sync.RWMutex{},
		templatePaths:      conf.Templates,
		reloadBeforeRender: conf.ReloadBeforeRender,
	}

	return e, e.init()
}
