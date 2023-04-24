package rpc

import (
	"context"
	"fmt"
	libslog "monorepo/src/libs/log"
	"monorepo/src/review_service/configs"
	"monorepo/src/review_service/entity"
	"monorepo/src/review_service/mocks"

	pb "monorepo/src/idl/review_service"
	"os"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/test-go/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	cl      pb.ReviewServiceClient
	require *assert.Assertions
	logger  *zap.Logger
	tracer  opentracing.Tracer
	conf    *configs.Configuration
	l       libslog.Factory
	ctrl    *gomock.Controller
)

func loadReviewServiceClient(conf *configs.Configuration, tracer opentracing.Tracer) pb.ReviewServiceClient {
	connReview, err := grpc.Dial(
		fmt.Sprintf("%s%s", conf.ServerHost, conf.RPCPort),
		grpc.WithTransportCredentials(
			insecure.NewCredentials()),
	)
	if err != nil {
		panic(fmt.Errorf("kitchen service dial host:%s port:%s err: %s",
			conf.ServerHost, conf.RPCPort, err))
	}

	return pb.NewReviewServiceClient(connReview)
}

func TestMain(m *testing.M) {
	conf = configs.TestConfig()

	logger, _ = zap.NewDevelopment(
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
	)
	tracer = opentracing.GlobalTracer()

	zapLogger := logger.With(zap.String("service", "review_service"))
	l = libslog.NewFactory(zapLogger)

	cl = loadReviewServiceClient(conf, tracer)
	os.Exit(m.Run())
}

func TestHealth(t *testing.T) {

	require = assert.New(t)
	ctrl = gomock.NewController(t)
	defer ctrl.Finish()
	ct := mocks.NewMockController(ctrl)
	// ct.EXPECT()

	h := New(ct, l, tracer)
	server := Start(conf, tracer, h)
	defer server.Stop()
	fmt.Println("Health")

	ctx := context.Background()
	_, err := cl.HealthCheck(ctx, &pb.Empty{})
	require.NoError(err)

}

func TestCreate(t *testing.T) {

	require = assert.New(t)
	ctrl = gomock.NewController(t)
	defer ctrl.Finish()
	ct := mocks.NewMockController(ctrl)
	ctx := context.Background()
	id := uuid.NewString()
	name := "doston"
	x := "1232134234"
	x1 := "some message"

	ct.EXPECT().CreateReview(gomock.Any(), entity.CreateReviewReq{
		CreatorID:          id,
		CreatorName:        name,
		CreatorPhoneNumber: x,
		Message:            x1,
	}).Return(nil)
	h := New(ct, l, tracer)
	server := Start(conf, tracer, h)
	defer server.Stop()
	fmt.Println("Create")

	_, err := cl.Create(ctx, &pb.CreateReviewReq{
		CreatorId:          id,
		CreatorName:        name,
		CreatorPhoneNumber: x,
		Message:            x1,
	})
	require.NoError(err)
}

func TestAssign(t *testing.T) {

	require = assert.New(t)
	ctrl = gomock.NewController(t)
	defer ctrl.Finish()
	ct := mocks.NewMockController(ctrl)
	ctx := context.Background()
	id := uuid.NewString()

	ct.EXPECT().AssignReview(gomock.Any(), entity.AssignReviewReq{
		ID:           id,
		AssignedToID: id,
	}).Return(nil)
	h := New(ct, l, tracer)
	server := Start(conf, tracer, h)
	defer server.Stop()
	fmt.Println("Assign")

	_, err := cl.Assign(ctx, &pb.AssignReviewReq{
		Id:           id,
		AssignedToId: id,
	})
	require.NoError(err)
}
