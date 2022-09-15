package main

import "fmt"

/* Interfaces */

type IFAppFactory interface {
	GetImage() IFImage
	GetBin() IFBinary
}

type IFImage interface {
	InjectBin(IFBinary)
	Run()
}

type IFBinary interface {
	Exec()
}

/* Image Struct */

type OSImage struct {
	EntryPoint IFBinary
}

func (o *OSImage) InjectBin(bin IFBinary) {
	o.EntryPoint = bin
}

func (o *OSImage) Run() {
	if o.EntryPoint == nil {
		fmt.Println("[WARN] Entrypoint not set yet !")
		return
	}

	o.EntryPoint.Exec()
}

type MacOSImage struct {
	OSImage
}

type LinuxImage struct {
	OSImage
}

/* Bin struct */
type AppBin struct {
	Content string
}

func (a *AppBin) Exec() {
	fmt.Println(a.Content)
}

type BinForMac struct {
	AppBin
}

type BinForLinux struct {
	AppBin
}

/* Factory struct */
type MacAppFactory struct {}

type LinuxAppFactory struct {}

func (m *MacAppFactory) GetImage() IFImage {
	return &MacOSImage{}
}

func (m *MacAppFactory) GetBin() IFBinary {
	return &BinForMac{
		AppBin {
			Content: "Running binary built for Mac OS ...",
		},
	}
}

func (l *LinuxAppFactory) GetImage() IFImage {
	return &LinuxImage{}
}

func (l *LinuxAppFactory) GetBin() IFBinary {
	return &BinForLinux{
		AppBin {
			Content: "Running binary built for Linux ...",
		},
	}
}

/* Method exposed to user */

func GetAppFactory(os string) (IFAppFactory, error) {
	if os == "darwin" {
		return &MacAppFactory{}, nil
	}

	if os == "linux" {
		return &LinuxAppFactory{}, nil
	}

	return nil, fmt.Errorf("Unsupport OS type: %s", os)
}

func main() {
	wanted := []string{"darwin", "linux", "windows"}
	factories := []IFAppFactory{}

	/* Init factories */
	for _, os := range wanted {
		if fac, err := GetAppFactory(os); err == nil {
			factories = append(factories, fac)
		}else{
			fmt.Println(err)
		}
	}

	for _, fac := range factories {
		img := fac.GetImage()
		img.InjectBin(fac.GetBin())
		img.Run()
	}
}
