package main

import (
	"fmt"

	"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AppImages struct {
	API      *dockerbuild.Image
	Migrator *dockerbuild.Image
}

func buildInternalImages(ctx *pulumi.Context) (*AppImages, error) {
	api, err := dockerbuild.NewImage(ctx, "api-image", &dockerbuild.ImageArgs{
		Context:    &dockerbuild.BuildContextArgs{Location: pulumi.String("../../")},
		Dockerfile: &dockerbuild.DockerfileArgs{Location: pulumi.String("../../Dockerfile")},
		BuildArgs:  pulumi.StringMap{"APP_NAME": pulumi.String("api")},
		Tags:       pulumi.StringArray{pulumi.String("my-go-playground/api:latest")},
		Push:       pulumi.Bool(false),
	})
	if err != nil {
		return nil, fmt.Errorf("new api image: %w", err)
	}

	migrator, err := dockerbuild.NewImage(ctx, "migrator-image", &dockerbuild.ImageArgs{
		Context:    &dockerbuild.BuildContextArgs{Location: pulumi.String("../../")},
		Dockerfile: &dockerbuild.DockerfileArgs{Location: pulumi.String("../../Dockerfile")},
		BuildArgs:  pulumi.StringMap{"APP_NAME": pulumi.String("migrator")},
		Tags:       pulumi.StringArray{pulumi.String("my-go-playground/migrator:latest")},
		Push:       pulumi.Bool(false),
	})
	if err != nil {
		return nil, fmt.Errorf("new migrator image: %w", err)
	}

	return &AppImages{API: api, Migrator: migrator}, nil
}
