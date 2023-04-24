


cover: 
	./cover.sh

mockgen:
	mockgen -package=mocks -source=./controller/controller.go >> ./mocks/controller_mock.go
	mockgen -package=mocks -source=./storage/repo/review.go >> ./mocks/repo_mock.go