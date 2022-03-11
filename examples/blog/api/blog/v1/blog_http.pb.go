// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// protoc-gen-go-http v2.0.0

package v1

import (
	context "context"
	http "github.com/SeeMusic/kratos/v2/transport/http"
	binding "github.com/SeeMusic/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

type BlogServiceHTTPServer interface {
	CreateArticle(context.Context, *CreateArticleRequest) (*CreateArticleReply, error)
	DeleteArticle(context.Context, *DeleteArticleRequest) (*DeleteArticleReply, error)
	GetArticle(context.Context, *GetArticleRequest) (*GetArticleReply, error)
	ListArticle(context.Context, *ListArticleRequest) (*ListArticleReply, error)
	UpdateArticle(context.Context, *UpdateArticleRequest) (*UpdateArticleReply, error)
}

func RegisterBlogServiceHTTPServer(s *http.Server, srv BlogServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/article/", _BlogService_CreateArticle0_HTTP_Handler(srv))
	r.PUT("/v1/article/{id}", _BlogService_UpdateArticle0_HTTP_Handler(srv))
	r.DELETE("/v1/article/{id}", _BlogService_DeleteArticle0_HTTP_Handler(srv))
	r.GET("/v1/article/{id}", _BlogService_GetArticle0_HTTP_Handler(srv))
	r.GET("/v1/article/", _BlogService_ListArticle0_HTTP_Handler(srv))
}

func _BlogService_CreateArticle0_HTTP_Handler(srv BlogServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateArticleRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/blog.api.v1.BlogService/CreateArticle")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateArticle(ctx, req.(*CreateArticleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateArticleReply)
		return ctx.Result(200, reply)
	}
}

func _BlogService_UpdateArticle0_HTTP_Handler(srv BlogServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateArticleRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/blog.api.v1.BlogService/UpdateArticle")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateArticle(ctx, req.(*UpdateArticleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateArticleReply)
		return ctx.Result(200, reply)
	}
}

func _BlogService_DeleteArticle0_HTTP_Handler(srv BlogServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteArticleRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/blog.api.v1.BlogService/DeleteArticle")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteArticle(ctx, req.(*DeleteArticleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteArticleReply)
		return ctx.Result(200, reply)
	}
}

func _BlogService_GetArticle0_HTTP_Handler(srv BlogServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetArticleRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/blog.api.v1.BlogService/GetArticle")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetArticle(ctx, req.(*GetArticleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetArticleReply)
		return ctx.Result(200, reply)
	}
}

func _BlogService_ListArticle0_HTTP_Handler(srv BlogServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListArticleRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/blog.api.v1.BlogService/ListArticle")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListArticle(ctx, req.(*ListArticleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListArticleReply)
		return ctx.Result(200, reply)
	}
}

type BlogServiceHTTPClient interface {
	CreateArticle(ctx context.Context, req *CreateArticleRequest, opts ...http.CallOption) (rsp *CreateArticleReply, err error)
	DeleteArticle(ctx context.Context, req *DeleteArticleRequest, opts ...http.CallOption) (rsp *DeleteArticleReply, err error)
	GetArticle(ctx context.Context, req *GetArticleRequest, opts ...http.CallOption) (rsp *GetArticleReply, err error)
	ListArticle(ctx context.Context, req *ListArticleRequest, opts ...http.CallOption) (rsp *ListArticleReply, err error)
	UpdateArticle(ctx context.Context, req *UpdateArticleRequest, opts ...http.CallOption) (rsp *UpdateArticleReply, err error)
}

type BlogServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewBlogServiceHTTPClient(client *http.Client) BlogServiceHTTPClient {
	return &BlogServiceHTTPClientImpl{client}
}

func (c *BlogServiceHTTPClientImpl) CreateArticle(ctx context.Context, in *CreateArticleRequest, opts ...http.CallOption) (*CreateArticleReply, error) {
	var out CreateArticleReply
	pattern := "/v1/article/"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/blog.api.v1.BlogService/CreateArticle"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *BlogServiceHTTPClientImpl) DeleteArticle(ctx context.Context, in *DeleteArticleRequest, opts ...http.CallOption) (*DeleteArticleReply, error) {
	var out DeleteArticleReply
	pattern := "/v1/article/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/blog.api.v1.BlogService/DeleteArticle"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *BlogServiceHTTPClientImpl) GetArticle(ctx context.Context, in *GetArticleRequest, opts ...http.CallOption) (*GetArticleReply, error) {
	var out GetArticleReply
	pattern := "/v1/article/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/blog.api.v1.BlogService/GetArticle"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *BlogServiceHTTPClientImpl) ListArticle(ctx context.Context, in *ListArticleRequest, opts ...http.CallOption) (*ListArticleReply, error) {
	var out ListArticleReply
	pattern := "/v1/article/"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/blog.api.v1.BlogService/ListArticle"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *BlogServiceHTTPClientImpl) UpdateArticle(ctx context.Context, in *UpdateArticleRequest, opts ...http.CallOption) (*UpdateArticleReply, error) {
	var out UpdateArticleReply
	pattern := "/v1/article/{id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/blog.api.v1.BlogService/UpdateArticle"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
