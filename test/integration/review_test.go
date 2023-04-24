package main

import (
	"context"
	"fmt"
	"log"
	pb "monorepo/src/idl/review_service"
	"monorepo/src/review_service/app"
	"monorepo/src/review_service/configs"

	"github.com/google/uuid"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func loadReviewServiceClient() pb.ReviewServiceClient {
	tracer := opentracing.GlobalTracer()
	conf := configs.TestConfig()
	connReview, err := grpc.Dial(
		fmt.Sprintf("%s:%s", conf.ServerHost, conf.RPCPort),
		grpc.WithTransportCredentials(
			insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(tracer)),
		grpc.WithStreamInterceptor(
			otgrpc.OpenTracingStreamClientInterceptor(tracer),
		),
	)
	if err != nil {
		panic(fmt.Errorf("kitchen service dial host: %s port:%s err: %s",
			conf.ServerHost, conf.RPCPort, err))
	}

	return pb.NewReviewServiceClient(connReview)
}

func main() {

	app := app.New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app.Start(ctx)

	cl := loadReviewServiceClient()

	fmt.Println("Health")
	_, err := cl.HealthCheck(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("HealthCheck: %v", err)
	}

	fmt.Println("CreateReview")
	id := uuid.NewString()
	_, err = cl.Create(ctx, &pb.CreateReviewReq{
		CreatorPhoneNumber: "1232134234",
		CreatorName:        "doston",
		CreatorId:          id,
		Message:            "some message",
	})
	if err != nil {
		log.Fatalf("Create: %v", err)
	}

	fmt.Println("GetAll")
	res, err := cl.GetAll(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("GetAll: %v", err)
	}
	if res.Count != 1 {
		log.Fatalf("GetAll: Count mismatch: got %v want %v", res.Count, 1)
	}

	fmt.Println("Assign")
	reviewID := res.Reviews[0].Id
	assignedToID := uuid.NewString()
	_, err = cl.Assign(ctx, &pb.AssignReviewReq{Id: reviewID, AssignedToId: assignedToID})
	if err != nil {
		log.Fatalf("Assign: %v", err)
	}
	res, err = cl.GetAll(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("GetAll: %v", err)
	}
	for _, r := range res.Reviews {
		if r.Id == reviewID {
			if r.AssignedToId != assignedToID {
				log.Printf("reviewID %v is not assigned to %v\n", r.Id, assignedToID)
				log.Fatalf("Assign: AssignedToId mismatch: got %v want %v", r.AssignedToId, assignedToID)
			}
		}
	}

	//resolve
	action := "made phone call"
	_, err = cl.Resolve(ctx, &pb.ResolveReviewReq{
		Id:             reviewID,
		AssignedToId:   assignedToID,
		TakenAction:    action,
		WithSuccess:    true,
		CustomerRating: 4,
	})
	if err != nil {
		log.Fatalf("Resolve: %v", err)
	}

	res, err = cl.GetAll(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("GetAll: %v", err)
	}
	for _, r := range res.Reviews {
		if r.Id == reviewID {
			if !r.IsResolved {
				log.Fatalf("Resolve: reviewID %v is not resolved", r.Id)
			}
		}
	}
}
