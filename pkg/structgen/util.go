// Package structgen
// @author Daud Valentino
package structgen

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/util"
)

func fileExist(fName string) bool {
	_, err := os.Stat(fName)
	return errors.Is(err, os.ErrNotExist) == false
}

func contractName(v string) string {
	if util.SubStringRight(v, 1) == "e" {
		return v + "r"
	}
	return v + "er"
}

func createUseCaseList(cols []Column, pkgName, tableName string) {
	tName := tableName
	if util.SubStringRight(tName, 1) == "s" {
		tName = util.SubStringLeft(tName, len(tName)-1)
	}

	pkgName = strings.ReplaceAll(pkgName, "_", "")

	if len(pkgName) == 0 {
		pkgName = strings.ReplaceAll(util.ToSnakeCase(tableName), "_", "")
	}

	pkgName = fmt.Sprintf("%s/%s", uCasePath, pkgName)
	if !util.PathExist(pkgName) {
		os.MkdirAll(pkgName, os.ModePerm)
	}

	fName := fmt.Sprintf("%s/list.go", pkgName)
	if fileExist(fName) {
		fmt.Println(fmt.Sprintf("file repo already exist %s", fName))
		return
	}

	eFile, err := os.Create(fName)
	defer eFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	tpl, err := template.ParseFiles(`./pkg/structgen/ucase_list.tpl`)
	if err != nil {
		log.Fatal(err)
	}

	pn := strings.Split(pkgName, "/")

	pNameSpace := pkgName
	if len(pn) > 1 {
		pNameSpace = pn[len(pn)-1]
	}

	err = tpl.Execute(eFile, UseCaseTemplate{
		FileName:         "list",
		TableName:        tableName,
		StructName:       util.ToCamelCase(tName),
		EntityName:       util.UpperFirst(util.ToCamelCase(tName)),
		PackageName:      pNameSpace,
		ModuleName:       util.GetModuleName(),
		RepoContractName: contractName(util.UpperFirst(util.ToCamelCase(tName))),
		Columns:          cols,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("success create use case list ", fName)
}

func createUseCaseStorer(cols []Column, importPackage []string, packageName, tableName string) {
	tName := tableName
	if util.SubStringRight(tName, 1) == "s" {
		tName = util.SubStringLeft(tName, len(tName)-1)
	}

	//pkgName := strings.ReplaceAll(util.ToSnakeCase(tableName), "_", "")
	pkgName := strings.ReplaceAll(packageName, "_", "")
	pkgName = fmt.Sprintf("%s/%s", uCasePath, packageName)
	if !util.PathExist(pkgName) {
		os.MkdirAll(pkgName, os.ModePerm)
	}

	fName := fmt.Sprintf("%s/store.go", pkgName)
	if fileExist(fName) {
		fmt.Println(fmt.Sprintf("file repo already exist %s", fName))
		return
	}

	eFile, err := os.Create(fName)
	defer eFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	tpl, err := template.ParseFiles(`./pkg/structgen/ucase_store.tpl`)
	if err != nil {
		log.Fatal(err)
	}

	pn := strings.Split(pkgName, "/")

	pNameSpace := pkgName
	if len(pn) > 1 {
		pNameSpace = pn[len(pn)-1]
	}

	en := util.UpperFirst(util.ToCamelCase(tName))
	err = tpl.Execute(eFile, UseCaseTemplate{
		FileName:         "store",
		TableName:        tableName,
		StructName:       util.ToCamelCase(tName),
		EntityName:       util.UpperFirst(en),
		PackageName:      pNameSpace,
		RepoContractName: contractName(util.UpperFirst(en)),
		ModuleName:       util.GetModuleName(),
		Columns:          cols,
		ImportPackage:    importPackage,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("success create use case storer ", fName)
}

func createUseCaseLogic(packageName, tableName string) {
	tName := tableName
	if util.SubStringRight(tName, 1) == "s" {
		tName = util.SubStringLeft(tName, len(tName)-1)
	}

	//pkgName := strings.ReplaceAll(util.ToSnakeCase(tableName), "_", "")
	pkgName := strings.ReplaceAll(packageName, "_", "")
	pkgName = fmt.Sprintf("%s/%s", uCasePath, packageName)
	if !util.PathExist(pkgName) {
		os.MkdirAll(pkgName, os.ModePerm)
	}

	fln := util.ToSnakeCase(tName)
	fName := fmt.Sprintf("%s/%s.go", pkgName, fln)
	if fileExist(fName) {
		fmt.Println(fmt.Sprintf("file repo already exist %s", fName))
		return
	}

	eFile, err := os.Create(fName)
	defer eFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	tpl, err := template.ParseFiles(`./pkg/structgen/ucase_store.tpl`)
	if err != nil {
		log.Fatal(err)
	}

	pn := strings.Split(pkgName, "/")

	pNameSpace := pkgName
	if len(pn) > 1 {
		pNameSpace = pn[len(pn)-1]
	}

	en := util.UpperFirst(util.ToCamelCase(tName))
	err = tpl.Execute(eFile, UseCaseTemplate{
		FileName:         fln,
		TableName:        tableName,
		StructName:       util.ToCamelCase(tName),
		EntityName:       util.UpperFirst(en),
		PackageName:      pNameSpace,
		RepoContractName: contractName(util.UpperFirst(en)),
		ModuleName:       util.GetModuleName(),
		Columns:          nil,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("success create use case storer ", fName)
}

func createUseCaseUpdater(cols []Column, packageName, tableName string) {
	tName := tableName
	if util.SubStringRight(tName, 1) == "s" {
		tName = util.SubStringLeft(tName, len(tName)-1)
	}

	//pkgName := strings.ReplaceAll(util.ToSnakeCase(tableName), "_", "")
	pkgName := strings.ReplaceAll(packageName, "_", "")
	pkgName = fmt.Sprintf("%s/%s", uCasePath, packageName)
	if !util.PathExist(pkgName) {
		os.MkdirAll(pkgName, os.ModePerm)
	}

	//fln := util.ToSnakeCase("update")
	fName := fmt.Sprintf("%s/update.go", pkgName)
	if fileExist(fName) {
		fmt.Println(fmt.Sprintf("file repo already exist %s", fName))
		return
	}

	eFile, err := os.Create(fName)
	defer eFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	tpl, err := template.ParseFiles(`./pkg/structgen/ucase_update.tpl`)
	if err != nil {
		log.Fatal(err)
	}

	pn := strings.Split(pkgName, "/")

	pNameSpace := pkgName
	if len(pn) > 1 {
		pNameSpace = pn[len(pn)-1]
	}

	en := util.UpperFirst(util.ToCamelCase(tName))
	err = tpl.Execute(eFile, UseCaseTemplate{
		FileName:         "update",
		TableName:        tableName,
		StructName:       util.ToCamelCase(tName),
		EntityName:       util.UpperFirst(en),
		PackageName:      pNameSpace,
		RepoContractName: contractName(util.UpperFirst(en)),
		ModuleName:       util.GetModuleName(),
		Columns:          cols,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("success create use case updater ", fName)
}

func packageName(tableName string) string {
	pkgName := strings.ReplaceAll(util.ToSnakeCase(tableName), "_", "")

	if util.SubStringRight(pkgName, 1) != "s" {
		pkgName = fmt.Sprintf("%s%s", pkgName, "s")
	}

	return pkgName
}