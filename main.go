package main

import (
    "context"
    "fmt"
    

    "dagger.io/dagger"
)

func main() {
    if err := build(context.Background()); err != nil {
        fmt.Println(err)
    }
}

func build(ctx context.Context) error {
    fmt.Println("Construyendo con Dagger")

    // Inicializa el cliente Dagger
    client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
    if err != nil {
        return err
    }
    defer client.Close()

    // Obtiene la referencia al proyecto local
    src := client.Host().Directory(".")

    // get `golang` image
    golang := client.Container().From("golang:latest")

    // Montar el repositorio clonado en la imagen `golang`
    golang = golang.WithMountedDirectory("/src", src).WithWorkdir("/src")

    // Define el comando de build de aplicación
    path := "build/"
    golang = golang.WithExec([]string{"go", "build", "-o", path})

    // Obtienen la referencia del directorio de build en el contenedor
    output := golang.Directory(path)

    // Escribe el contenido del directorio build/ al host
    _, err = output.Export(ctx, path)

    return nil
}