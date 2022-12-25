package main

import (
    "context"
    "fmt"
    "os"

    "dagger.io/dagger"
)

func main() {
    if err := build(context.Background()); err != nil {
        fmt.Println(err)
    }
}

func build(ctx context.Context) error {
    fmt.Println("Construyendo con Dagger")

    // Define una matriz de build
    oses := []string{"linux", "darwin"}
    arches := []string{"amd64", "arm64"}

    // Inicializa el cliente Dagger
    client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
    if err != nil {
        return err
    }
    defer client.Close()

    // Obtiene la referencia al proyecto local
    src := client.Host().Directory(".")

    // Crear un directorio vacio para almacenar las salidas de los builds
    outputs := client.Directory()

    // get `golang` image
    golang := client.Container().From("golang:latest")

    // Montar el repositorio clonado en la imagen `golang`
    golang = golang.WithMountedDirectory("/src", src).WithWorkdir("/src")

    for _, goos := range oses {
        for _, goarch := range arches {
            // Crear un directorio para cada os y arquitectura
            path := fmt.Sprintf("build/%s/%s/", goos, goarch)

            // Configurar GOARCH y GOOS en el ambiente de construcción
            build := golang.WithEnvVariable("GOOS", goos)
            build = build.WithEnvVariable("GOARCH", goarch)

            // construir la aplicación
            build = build.WithExec([]string{"go", "build", "-o", path})

            // get reference to build output directory in container
            outputs = outputs.WithDirectory(path, build.Directory(path))
        }
    }
    _, err = outputs.Export(ctx, ".")
    if err != nil {
        return err
    }

    return nil
}