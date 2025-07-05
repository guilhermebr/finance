package example

type UseCase struct {
	R Repository
}

func New(repo Repository) UseCase {
	return UseCase{R: repo}
}
