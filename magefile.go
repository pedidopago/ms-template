// +build mage

package main

// https://magefile.org

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/magefile/mage/sh"
	"github.com/mozillazg/go-slugify"
	"go.pedidopago.com.br/microservices/defaults"
	"go.pedidopago.com.br/microservices/xcmd"
)

const (
	dbfiles = "file://%GIT_ROOT%/internal/database/migrations"
	dbcs    = "mysql://testuser:123456789@tcp(localhost)/ms_xyz?parseTime=true"
)

func Migrate(version string) error {
	loadEnvs()
	m, err := migrate.New(xs(defaults.String(os.Getenv("DB_FILES"), dbfiles)), xs(defaults.String(os.Getenv("DB_CS"), dbcs)))
	if err != nil {
		return err
	}
	if version == "" || version == "up" {
		return m.Up()

	}
	return nil
}

func Newmigration(name string) error {
	loadEnvs()
	//
	if name == "" {
		return errors.New("migration name is required")
	}
	//
	u0, err := url.Parse(xs(defaults.String(os.Getenv("DB_FILES"), dbfiles)))
	if err != nil {
		return err
	}
	t0 := time.Now()
	mformat := "%d_%s.%s.sql"
	//
	migname := slugify.Slugify(name)
	upname := fmt.Sprintf(mformat, t0.UnixNano(), migname, "up")
	downname := fmt.Sprintf(mformat, t0.UnixNano(), migname, "down")
	//
	f, err := os.Create(path.Join(u0.Path, upname))
	if err != nil {
		return fmt.Errorf("unable to create %s %w", upname, err)
	}
	f.Close()
	//
	f, err = os.Create(path.Join(u0.Path, downname))
	if err != nil {
		return fmt.Errorf("unable to create %s %w", downname, err)
	}
	f.Close()
	//
	println("created " + upname)
	println("created " + downname)
	if v := os.Getenv("OPEN_CMD"); v != "" {
		xcmd.Run(v, path.Join(u0.Path, upname))
		xcmd.Run(v, path.Join(u0.Path, downname))
	}
	return nil
}

func Setup() error {
	// fetch dependencies!
	if err := setupDependencies(); err != nil {
		return err
	}
	// # remove >>
	if err := setupInstall(); err != nil {
		return err
	}
	// # remove <<
	if _, e := os.Stat("dbm.env"); e != nil {
		nd, _ := ioutil.ReadFile(".name")
		name := strings.TrimSpace(string(nd))
		print("mariadb user (testuser): ")

		rdr := bufio.NewReader(os.Stdin)
		str, err := rdr.ReadString('\n')
		if err != nil {
			return err
		}
		str = strings.TrimSpace(str)
		if str == "" {
			str = "testuser"
		}
		dbuser := "testuser"

		println("")
		print("mariadb password (123456789): ")

		str, err = rdr.ReadString('\n')
		if err != nil {
			return err
		}
		str = strings.TrimSpace(str)
		if str == "" {
			str = "123456789"
		}
		dbpass := str

		println("please run this:")
		println("CREATE DATABASE IF NOT EXISTS `ms_" + name + "`;")

		ndbcs := dbuser + ":" + dbpass + "@tcp(localhost)/ms_" + name + "?parseTime=true"

		if err := replaceStringInFile("cmd/"+name+"/main.go", "testuser:123456789@tcp(localhost)/ms_xyz?parseTime=true", ndbcs); err != nil {
			return err
		}

		if err := ioutil.WriteFile("dbm.env", []byte("DB_FILES="+dbfiles+"\nDB_CS="+ndbcs+"\nOPEN_CMD=code\nDOCKER_REGISTRY=registry.docker.pedidopago.com.br/ms/"+name), 0644); err != nil {
			return err
		}
	}
	return nil
}

func Devbuild() error {
	//
	nd, _ := ioutil.ReadFile(".name")
	name := strings.TrimSpace(string(nd))
	//
	_ = sh.Run("mkdir", "-p", "tmp")
	_, err := sh.Exec(map[string]string{
		"GOOS":   "linux",
		"GOARCH": "amd64",
	}, os.Stdout, os.Stderr, "go", "build", "-o", "tmp/service_linux_x64", "cmd/"+name+"/main.go")
	return err
}

func Devdocker() error {
	loadEnvs()
	//
	nd, _ := ioutil.ReadFile(".name")
	name := strings.TrimSpace(string(nd))
	//
	version := defaults.String(os.Getenv("VERSION"), "dev")
	tags := []string{version}
	registry := defaults.String(os.Getenv("DOCKER_REGISTRY"), "registry.docker.pedidopago.com.br/ms/"+name)

	//
	if v := os.Getenv("DOCKER_TAGS"); v != "" {
		ts := strings.Split(v, ",")
		tags = append(tags, ts...)
	}
	if err := Devbuild(); err != nil {
		return err
	}

	// docker build --build-arg VERSION=${VERSION} -t ${REGISTRY}:${VERSION} .
	args := []string{
		"build",
		"--build-arg",
		"VERSION=" + version,
	}
	for _, tag := range tags {
		args = append(args, "-t", registry+":"+tag)
	}
	args = append(args, "-f", "dev.Dockerfile", ".")
	return sh.Run("docker", args...)
}

func Devdeploy() error {
	if err := Devdocker(); err != nil {
		return err
	}
	//
	nd, _ := ioutil.ReadFile(".name")
	name := strings.TrimSpace(string(nd))
	//
	registry := defaults.String(os.Getenv("DOCKER_REGISTRY"), "registry.docker.pedidopago.com.br/ms/"+name)
	return sh.Run("docker", "push", registry)
}

func Gen() error {
	return sh.Run("go", "generate", "./...")
}

func Migrationtest(dbcs, mpath string) error {
	// loadEnvs()
	m, err := migrate.New(mpath, dbcs)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		return err
	}
	if err := m.Down(); err != nil {
		return err
	}
	return nil
}

func setupDependencies() error {
	if err := sh.Run("mockery", "--help"); err != nil {
		// install mockery
		if runtime.GOOS == "darwin" {
			if err := sh.Run("brew", "help"); err == nil {
				// install mockery via homebrew
				if err := sh.RunV("brew", "install", "mockery"); err != nil {
					return err
				}
			} else {
				// install vi go get
				if _, err := sh.Exec(map[string]string{
					"GO111MODULE": "off",
				}, os.Stdout, os.Stderr, "go", "get", "github.com/vektra/mockery/v2/.../"); err != nil {
					return err
				}
			}
		} else {
			if _, err := sh.Exec(map[string]string{
				"GO111MODULE": "off",
			}, os.Stdout, os.Stderr, "go", "get", "github.com/vektra/mockery/v2/.../"); err != nil {
				return err
			}
		}
	}
	return nil
}

// # remove >>
func setupInstall() error {
	if _, e := os.Stat(".name"); e != nil {
		print("service name: ")
		rdr := bufio.NewReader(os.Stdin)
		str, err := rdr.ReadString('\n')
		if err != nil {
			return err
		}
		str = strings.TrimSpace(str)
		name := slugify.Slugify(str)
		if name == "" {
			return errors.New("invalid service name")
		}
		print("\nmodule name+path (e.g. github.com/pedidopago/ms-template): ")
		str, err = rdr.ReadString('\n')
		if err != nil {
			return err
		}
		str = strings.TrimSpace(str)
		module := str
		if module == "" {
			return errors.New("invalid module name")
		}
		// ungit it
		if err := sh.Run("rm", "-rf", ".git"); err != nil {
			return err
		}
		// git it again
		if err := sh.Run("git", "init"); err != nil {
			return err
		}
		// rename folders
		if err := sh.Run("mv", "cmd/xyzservice", "cmd/"+name); err != nil {
			return err
		}
		if err := sh.Run("mv", "internal/xyzservice/xyzservice.go", "internal/xyzservice/"+name+".go"); err != nil {
			return err
		}
		if err := sh.Run("mv", "internal/xyzservice", "internal/"+name); err != nil {
			return err
		}
		if err := sh.Run("mv", "protos/xyzpb", "protos/"+name+"pb"); err != nil {
			return err
		}
		// replace strings
		if err := replaceStringInFile("cmd/"+name+"/main.go", "xyzservice", name); err != nil {
			return err
		}
		if err := replaceStringInFile("internal/"+name+"/"+name+".go", "xyzservice", name); err != nil {
			return err
		}
		if err := replaceStringInFile("internal/"+name+"/"+name+".go", "xyzpb", name+"pb"); err != nil {
			return err
		}
		if err := replaceStringInFile(".Dockerfile", "xyzservice", name); err != nil {
			return err
		}
		if err := replaceStringInFile(".Dockerfile", "github.com/pedidopago/ms-template", module); err != nil {
			return err
		}

		if err := replaceStringInFile("internal/meta/meta.go", "xyzservice", name); err != nil {
			return err
		}

		_ = sh.Run("rm", "protos/"+name+"pb"+"/service.pb.go")

		if err := replaceStringInFile("protos/"+name+"pb/service.proto", "xyzpb", name+"pb"); err != nil {
			return err
		}
		if err := replaceStringInFile("protos/"+name+"pb/service.proto", "XYZService", name+"Service"); err != nil {
			return err
		}
		if err := replaceStringInFile("protos/"+name+"pb/gen.go", "xyzpb", name+"pb"); err != nil {
			return err
		}
		if err := replaceStringInFile("protos/"+name+"pb/client/client.go", "github.com/pedidopago/ms-template", module); err != nil {
			return err
		}
		if err := replaceStringInFile("protos/"+name+"pb/client/grpcdclient.go", "github.com/pedidopago/ms-template", module); err != nil {
			return err
		}
		if err := replaceStringInFile("protos/"+name+"pb/client/mockclient.go", "github.com/pedidopago/ms-template", module); err != nil {
			return err
		}
		if err := replaceStringInFile("protos/"+name+"pb/client/client.go", "xyzpb", name+"pb"); err != nil {
			return err
		}
		if err := replaceStringInFile("protos/"+name+"pb/client/grpcdclient.go", "xyzpb", name+"pb"); err != nil {
			return err
		}
		if err := replaceStringInFile("protos/"+name+"pb/client/mockclient.go", "xyzpb", name+"pb"); err != nil {
			return err
		}
		if err := replaceStringInFile("protos/"+name+"pb/client/client.go", "XYZServiceClient", strings.Title(name)+"ServiceClient"); err != nil {
			return err
		}
		if err := replaceStringInFile("protos/"+name+"pb/client/grpcdclient.go", "XYZServiceClient", strings.Title(name)+"ServiceClient"); err != nil {
			return err
		}
		if err := replaceStringInFile("protos/"+name+"pb/client/mockclient.go", "XYZServiceClient", strings.Title(name)+"ServiceClient"); err != nil {
			return err
		}
		if err := replaceStringInFile("internal/"+name+"/"+name+".go", "XYZService", strings.Title(name)); err != nil {
			return err
		}

		// replace module path
		if err := replaceStringInFile("cmd/"+name+"/main.go", "github.com/pedidopago/ms-template", module); err != nil {
			return err
		}
		if err := replaceStringInFile("go.mod", "github.com/pedidopago/ms-template", module); err != nil {
			return err
		}
		if err := replaceStringInFile("protos/"+name+"pb/service.proto", "github.com/pedidopago/ms-template", module); err != nil {
			return err
		}

		if err := ioutil.WriteFile(".name", []byte(name), 0644); err != nil {
			return err
		}

		// remove self install
		if err := regexpRemoveStringInFile("magefile.go", regexp.MustCompile(`(?s)// # rem`+`ove >>.*?// #`+` remo`+`ve <<`)); err != nil {
			return err
		}
	}

	return nil
}

// # remove <<

func replaceStringInFile(fname, oldv, newv string) error {
	fi, err := os.Stat(fname)
	if err != nil {
		return err
	}
	d, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	ds := strings.Replace(string(d), oldv, newv, -1)
	return ioutil.WriteFile(fname, []byte(ds), fi.Mode().Perm())
}

func regexpRemoveStringInFile(fname string, exp *regexp.Regexp) error {
	fi, err := os.Stat(fname)
	if err != nil {
		return err
	}
	d, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	d = exp.ReplaceAll(d, []byte(""))
	return ioutil.WriteFile(fname, d, fi.Mode().Perm())
}

func xs(v string) string {
	if strings.Contains(v, "%GIT_ROOT%") {
		rootp, _ := xcmd.CombinedString("git", "rev-parse", "--show-toplevel")
		v = strings.Replace(v, "%GIT_ROOT%", rootp, -1)
	}
	return v
}

func loadEnvs() {
	if _, e := os.Stat("dbm.env"); e == nil {
		godotenv.Overload("dbm.env")
	}
}
