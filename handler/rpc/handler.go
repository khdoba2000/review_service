package rpc

import (
	"context"
	"fmt"
	pb "monorepo/src/idl/review_service"
	"monorepo/src/libs/constants"
	"monorepo/src/libs/log"
	"monorepo/src/libs/utils"
	"monorepo/src/review_service/configs"
	"monorepo/src/review_service/controller"
	"monorepo/src/review_service/entity"
	"net"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Handler defines a review service handler.
type Handler struct {
	ctrl   controller.Controller
	logger log.Factory
	tracer opentracing.Tracer
	pb.UnimplementedReviewServiceServer
}

// Start ...
func Start(config *configs.Configuration, tracer opentracing.Tracer, handler *Handler) *grpc.Server {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		otgrpc.OpenTracingServerInterceptor(tracer)),
		grpc.StreamInterceptor(
			otgrpc.OpenTracingStreamServerInterceptor(tracer)))

	pb.RegisterReviewServiceServer(grpcServer, handler)

	fmt.Println("Listen: ")

	//listenting tcp rpcport
	lis, err := net.Listen("tcp", config.RPCPort)
	if err != nil {
		fmt.Println("listening tcp error: ", err)
	}

	fmt.Println("crm server running on port ", config.RPCPort)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			panic(fmt.Errorf("failed to serve: %w", err))
		}
	}()
	return grpcServer
}

// New creates a metadata service controller.
func New(ctrl controller.Controller, logger log.Factory, tracer opentracing.Tracer) *Handler {
	return &Handler{
		ctrl:   ctrl,
		logger: logger,
		tracer: tracer,
	}
}

func (h *Handler) HealthCheck(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (h *Handler) Create(ctx context.Context, req *pb.CreateReviewReq) (*pb.Empty, error) {

	//TODO: use mapper to map the request
	err := h.ctrl.CreateReview(ctx, entity.CreateReviewReq{
		CreatorID:          req.CreatorId,
		CreatorName:        req.CreatorName,
		CreatorPhoneNumber: req.CreatorPhoneNumber,
		Message:            req.Message,
	})
	if err != nil {
		h.logger.For(ctx).Error("ctrl.CreateReview failed", zap.Error(err))
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (h *Handler) Assign(ctx context.Context, req *pb.AssignReviewReq) (*pb.Empty, error) {

	//TODO: use mapper to map the request
	err := h.ctrl.AssignReview(ctx, entity.AssignReviewReq{
		ID:           req.Id,
		AssignedToID: req.AssignedToId,
	})
	if err != nil {
		h.logger.For(ctx).Error("ctrl.AssignReview failed", zap.Error(err))
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (h *Handler) Resolve(ctx context.Context, req *pb.ResolveReviewReq) (*pb.Empty, error) {

	//TODO: use mapper to map the request
	err := h.ctrl.ResolveReview(ctx, entity.ResolveReviewReq{
		ID:             req.Id,
		AssignedToID:   req.AssignedToId,
		TakenAction:    req.TakenAction,
		CustomerRating: uint8(req.CustomerRating),
		WithSuccess:    req.WithSuccess,
	})
	if err != nil {
		h.logger.For(ctx).Error("ctrl.ResolveReview failed", zap.Error(err))
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (h *Handler) GetAll(ctx context.Context, r *pb.Empty) (*pb.GetAllResponse, error) {
	//TODO: add pagination
	//TODO: use mapper to map the request
	res, err := h.ctrl.GetAll(ctx)
	if err != nil {
		h.logger.For(ctx).Error("ctrl.CreateReview failed", zap.Error(err))
		return nil, err
	}

	reviews := utils.MapSlice(res, func(r entity.ReviewOut) *pb.ReviewOut {
		return &pb.ReviewOut{
			Id:           r.ID,
			CreatorId:    r.CreatorID,
			CreatorName:  r.CreatorName,
			Message:      r.Message,
			IsResolved:   r.IsResolved,
			AssignedToId: r.AssignedToID,
			ResolvedAt:   r.ResolvedAt.Format(constants.OnlyDateAndTimeWithDots),
		}
	})

	return &pb.GetAllResponse{Reviews: reviews, Count: uint32(len(reviews))}, nil
}
