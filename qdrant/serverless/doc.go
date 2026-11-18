/*
Package serverless provides a client for Qdrant Serverless.

Experimental: the API may change without notice.

Point operations take the same requests as [qdrant.Client]:

	client, err := serverless.NewClient(&qdrant.Config{
		Host:   "serverless.example.cloud.qdrant.io",
		APIKey: "<your api key>",
	})
	...
	_, err = client.CreateCollection(ctx, "my-collection", &serverless.CollectionConfig{
		DenseVectors: map[string]*serverless.DenseVectorConfig{
			"": {Size: 1536, Distance: serverless.Distance_COSINE},
		},
	})

The generated types use the same proto names as the Go stubs of
github.com/qdrant/qdrant-cloud-public-api. A program importing both panics at
startup unless GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn is set.
*/
package serverless
