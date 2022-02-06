//go:build mage
// +build mage

package main

// https://magefile.org

import (
	"bufio"
	"errors"
	"fmt"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/magefile/mage/sh"
	"github.com/mozillazg/go-slugify"
	"go.pedidopago.com.br/microservices/defaults"
	"go.pedidopago.com.br/microservices/xcmd"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
)

const (
	dbfiles = "file://%GIT_ROOT%/internal/database/migrations"
	dbcs    = "mysql://testuser:123456789@tcp(localhost)/ms_xyz?parseTime=true"
)

func Clean() {
	magelib.Clean()
}

func Migrate(version string) error {
	return magelib.Migrate(version)
}

func Newmigration(name string) error {
	return magelib.Newmigration(name)
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

func Docker() error {
	return magelib.Docker()
}

func DockerPush() error {
	return magelib.DockerPush()
}

func Generate() error {
	return magelib.Generate()
}

func Run() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("run.env")
	//
	nd, _ := ioutil.ReadFile(".name")
	name := strings.TrimSpace(string(nd))
	//
	dbcs := defaults.String(dbcs[8:], os.Getenv("DEV_DB_CS"), os.Getenv("DB_CS"), os.Getenv("DBCS"))
	//

	envs := make(map[string]string)
	envs["DBCS"] = dbcs

	if os.Getenv("GRPCD") == "" {
		envs["GRPCD"] = "dev"
	}
	if os.Getenv("GRPC_DISABLE_TLS") == "" {
		envs["GRPC_DISABLE_TLS"] = "1"
	}
	if os.Getenv("LISTEN") == "" {
		envs["LISTEN"] = ":15055"
	}
	if os.Getenv("AUTOMIGRATE") == "" {
		envs["AUTOMIGRATE"] = "1"
	}

	//
	_ = sh.RunWithV(envs, "go", "run", fmt.Sprintf("cmd/%s/main.go", name))
}

func RunAllTests() error {
	return magelib.RunAllTests()
}

func RunTestCoverage() error {
	return magelib.RunTestCoverage()
}

func Composerun() {
	if err := Devbuild(); err != nil {
		println("Composerun::Devbuild: " + err.Error())
		os.Exit(1)
	}
	_ = sh.RunV("docker-compose", "up")
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
		if err := sh.Run("mv", "proto/pedidopago/xyz", "proto/pedidopago/"+name); err != nil {
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

		_ = sh.Run("rm", "gen/proto/go/pedidopago/"+name+"/v1"+"/xyz_service.pb.go")
		_ = sh.Run("rm", "gen/proto/go/pedidopago/"+name+"/v1"+"/xyz_service_grpc.pb.go")

		if err := replaceStringInFile("proto/pedidopago/"+name+"v1/xyz_service.proto", "xyz", name); err != nil {
			return err
		}
		if err := replaceStringInFile("proto/pedidopago/"+name+"v1/xyz_service.proto", "XYZService", name+"Service"); err != nil {
			return err
		}
		if err := replaceStringInFile("gen/proto/go/pedidopago/"+name+"/v1/gen.go", "xyz", name); err != nil {
			return err
		}
		if err := replaceStringInFile("gen/proto/go/pedidopago/"+name+"/v1/helpers.go", "xyzpb", name+"pb"); err != nil {
			return err
		}
		if err := replaceStringInFile("gen/proto/go/pedidopago/"+name+"/v1/helpers.go", "xyzservice", name); err != nil {
			return err
		}
		if err := replaceStringInFile("gen/proto/go/pedidopago/"+name+"/v1/helpers.go", "XYZServiceClient", strings.Title(name)+"ServiceClient"); err != nil {
			return err
		}
		if err := replaceStringInFile("gen/proto/go/pedidopago/"+name+"/v1/helpers.go", "XYZServiceServer", strings.Title(name)+"ServiceServer"); err != nil {
			return err
		}
		if err := replaceStringInFile("gen/proto/go/pedidopago/"+name+"/v1/client/client.go", "github.com/pedidopago/ms-template", module); err != nil {
			return err
		}
		if err := replaceStringInFile("gen/proto/go/pedidopago/"+name+"/v1/client/grpcdclient.go", "github.com/pedidopago/ms-template", module); err != nil {
			return err
		}
		if err := replaceStringInFile("gen/proto/go/pedidopago/"+name+"/v1/client/mockclient.go", "github.com/pedidopago/ms-template", module); err != nil {
			return err
		}
		if err := replaceStringInFile("gen/proto/go/pedidopago/"+name+"/v1/client/client.go", "xyzv1", name+"v1"); err != nil {
			return err
		}
		if err := replaceStringInFile("gen/proto/go/pedidopago/"+name+"/v1/client/grpcdclient.go", "xyzv1", name+"v1"); err != nil {
			return err
		}
		if err := replaceStringInFile("gen/proto/go/pedidopago/"+name+"/v1/client/mockclient.go", "xyzv1", name+"v1"); err != nil {
			return err
		}
		if err := replaceStringInFile("gen/proto/go/pedidopago/"+name+"/v1/client/client.go", "XYZServiceClient", strings.Title(name)+"ServiceClient"); err != nil {
			return err
		}
		if err := replaceStringInFile("gen/proto/go/pedidopago/"+name+"/v1/client/grpcdclient.go", "XYZServiceClient", strings.Title(name)+"ServiceClient"); err != nil {
			return err
		}
		if err := replaceStringInFile("gen/proto/go/pedidopago/"+name+"/v1/client/mockclient.go", "XYZServiceClient", strings.Title(name)+"ServiceClient"); err != nil {
			return err
		}
		if err := replaceStringInFile("internal/"+name+"/"+name+".go", "XYZService", strings.Title(name)); err != nil {
			return err
		}
		if err := replaceStringInFile("README.md", "xyzservice", name); err != nil {
			return err
		}
		if err := replaceStringInFile("docker-compose.yml", "xyzservice", name); err != nil {
			return err
		}

		// replace module path
		if err := replaceStringInFile("cmd/"+name+"/main.go", "github.com/pedidopago/ms-template", module); err != nil {
			return err
		}
		if err := replaceStringInFile("go.mod", "github.com/pedidopago/ms-template", module); err != nil {
			return err
		}
		if err := replaceStringInFile("proto/pedidopago/"+name+"/v1/"+name+"_service.proto", "github.com/pedidopago/ms-template", module); err != nil {
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
